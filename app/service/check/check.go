package check

import (
	"errors"
	"fmt"
	"github.com/vivid-vvo/vvbot/app/model/boss_data"
	"github.com/vivid-vvo/vvbot/app/model/clan_group"
	"github.com/vivid-vvo/vvbot/app/model/clan_member"
	"github.com/vivid-vvo/vvbot/app/model/gvg_challenge"
	"github.com/vivid-vvo/vvbot/app/model/gvg_group"
	"github.com/vivid-vvo/vvbot/library/tools"
)

// GetBossStateStr 获取当前BOSS状态文本
func GetBossStateStr(groupId int) string {
	clanGroup, gvgGroup, err := GetGvgGroupDataAtGroupIdToCheck(groupId)
	if err != nil {
		return err.Error()
	}
	BossHpCount(gvgGroup.GvgId, gvgGroup.GameServer)
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

// GetGvgGroupDataAtGroupId 获取公会战信息 通过公会组ID
func GetGvgGroupDataAtGroupIdToCheck(groupID int) (*clan_group.Entity, *gvg_group.Entity, error) {
	clanGroup, err := GetClanGroupAndChack(groupID)
	if err != nil {
		return nil, nil, err
	}
	gvgGroupData, err := gvg_group.GetGvgGroupData(clanGroup.GvgId)
	if err != nil {
		return nil, nil, err
	}
	if gvgGroupData == nil {
		return nil, nil, errors.New("当前未开启公会战，请先进入群内发送“开启公会战 公会战名“开启一场全新的公会战")
	}
	return clanGroup, gvgGroupData, nil
}

// GetClanGroupAndChack 获取公会并检查公会是否存在
func GetClanGroupAndChack(clanGroupID int) (*clan_group.Entity, error) {
	clanGroupData, err := clan_group.GetClanGroup(clanGroupID)
	if err != nil {
		return nil, err
	}
	if clanGroupData == nil {
		return nil, errors.New("公会不存在")
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
		return nil, errors.New("未加入公会，请进入群内发送“加入公会”")
	}
	return userClanGroupData, nil
}

// GetClanGroupAndUserGroupToCheck 获取公户信息和用户信息、并检测是否已加入公会
func GetClanGroupAndUserGroupToCheck(groupId int, qqid int64) (*clan_group.Entity, *clan_member.Entity, error) {
	clanGroupData, err := GetClanGroupAndChack(groupId)
	if err != nil {
		return nil, nil, err
	}
	userClanGroupData, err := GetClanUserGroupToCheck(qqid, clanGroupData.GroupId)
	if err != nil {
		return nil, nil, err
	}
	if userClanGroupData.GroupId != clanGroupData.GroupId {
		return nil, nil, errors.New("您当前未加入本公会")
	}
	return clanGroupData, userClanGroupData, nil
}

// BossHpCount 计算BOSS当前周目和血量  CN 国服, TW 台服, JP 日服, KR 韩服
func BossHpCount(gvgid int, gameServer string) error {
	var cycle, bossNum, nowBossHp int
	bossNum = 1
	cycle = 1
	entity, err := gvg_challenge.GetAllChallengeAndRepair(gvgid)
	if err != nil {
		return err
	}
	var nowBossAllDamage [100][5]int
	for _, v := range entity {
		if v.RepairType == 1 {
			nowBossAllDamage[cycle-1][bossNum-1] += +v.ChallengeDamage
		} else if v.RepairType == 2 {
			cycle = v.RepairCycle
		} else if v.RepairType == 3 {
			bossNum = v.RepairNum
		} else if v.IsContinue == 1 {
			nowBossAllDamage[v.BossCycle-1][v.BossNum-1] += +v.ChallengeDamage
			bossNum++
			if bossNum == 6 {
				bossNum = 1
				cycle++
			}
		} else {
			nowBossAllDamage[v.BossCycle-1][v.BossNum-1] += +v.ChallengeDamage
		}
	}
	gvgBossData := boss_data.GetBossHpList(gameServer, cycle)
	nowBossHp = gvgBossData[bossNum-1] - nowBossAllDamage[cycle-1][bossNum-1]
	if err = gvg_group.UpdateGvgBossData(gvgid, cycle, bossNum, nowBossHp, gvgBossData[bossNum-1]); err != nil {
		return err
	}
	return err
}
