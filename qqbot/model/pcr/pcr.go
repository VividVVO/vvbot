package pcr

import (
	"errors"
	"fmt"
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/util/grand"
	"github.com/vivid-vvo/vvbot/app/model/clan_group"
	"github.com/vivid-vvo/vvbot/app/model/clan_member"
	"github.com/vivid-vvo/vvbot/app/model/gvg_group"
	"github.com/vivid-vvo/vvbot/app/model/user"
	"github.com/vivid-vvo/vvbot/app/service/check"
	"github.com/vivid-vvo/vvbot/library/tools"
)

// GetClanGroupAndChack 获取公会并检查公会是否存在
func GetClanGroupAndChack(qqGroupId int64) (*clan_group.Entity, error) {
	clanGroupData, err := clan_group.GetClanGroupAtQqGroupId(qqGroupId)
	if err != nil {
		return nil, err
	}
	if clanGroupData == nil {
		return nil, errors.New("当前公会不存在，请输入”创建[日台韩国]服公会 公会名“ 例：“创建国服公会 拉胯会长拉胯记”")
	}
	return clanGroupData, nil
}

// GetClanUserGroupToCheck 获取公会成员信息，并检查是否已加入公会
func GetClanUserGroupToCheck(qqid int64, groupID int) (*clan_member.Entity, error) {
	userClanGroupData, err := clan_member.GetClanMember(qqid, groupID)
	if err != nil {
		return nil, err
	}
	if userClanGroupData == nil {
		return nil, errors.New("未加入公会，请先发送“加入公会”")
	}
	return userClanGroupData, nil
}

// GetClanGroupAndUserGroupToCheck 获取公户信息和用户信息、并检测是否已加入公会
func GetClanGroupAndUserGroupToCheck(qqGroupId int64, qqid int64) (*clan_group.Entity, *clan_member.Entity, error) {
	clanGroupData, err := GetClanGroupAndChack(qqGroupId)
	if err != nil {
		return nil, nil, err
	}
	userClanGroupData, err := GetClanUserGroupToCheck(qqid, clanGroupData.GroupId)
	if err != nil {
		return nil, nil, err
	}

	if userClanGroupData.GroupId != clanGroupData.GroupId {
		return nil, nil, errors.New("您当前未加入本公会，请发送“加入公会”")
	}
	return clanGroupData, userClanGroupData, nil
}

// GetGvgGroupDataAtGroupIdToCheck 获取公会信息和公会战信息，并检查公会是否存在和公会战是否已开启
func GetGvgGroupDataAtGroupIdToCheck(qqGroupId int64) (*clan_group.Entity, *gvg_group.Entity, error) {
	clanGroup, err := GetClanGroupAndChack(qqGroupId)
	if err != nil {
		return nil, nil, err
	}
	gvgGroupData, err := gvg_group.GetGvgGroupData(clanGroup.GvgId)
	if err != nil {
		return nil, nil, err
	}
	if gvgGroupData == nil {
		return nil, nil, errors.New("当前未开启公会战，请先使用“开启公会战 公会战名“开启一场全新的公会战")
	}
	return clanGroup, gvgGroupData, nil
}

// GetBossStateStr 获取当前BOSS状态文本
func GetBossStateStr(qqGroupId int64) string {
	clanGroup, gvgGroup, err := GetGvgGroupDataAtGroupIdToCheck(qqGroupId)
	if err != nil {
		return err.Error()
	}
	check.BossHpCount(gvgGroup.GvgId, gvgGroup.GameServer)
	gvgGroup, err = gvg_group.GetGvgGroupData(gvgGroup.GvgId)
	if err != nil {
		return err.Error()
	}
	msg := fmt.Sprintf("公会：%s\n当前公会战：%s\n现在%d周目，%d号boss\n生命值%s", clanGroup.GroupName, gvgGroup.GvgName, gvgGroup.BossCycle, gvgGroup.BossNum, tools.NumberFormat(gvgGroup.BossHp))
	if gvgGroup.ChallengeStratQqid != 0 {
		member, err := clan_member.GetClanMember(gvgGroup.ChallengeStratQqid, clanGroup.GroupId)
		if err != nil {
			return err.Error()
		}
		msg = fmt.Sprintf("%s\n%s正在挑战boss", msg, member.GameName)
	}
	if gvgGroup.BossLockType == 1 {
		member, err := clan_member.GetClanMember(gvgGroup.BossLockQqid, clanGroup.GroupId)
		if err != nil {
			return err.Error()
		}
		msg = fmt.Sprintf("%s\n%s锁定了boss\n留言：%s", msg, member.GameName, gvgGroup.BossLockMsg)
	}
	return msg
}

func UserJoinGroup(qqGroupId int64, qqid int64, qqNickName string) error {
	clanGroupData, err := GetClanGroupAndChack(qqGroupId)
	if err != nil {
		return err
	}
	userInfo, err := user.GetProfile(qqid)
	if err != nil {
		return err
	}
	if userInfo == nil {
		// 注册
		err := user.SignUp(qqid, grand.Letters(12), qqNickName, "")
		if err != nil {
			return err
		}
	}
	clanMember, err := clan_member.GetClanMember(qqid, clanGroupData.GroupId)
	if err != nil {
		return err
	}
	if clanMember != nil {
		return errors.New("成员已加入当前公会")
	}
	err = clan_member.MemberExitGroupAtQqid(qqid, clanGroupData.GroupId)
	if err != nil {
		return err
	}
	entity := new(clan_member.Entity)
	entity.GroupId = clanGroupData.GroupId
	entity.Qqid = qqid
	entity.GameName = qqNickName
	entity.JoinTime = gtime.Now().Unix()
	err = clan_member.JoinClan(entity)
	if err != nil {
		return err
	}
	if err := user.ChangeClanGroupId(qqid, entity.GroupId); err != nil {
		return err
	}
	return nil
}
