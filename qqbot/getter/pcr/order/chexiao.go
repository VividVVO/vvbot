package order

import (
	"fmt"
	"github.com/vivid-vvo/vvbot/app/model/clan_member"
	"github.com/vivid-vvo/vvbot/app/model/gvg_challenge"
	"github.com/vivid-vvo/vvbot/app/model/user"
	"github.com/vivid-vvo/vvbot/app/service/gvg"
	"github.com/vivid-vvo/vvbot/library/Tools"
	time2 "github.com/vivid-vvo/vvbot/library/time"
	"github.com/vivid-vvo/vvbot/qqbot/getter"
	"github.com/vivid-vvo/vvbot/qqbot/model/bot"
	"github.com/vivid-vvo/vvbot/qqbot/model/pcr"
)

type Chexiao struct {
}

func (l Chexiao) CheckOrder(cm string) bool {
	return true
}

func (l Chexiao) IsNotCheckOrder() bool {
	return false
}

func (l Chexiao) GetOrders() []string {
	return []string{
		"^撤销( +)?$",
	}
}

func (l Chexiao) Run(mess getter.MeassageData, ms string, atqq int) {
	qqid := int(mess.FromUserID)
	clanGroup, err := pcr.GetClanGroupAndChack(mess.FromGroupID)
	if err != nil {
		bot.Send(mess.FromGroupID, 2, err.Error())
		return
	}
	// 检测公会战是否开启
	_, _, err = pcr.GetGvgGroupDataAtGroupIdToCheck(mess.FromGroupID)
	if err != nil {
		bot.Send(mess.FromGroupID, 2, err.Error())
		return
	}
	lostChallenge, err := gvg_challenge.GetLostChallengeAndRepair(clanGroup.GvgId)
	if err != nil {
		bot.Send(mess.FromGroupID, 2, err.Error())
		return
	}
	if lostChallenge == nil {
		bot.Send(mess.FromGroupID, 2, "撤销失败,无可撤销记录")
		return
	}
	if lostChallenge.Qqid != qqid || lostChallenge.RepairType != 0 {
		if !user.CheckUserAuthorityGroup(int(mess.FromUserID), user.AuthAdmin) {
			bot.Send(mess.FromGroupID, 2, "权限不足")
			return
		}
	}
	err = gvg_challenge.BackChallenge(lostChallenge.ChallengeId)
	if err != nil {
		bot.Send(mess.FromGroupID, 2, err.Error())
		return
	}
	// BOSS数据计算刷新
	if err := gvg.BossHpCount(clanGroup.GvgId, clanGroup.GameServer); err != nil {
		bot.Send(mess.FromGroupID, 2, err.Error())
		return
	}
	clanMember, err := clan_member.GetClanMember(lostChallenge.Qqid, clanGroup.GroupId)
	if err != nil {
		bot.Send(mess.FromGroupID, 2, err.Error())
		return
	}
	var gameName string
	if clanMember != nil {
		gameName = clanMember.GameName
	}
	if lostChallenge.RepairType != 0 {
		switch lostChallenge.RepairType {
		case 1:
			bot.Send(mess.FromGroupID, 2, fmt.Sprintf("%s在%s\n修正血量%s已撤销",
				gameName, time2.GetTimeAtUnixToZone(clanGroup.GameServer,
					lostChallenge.ChallengeTime).Format("2006-01-02 15:04:05"),
				Tools.NumberFormat(lostChallenge.RepairHp)))
		case 2:
			bot.Send(mess.FromGroupID, 2, fmt.Sprintf("%s在%s\n修正周目%d已撤销",
				gameName, time2.GetTimeAtUnixToZone(clanGroup.GameServer,
					lostChallenge.ChallengeTime).Format("2006-01-02 15:04:05"), lostChallenge.RepairCycle))
		case 3:
			bot.Send(mess.FromGroupID, 2, fmt.Sprintf("%s在%s\n修正Boss%d已撤销",
				gameName, time2.GetTimeAtUnixToZone(clanGroup.GameServer,
					lostChallenge.ChallengeTime).Format("2006-01-02 15:04:05"), lostChallenge.RepairNum))
		}
	} else {
		var strEnd string
		if lostChallenge.IsSurplus == 1 {
			if lostChallenge.IsContinue == 1 {
				strEnd = "余尾刀"
			} else {
				strEnd = "余刀"
			}
		} else if lostChallenge.IsContinue == 1 {
			strEnd = "收尾刀"
		} else {
			strEnd = "完整刀"
		}
		bot.Send(mess.FromGroupID, 2, fmt.Sprintf("%s在%s\n对%d周目%d号boss\n造成的%s%s伤害已撤销",
			gameName, time2.GetTimeAtUnixToZone(clanGroup.GameServer,
				lostChallenge.ChallengeTime).Format("2006-01-02 15:04:05"),
			lostChallenge.BossCycle, lostChallenge.BossNum,
			Tools.NumberFormat(lostChallenge.ChallengeDamage), strEnd))
	}
	bot.Send(mess.FromGroupID, 2, pcr.GetBossStateStr(mess.FromGroupID))
}
