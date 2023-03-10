package gvg

import (
	"errors"
	"fmt"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/util/gconv"
	"github.com/gogf/gf/util/gvalid"
	"github.com/vivid-vvo/vvbot/app/model/clan_group"
	"github.com/vivid-vvo/vvbot/app/model/gvg_challenge"
	"github.com/vivid-vvo/vvbot/app/model/gvg_group"
	"github.com/vivid-vvo/vvbot/app/model/gvg_member_extra"
	"github.com/vivid-vvo/vvbot/app/service/check"
	"github.com/vivid-vvo/vvbot/app/service/clan"
	"github.com/vivid-vvo/vvbot/app/service/user"
	time2 "github.com/vivid-vvo/vvbot/library/time"
	"github.com/vivid-vvo/vvbot/qqbot/model/bot"
	"time"
)

// GvgGroupCreateInput
type GvgGroupCreateInput struct {
	GroupName  string `v:"required#公会组名不能为空"`
	GameServer string `v:"in:CN,TW,JP,KR#游戏服务器应当在国,台,日,韩选择其一"`
	GvgName    string `v:"required#请输入本次公会战名称"`
}

// ReportCauseHarmInput
type ReportCauseHarmInput struct {
	Qqid   int64
	Damage int
	// 是否尾刀
	IsContinue int
	// 是否昨日刀
	IsYesterday int
	// 留言
	Message string
}

// GvgGroupCreate 创建一个新的公会战
func GvgGroupCreate(r *ghttp.Request, ggc *GvgGroupCreateInput, qqid int64) error {
	// 输入参数检查
	if e := gvalid.CheckStruct(ggc, nil); e != nil {
		return errors.New(e.FirstString())
	}
	if qqid == 0 {
		qqid = user.GetLoginData2(r).Qqid
	}
	var entity *gvg_group.Entity
	if err := gconv.Struct(ggc, &entity); err != nil {
		return err
	}

	data, err := clan_group.GetClanGroupAtName(ggc.GroupName)
	if err != nil {
		return err
	}
	if data == nil {
		return errors.New(fmt.Sprintf("指定公会不存在！"))
	}
	entity.CreateQqid = qqid
	entity.GroupId = data.GroupId
	entity.Time = gtime.Now().Unix()
	if err = gvg_group.GvgGroupCreate(entity); err != nil {
		return err
	}
	return nil
}

// ReportCauseHarm 报刀
func ReportCauseHarm(r *ghttp.Request, rch *ReportCauseHarmInput) error {
	// 输入参数检查
	if e := gvalid.CheckStruct(rch, nil); e != nil {
		return errors.New(e.FirstString())
	}
	var agentQqid int64
	if rch.Qqid == 0 {
		rch.Qqid = user.GetLoginData2(r).Qqid
	} else {
		agentQqid = rch.Qqid
	}
	// 查询公会信息
	clanData, err := clan.GetClanDataAtQqid(r, rch.Qqid)
	if err != nil {
		return errors.New(fmt.Sprintf("内部错误"))
	}
	gvgData, err := gvg_group.GetGvgGroupData(clanData.GvgId)
	if err != nil {
		return errors.New(fmt.Sprintf("内部错误"))
	}
	if gvgData == nil {
		return errors.New(fmt.Sprintf("当前公会未开启公会战"))
	}
	if rch.IsContinue == 1 {
		rch.Damage = gvgData.BossHp
	}
	if rch.Damage > gvgData.BossHp {
		return errors.New(fmt.Sprintf("报刀伤害大于BOSS当前血量，请使用尾刀"))
	}
	// 每天5点刷新
	dayTimeS := time2.GetPcrDayStartTimeToUnix(clanData.GameServer)
	dayTimeE := time2.GetPcrDayEndTimeToUnix(clanData.GameServer)
	// 获取今日刀，过滤尾刀
	num, err := gvg_challenge.FindCount("qqid=? and gvg_id=? and is_continue=0 and challenge_time>=? and challenge_time<=?", rch.Qqid, clanData.GvgId, dayTimeS, dayTimeE)
	if err != nil {
		return errors.New(fmt.Sprintf("内部错误"))
	}
	if num >= 3 {
		return errors.New(fmt.Sprintf("今日上报次数已达到3次"))
	}
	var challengeTime int64

	// 刀是否加入到昨日
	if rch.IsYesterday == 1 {
		challengeTime = dayTimeS - 1
	} else {
		challengeTime = time.Now().Unix()
	}
	gvgChallenge := gvg_challenge.Entity{
		Qqid:            rch.Qqid,
		GvgId:           clanData.GvgId,
		ClanGroupId:     clanData.GroupId,
		AgentQqid:       agentQqid,
		ChallengeDamage: rch.Damage,
		IsContinue:      rch.IsContinue,
		Message:         rch.Message,
		ChallengeTime:   challengeTime,
	}
	if _, err := gvg_challenge.Model.FieldsEx("challenge_id").Insert(gvgChallenge); err != nil {
		return errors.New(fmt.Sprintf("内部错误"))
	}
	// BOSS数据计算刷新
	if err := check.BossHpCount(clanData.GvgId, clanData.GameServer); err != nil {
		return errors.New(fmt.Sprintf("报刀成功，但BOSS血量修改失败"))
	}
	return nil
}

// BackCauseHarm 撤销指定成员一刀
func BackCauseHarm(r *ghttp.Request, qqid int64) error {
	if qqid == 0 {
		qqid = user.GetLoginData2(r).Qqid
	}
	clanData, err := clan.GetClanDataAtQqid(r, qqid)
	if err != nil {
		return errors.New(fmt.Sprintf("内部错误"))
	}
	causeHarmLostAtMember, err := GetCauseHarmLostAtMember(r, qqid)
	if err != nil {
		return err
	}
	if causeHarmLostAtMember == nil {
		return errors.New(fmt.Sprintf("无出刀记录"))
	}
	/*
		dayTimeS := time2.GetPcrDayStartTimeToUnix(clanData.GameServer)
		dayTimeE := time2.GetPcrDayEndTimeToUnix(clanData.GameServer)
	*/
	/*
		if _, err := gvg_challenge.Model.Limit(1).Order("challenge_time dec").Delete("qqid=? and time>? and time<? and gvg_id=?", qqid, dayTimeS, dayTimeE, clanData.GvgId); err != nil {
			return errors.New(fmt.Sprintf("内部错误"))
		}*/
	// BOSS数据计算刷新
	if err := check.BossHpCount(clanData.GvgId, clanData.GameServer); err != nil {
		return errors.New(fmt.Sprintf("报刀成功，但BOSS血量修改失败"))
	}
	return nil
}

// BackCauseHarmLost 撤销最近一刀
func BackCauseHarmLost(r *ghttp.Request, qqid int64) error {
	if qqid == 0 {
		qqid = user.GetLoginData2(r).Qqid
	}
	clanData, err := clan.GetClanDataAtQqid(r, qqid)
	if err != nil {
		return errors.New(fmt.Sprintf("内部错误"))
	}
	causeHarmLost, err := GetCauseHarmLost(r, qqid)
	if err != nil {
		return errors.New(fmt.Sprintf("内部错误"))
	}
	if causeHarmLost == nil {
		return errors.New(fmt.Sprintf("无出刀数据"))
	}
	/*dayTimeS := time2.GetPcrDayStartTimeToUnix(clanData.GameServer)
	dayTimeE := time2.GetPcrDayEndTimeToUnix(clanData.GameServer)*/
	/*if _, err := gvg_challenge.Model.Limit(1).Order("challenge_time dec").Delete("time>? and time<? and gvg_id=?", dayTimeS, dayTimeE, clanData.GvgId); err != nil {
		return errors.New(fmt.Sprintf("内部错误"))
	}*/
	// BOSS数据计算刷新
	if err := check.BossHpCount(clanData.GvgId, clanData.GameServer); err != nil {
		return errors.New(fmt.Sprintf("报刀成功，但BOSS血量修改失败"))
	}
	return nil
}

// GetCauseHarmLost 获取最近一刀
func GetCauseHarmLost(r *ghttp.Request, qqid int64) (*gvg_challenge.Entity, error) {
	if qqid == 0 {
		qqid = user.GetLoginData2(r).Qqid
	}
	clanData, err := clan.GetClanDataAtQqid(r, qqid)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("内部错误"))
	}
	dayTimeS := time2.GetPcrDayStartTimeToUnix(clanData.GameServer)
	dayTimeE := time2.GetPcrDayEndTimeToUnix(clanData.GameServer)
	gvgChallenge, err := gvg_challenge.Model.Limit(1).Order("challenge_time dec").FindOne("time>? and time<? and gvg_id=?", dayTimeS, dayTimeE, clanData.GvgId)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("内部错误"))
	}
	return gvgChallenge, nil
}

// GetCauseHarmLostAtMember 获取指定成员最近一刀
func GetCauseHarmLostAtMember(r *ghttp.Request, qqid int64) (*gvg_challenge.Entity, error) {
	if qqid == 0 {
		qqid = user.GetLoginData2(r).Qqid
	}
	clanData, err := clan.GetClanDataAtQqid(r, qqid)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("内部错误"))
	}
	dayTimeS := time2.GetPcrDayStartTimeToUnix(clanData.GameServer)
	dayTimeE := time2.GetPcrDayEndTimeToUnix(clanData.GameServer)
	gvgChallenge, err := gvg_challenge.Model.Limit(1).Order("challenge_time dec").FindOne("qqid=? and time>? and time<? and gvg_id=?", qqid, dayTimeS, dayTimeE, clanData.GvgId)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("内部错误"))
	}
	return gvgChallenge, nil
}

// ChangeBossHp 补刀调血
func ChangeBossHp(r *ghttp.Request, qqid int64, hp int) error {
	if qqid == 0 {
		qqid = user.GetLoginData2(r).Qqid
	}
	// 查询公会信息
	clanData, err := clan.GetClanDataAtQqid(r, qqid)
	if err != nil {
		return err
	}
	gvgData, err := gvg_group.GetGvgGroupData(clanData.GvgId)
	if err != nil {
		return err
	}
	if gvgData == nil {
		return errors.New(fmt.Sprintf("当前公会未开启公会战"))
	}
	damage := hp - gvgData.BossHp
	gvgChallenge := gvg_challenge.Entity{
		GvgId:           clanData.GvgId,
		ClanGroupId:     clanData.GroupId,
		ChallengeDamage: damage,
		RepairType:      1,
		IsContinue:      0,
		Message:         "补刀，BOSS血量调整",
		ChallengeTime:   time.Now().Unix(),
	}
	if err := gvg_challenge.ReportChallenge(gvgChallenge); err != nil {
		return err
	}

	// BOSS数据计算刷新
	bot.Send(clanData.BindQqGroup, 2, check.GetBossStateStr(clanData.GroupId))

	return nil
}

// ApplyChallenge 申请出道(刀)
func ApplyChallenge(r *ghttp.Request, qqid int64) error {
	if qqid == 0 {
		qqid = user.GetLoginData2(r).Qqid
	}
	clanData, err := clan.GetClanDataAtQqid(r, qqid)
	if err != nil {
		return errors.New(fmt.Sprintf("内部错误"))
	}
	isLook, err := gvg_group.BossIsLock(clanData.GvgId)
	if err != nil {
		return errors.New(fmt.Sprintf("内部错误"))
	}
	if isLook {
		return errors.New(fmt.Sprintf("BOSS已锁定"))
	}
	_, gvgGroupData, err := check.GetGvgGroupDataAtGroupIdToCheck(clanData.GroupId)
	if err != nil {
		return errors.New(fmt.Sprintf("内部错误"))
	}
	if gvgGroupData.ChallengeStratQqid != 0 {
		return errors.New(fmt.Sprintf("申请失败，当前有人正在挑战BOSS"))
	}
	isUpTree, err := GetIsUpTree(r, qqid)
	if err != nil {
		return err
	}
	if isUpTree {
		return errors.New(fmt.Sprintf("您在树上！无法继续申请出刀。请先下树！"))
	}
	if err = gvg_group.ApplyChallenge(qqid, gvgGroupData.GvgId); err != nil {
		return err
	}
	return nil
}

// CancelChallenge 取消申请出刀
func CancelChallenge(r *ghttp.Request, qqid int64) error {
	if qqid == 0 {
		qqid = user.GetLoginData2(r).Qqid
	}
	clanData, err := clan.GetClanDataAtQqid(r, qqid)
	if err != nil {
		return errors.New(fmt.Sprintf("内部错误"))
	}
	_, gvgGroupData, err := check.GetGvgGroupDataAtGroupIdToCheck(clanData.GroupId)
	if err != nil {
		return errors.New(fmt.Sprintf("内部错误"))
	}
	if gvgGroupData.ChallengeStratQqid == 0 {
		return errors.New(fmt.Sprintf("取消失败，当前没有出刀"))
	}
	if gvgGroupData.ChallengeStratQqid != qqid {
		return errors.New(fmt.Sprintf("取消失败，权限不足，不能取消他人出刀。管理员请使用”解锁“"))
	}
	if err = gvg_group.CancelChallenge(gvgGroupData.GvgId); err != nil {
		return errors.New(fmt.Sprintf("内部错误"))
	}
	return nil
}

func UnlChallenge(r *ghttp.Request, qqid int64) error {
	if qqid == 0 {
		qqid = user.GetLoginData2(r).Qqid
	}
	clanData, err := clan.GetClanDataAtQqid(r, qqid)
	if err != nil {
		return errors.New(fmt.Sprintf("内部错误"))
	}
	_, gvgGroupData, err := check.GetGvgGroupDataAtGroupIdToCheck(clanData.GroupId)
	if err != nil {
		return errors.New(fmt.Sprintf("内部错误"))
	}
	var isAdmin bool
	if gvgGroupData.ChallengeStratQqid > 0 {
		if gvgGroupData.ChallengeStratQqid != qqid {
			isAdmin = check.CheckAuthorityGroup(qqid, check.AuthGvgAdmin, clanData.GroupId)
			if err != nil {
				return errors.New(fmt.Sprintf("内部错误"))
			}
			// 必须管理员权限才能取消别人出刀, 3分钟后谁都可以解锁
			if !isAdmin && time.Now().Unix()-gvgGroupData.ChallengeStratTime > 60*3 {
				return errors.New(fmt.Sprintf("解锁失败，权限不足，不能解锁他人出刀。"))
			}
		}
		if err = gvg_group.CancelChallenge(gvgGroupData.GvgId); err != nil {
			return errors.New(fmt.Sprintf("内部错误"))
		}
	}
	if !isAdmin {
		return errors.New(fmt.Sprintf("解锁失败，权限不足。"))
	}
	if err = gvg_group.UnLockBoss(gvgGroupData.GvgId); err != nil {
		return errors.New(fmt.Sprintf("内部错误"))
	}
	return nil
}

// ReportDaySL 报告今日SL;
func ReportDaySL(r *ghttp.Request, qqid int64, state int) error {
	//var agentQqid          int64
	if qqid == 0 {
		qqid = user.GetLoginData2(r).Qqid
	} else {
		// agentQqid = qqid
	}
	if isSL, err := GetDaySL(r, qqid); err != nil || isSL {
		if isSL {
			return errors.New(fmt.Sprintf("今日已SL过"))
		}
		return errors.New(fmt.Sprintf("内部错误"))
	}
	return nil //gvg_member_extra.ReportMemberExtra(qqid, agentQqid, 1)
}

// CancelDaySL 取消SL
func CancelDaySL(r *ghttp.Request, qqid int64, state int) error {
	if qqid == 0 {
		qqid = user.GetLoginData2(r).Qqid
	}
	if isSL, err := GetDaySL(r, qqid); err != nil || !isSL {
		if !isSL {
			return errors.New(fmt.Sprintf("今日未sl"))
		}
		return errors.New(fmt.Sprintf("内部错误"))
	}
	clanData, err := clan.GetClanDataAtQqid(r, qqid)
	if err != nil {
		return errors.New(fmt.Sprintf("内部错误"))
	}
	return gvg_member_extra.CancelDaySL(qqid, clanData.GvgId, clanData.GameServer)
}

// GetDaySL 获取是今日否SL;
func GetDaySL(r *ghttp.Request, qqid int64) (bool, error) {
	if qqid == 0 {
		qqid = user.GetLoginData2(r).Qqid
	}

	/*clanData, err := clan.GetClanDataAtQqid(r, qqid)
	if err != nil {
		return false, errors.New(fmt.Sprintf("内部错误"))
	}
	dayTimeS := time2.GetPcrDayStartTimeToUnix(clanData.GameServer)
	dayTimeE := time2.GetPcrDayEndTimeToUnix(clanData.GameServer)
	return gvg_member_extra.GetMemberExtra(qqid, 1, int(dayTimeS), int(dayTimeE))*/

	return false, nil
}

// ReportUPTree 报告挂树
func ReportUPTree(r *ghttp.Request, qqid int64) error {
	//var agentQqid          int64
	if qqid == 0 {
		qqid = user.GetLoginData2(r).Qqid
	} else {
		//agentQqid = qqid
	}
	if isUpTree, err := GetIsUpTree(r, qqid); err != nil || isUpTree {
		if isUpTree {
			return errors.New(fmt.Sprintf("您已在树上！"))
		}
		return errors.New(fmt.Sprintf("内部错误"))
	}
	return nil // gvg_member_extra.ReportMemberExtra(qqid, agentQqid, 2)
}

// GetIsUpTree 获取是当前是否挂树;
func GetIsUpTree(r *ghttp.Request, qqid int64) (bool, error) {
	/*if qqid == 0 {
		qqid = user.GetLoginData2(r).Qqid
	}
	clanData, err := clan.GetClanDataAtQqid(r, qqid)
	if err != nil {
		return false, errors.New(fmt.Sprintf("内部错误"))
	}
	dayTimeE := time2.GetNowTimeToZone(clanData.GameServer).Unix()
	// 挂树最长时间为1小时
	dayTimeS := dayTimeE - 60*60
	return gvg_member_extra.GetMemberExtra(qqid, 2, int(dayTimeS), int(dayTimeE))
	*/
	return false, nil
}

// ReportDownTree 报告下树 state==1 结算下树, state==2 BOSS死亡后下树
func ReportDownTree(r *ghttp.Request, qqid int64, state int) error {
	if qqid == 0 {
		qqid = user.GetLoginData2(r).Qqid
	}
	if isUpTree, err := GetIsUpTree(r, qqid); err != nil || !isUpTree {
		if !isUpTree {
			return errors.New(fmt.Sprintf("您未已在树上"))
		}
		return errors.New(fmt.Sprintf("内部错误"))
	}
	clanData, err := clan.GetClanDataAtQqid(r, qqid)
	if err != nil {
		return errors.New(fmt.Sprintf("内部错误"))
	}
	return gvg_member_extra.ReportDownTree(qqid, clanData.GvgId, state)
}
