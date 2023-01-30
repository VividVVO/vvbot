package gvg

import (
	"errors"
	"fmt"

	"github.com/gogf/gf/net/ghttp"

	"github.com/gogf/gf/util/gconv"

	"github.com/vivid-vvo/vvbot/app/model/clan_member"
	"github.com/vivid-vvo/vvbot/app/model/gvg_challenge"
	"github.com/vivid-vvo/vvbot/app/model/gvg_group"
	"github.com/vivid-vvo/vvbot/app/model/gvg_member_extra"
	"github.com/vivid-vvo/vvbot/app/service/check"

	"github.com/vivid-vvo/vvbot/app/service/user"
	time2 "github.com/vivid-vvo/vvbot/library/time"
	"github.com/vivid-vvo/vvbot/qqbot/model/bot"
)

// GetChallengeAtQQ 获取指定QQ, 指定时间内所有刀，不包括修正刀
func GetChallengeAtQQ(r *ghttp.Request, qqid int64, clanGroupId int, TimeType string) ([]*gvg_challenge.Entity, error) {
	if qqid == 0 {
		qqid = user.GetLoginData2(r).Qqid
	}
	clanGroup, err := check.GetClanGroupAndChack(clanGroupId)
	if err != nil {
		return nil, err
	}
	var timeS, timeE int64
	switch TimeType {
	case "day":
		timeS = time2.GetPcrDayEndTimeToUnix(clanGroup.GameServer)
		timeE = time2.GetPcrDayEndTimeToUnix(clanGroup.GameServer)
	case "yesterday":
		timeS = time2.GetPcrYesterdayStartTimeToUnix(clanGroup.GameServer)
		timeE = time2.GetPcrYesterdayEndTimeToUnix(clanGroup.GameServer)
	default:
		// 默认所有刀
		timeS = 0
		timeE = time2.GetPcrDayEndTimeToUnix(clanGroup.GameServer)
	}
	gvgChallenges, err := gvg_challenge.GetChallengeAtQQ(qqid, clanGroup.GvgId, timeS, timeE)
	if err != nil {
		return nil, err
	}
	return gvgChallenges, nil
}

type GetAllChallengeInput struct {
	List    []*gvg_challenge.Entity `json:"list"`
	TimeStr string                  `json:"timeStr"`
}

func GetAllChallenge(r *ghttp.Request, clanGroupId int, timeType string, startTime string, endTime string) (*GetAllChallengeInput, error) {
	clanGroup, err := check.GetClanGroupAndChack(clanGroupId)
	if err != nil {
		return nil, err
	}
	var timeS, timeE int64
	switch timeType {
	case "day":
		timeS = time2.GetPcrDayStartTimeToUnix(clanGroup.GameServer)
		timeE = time2.GetPcrDayEndTimeToUnix(clanGroup.GameServer)
	case "yesterday":
		timeS = time2.GetPcrYesterdayStartTimeToUnix(clanGroup.GameServer)
		timeE = time2.GetPcrYesterdayEndTimeToUnix(clanGroup.GameServer)
	case "all":
		timeS = 0
		timeE = time2.GetPcrDayEndTimeToUnix(clanGroup.GameServer)
	case "time":
		timeS = time2.GetPcrStartTimeToUnixAtStr(clanGroup.GameServer, startTime)
		if endTime != "" {
			timeE = time2.GetPcrEndTimeToUnixAtStr(clanGroup.GameServer, endTime)
		} else {
			timeE = timeS
		}
		timeE = timeE + 24*60*60

	default:
		// 默认今日刀
		timeS = time2.GetPcrDayStartTimeToUnix(clanGroup.GameServer)
		timeE = time2.GetPcrDayEndTimeToUnix(clanGroup.GameServer)
	}
	gvgChallenges, err := gvg_challenge.GetAllChallengeAtTime(clanGroup.GvgId, timeS, timeE)
	if err != nil {
		return nil, err
	}
	allChallengeInput := &GetAllChallengeInput{
		List:    gvgChallenges,
		TimeStr: time2.GetTimeAtUnixToZone(clanGroup.GameServer, timeS).Format("2006-01-02"),
	}
	return allChallengeInput, nil
}

func GetAllSlState(r *ghttp.Request, clanGroupId int, TimeType string) (*[]gvg_member_extra.GetAllSlStateEntity, error) {
	clanGroup, err := check.GetClanGroupAndChack(clanGroupId)
	if err != nil {
		return nil, err
	}
	var timeS, timeE int64
	switch TimeType {
	case "day":
		timeS = time2.GetPcrDayStartTimeToUnix(clanGroup.GameServer)
		timeE = time2.GetPcrDayEndTimeToUnix(clanGroup.GameServer)
	case "yesterday":
		timeS = time2.GetPcrYesterdayStartTimeToUnix(clanGroup.GameServer)
		timeE = time2.GetPcrYesterdayEndTimeToUnix(clanGroup.GameServer)
	case "":
		// 默认所有刀
		timeS = 0
		timeE = time2.GetPcrDayEndTimeToUnix(clanGroup.GameServer)
	default:
		timeS = time2.GetPcrStartTimeToUnixAtStr(clanGroup.GameServer, TimeType)
		timeE = time2.GetPcrEndTimeToUnixAtStr(clanGroup.GameServer, TimeType)
	}

	allSlState, err := gvg_member_extra.GetAllSlState(clanGroup.GvgId, timeS, timeE)
	if err != nil {
		return nil, err
	}
	return allSlState, nil

}

// 公会战信息
type GetClanGvgData struct {
	GroupId   int       `orm:"group_id"      json:"group_id"`
	GroupName string    `orm:"group_name"    json:"group_name"`
	ClanGroup ClanGroup `orm:"clan_group"      json:"clan_group"`
	GvgGroup  GvgGroup  `orm:"gvg_group"      json:"gvg_group"`
}

// Entity is the golang structure for table clan_group.
type ClanGroup struct {
	GroupId      int    `orm:"group_id"      json:"group_id"`      //
	GroupName    string `orm:"group_name"    json:"group_name"`    //
	Privacy      int    `orm:"privacy"       json:"privacy"`       //
	CreatorQqid  int64  `orm:"creator_qqid"  json:"creator_qqid"`  //
	GameServer   string `orm:"game_server"   json:"game_server"`   //
	Notification string `orm:"notification"  json:"notification"`  //
	GvgId        int    `orm:"gvg_id"        json:"gvg_id"`        //
	BindQqGroup  int64  `orm:"bind_qq_group" json:"bind_qq_group"` //
}

// Entity is the golang structure for table gvg_group.
type GvgGroup struct {
	GvgId              int    `orm:"gvg_id"               json:"gvg_id"`               //
	GroupId            int    `orm:"group_id"             json:"group_id"`             //
	CreateQqid         int64  `orm:"create_qqid"          json:"create_qqid"`          //
	GameServer         string `orm:"game_server"          json:"game_server"`          //
	GvgName            string `orm:"gvg_name"             json:"gvg_name"`             //
	BossCycle          int    `orm:"boss_cycle"           json:"boss_cycle"`           //
	BossNum            int    `orm:"boss_num"             json:"boss_num"`             //
	BossHp             int    `orm:"boss_hp"              json:"boss_hp"`              //
	BossFullHp         int    `orm:"boss_full_hp"         json:"boss_full_hp"`         //
	BossLockQqid       int64  `orm:"boss_lock_qqid"       json:"boss_lock_qqid"`       //
	BossLockType       int    `orm:"boss_lock_type"       json:"boss_lock_type"`       //
	BossLockMsg        string `orm:"boss_lock_msg"        json:"boss_lock_msg"`        //
	BossLockTime       int64  `orm:"boss_lock_time"       json:"boss_lock_time"`       //
	ChallengeStratTime int64  `orm:"challenge_strat_time" json:"challenge_strat_time"` //
	ChallengeStratQqid int64  `orm:"challenge_strat_qqid" json:"challenge_strat_qqid"` //
	GvgStartTime       int64  `orm:"gvg_start_time"       json:"gvg_start_time"`       //
	GvgEndTime         int64  `orm:"gvg_end_time"         json:"gvg_end_time"`         //
}

func GetClanGvg(r *ghttp.Request, clanGroupID int) (*GetClanGvgData, error) {
	qqid := user.GetLoginData2(r).Qqid
	clan_member.Login(qqid, clanGroupID, r.GetClientIp())

	getClanGvgData := new(GetClanGvgData)
	clanGroup, err := check.GetClanGroupAndChack(clanGroupID)
	if err != nil {
		return nil, err
	}
	// BossHpCount(getClanGvgData.GroupId, clanGroup.GameServer)
	getClanGvgData.GroupName = clanGroup.GroupName
	getClanGvgData.GroupId = clanGroup.GroupId
	if err := gconv.Struct(clanGroup, &getClanGvgData.ClanGroup); err != nil {
		return nil, err
	}
	gvgGroup, err := gvg_group.GetGvgGroupData(clanGroup.GvgId)
	if err != nil {
		return nil, nil
	}
	if gvgGroup == nil {
		return getClanGvgData, nil
	}
	if err := gconv.Struct(gvgGroup, &getClanGvgData.GvgGroup); err != nil {
		return nil, err
	}
	return getClanGvgData, nil
}

// RemindChallenge 提醒出刀
func RemindChallenge(r *ghttp.Request, clanGroupID int, QQIDlist []int64, Type int) error {
	qqid := user.GetLoginData2(r).Qqid
	clanGroup, err := check.GetClanGroupAndChack(clanGroupID)
	if err != nil {
		return err
	}
	if clanGroup == nil {
		return errors.New("公会不存在")
	}
	if !check.CheckAuthorityGroup(qqid, check.AuthGvgAdmin, clanGroup.GroupId) {
		return errors.New("权限不足")
	}
	var msg string
	qqGroupID := clanGroup.BindQqGroup
	if Type == 2 {
		for _, qqid := range QQIDlist {
			msg += fmt.Sprintf("%s\n", bot.GetAtQQStr(qqid))
		}
		if Type == 0 {
			Type = 2
		}
		clanMember, _ := clan_member.GetClanMember(qqid, clanGroup.GroupId)
		if clanMember != nil {
			msg += fmt.Sprintf("=======\n%s提醒您及时完成今日出刀", clanMember.GameName)
		} else {
			msg += fmt.Sprintf("=======\n%d提醒您及时完成今日出刀", qqid)
		}
		bot.Send(qqGroupID, Type, msg)
		return nil
	}
	userData := bot.GetGroupUserData(qqGroupID, qqid)
	if userData != nil {
		msg = fmt.Sprintf("%s提醒您及时完成今日出刀", userData.NickName)
	} else {
		msg = fmt.Sprintf("%d提醒您及时完成今日出刀", qqid)
	}
	for _, qqid := range QQIDlist {
		bot.Send(qqid, Type, msg)
	}
	return nil
}

type ChangeUserChallengeInput struct {
	ChallengeId     int `v:"min:1#战斗ID不能为空"`
	GroupId         int `v:"min:1#公会ID不能为空"`
	ChallengeDamage int
	BossCycle       string `v:"between:1,100#周期错误"`
	BossNum         string `v:"between:1,5#BOSS序号错误"`
	Meassage        string
}

// 修改成员战斗记录
func ChangeUserChallenge(r *ghttp.Request, data ChangeUserChallengeInput) error {
	qqid := user.GetLoginData2(r).Qqid
	clanGroup, err := check.GetClanGroupAndChack(data.GroupId)
	if err != nil {
		return err
	}
	gvgChallenge, err := gvg_challenge.GetlUserChallengeAtId(data.ChallengeId)
	if err != nil {
		return err
	}
	if gvgChallenge == nil {
		return errors.New("战斗数据不存在")
	}
	if qqid != gvgChallenge.Qqid && !check.CheckAuthorityGroup(qqid, check.AuthGvgAdmin, clanGroup.GroupId) {
		return errors.New("权限不足")
	}
	err = gvg_challenge.ChangeUserChallenge(data.ChallengeId, data.ChallengeDamage, data.BossCycle, data.BossNum, data.Meassage)
	if err != nil {
		return err
	}
	bot.Send(clanGroup.BindQqGroup, 2, check.GetBossStateStr(data.GroupId))
	return nil
}

// 修改成员战斗记录
func DelUserChallenge(r *ghttp.Request, groupID int, challengeId int) error {
	qqid := user.GetLoginData2(r).Qqid
	clanGroup, err := check.GetClanGroupAndChack(groupID)
	if err != nil {
		return err
	}
	gvgChallenge, err := gvg_challenge.GetlUserChallengeAtId(challengeId)
	if err != nil {
		return err
	}
	if qqid != gvgChallenge.Qqid && !check.CheckAuthorityGroup(qqid, check.AuthGvgAdmin, clanGroup.GroupId) {
		return errors.New("权限不足")
	}
	err = gvg_challenge.DelUserChallenge(challengeId)
	if err != nil {
		return err
	}
	bot.Send(clanGroup.BindQqGroup, 2, check.GetBossStateStr(groupID))
	return nil
}
