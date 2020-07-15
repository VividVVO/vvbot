package order

import (
	"fmt"
	"github.com/vivid-vvo/vvbot/app/model/gvg_challenge"
	"github.com/vivid-vvo/vvbot/app/model/user"
	"github.com/vivid-vvo/vvbot/app/service/gvg"
	"github.com/vivid-vvo/vvbot/library/Tools"
	time2 "github.com/vivid-vvo/vvbot/library/time"
	"github.com/vivid-vvo/vvbot/qqbot/getter"
	"github.com/vivid-vvo/vvbot/qqbot/model/bot"
	"github.com/vivid-vvo/vvbot/qqbot/model/pcr"
)

type Chedao struct {
}

func (l Chedao) CheckOrder(cm string) bool {
	return true
}

func (l Chedao) IsNotCheckOrder() bool {
	return false
}

func (l Chedao) GetOrders() []string {
	return []string{
		"^撤刀( +)?(昨)?(天|日)?( +)?$",
		"^取消报刀( +)?(昨)?(天|日)?( +)?$",
		"^报刀取消( +)?(昨)?(天|日)?( +)?$",
	}
}

func (l Chedao) Run(mess getter.MeassageData, cm string, atqq int) {
	qqid := int(mess.FromUserID)
	if atqq != 0 {
		qqid = atqq
		if !user.CheckUserAuthorityGroup(int(mess.FromUserID), user.AuthAdmin) {
			bot.Send(mess.FromGroupID, 2, "权限不足")
			return
		}
	}

	clanGroup, clanMember, err := pcr.GetClanGroupAndUserGroupToCheck(mess.FromGroupID, qqid)
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
	isYesterday := Tools.Compare(cm, "昨")

	var lostChallenge *gvg_challenge.Entity
	// 是否撤销昨日刀
	if isYesterday {
		lostChallenge, err = gvg_challenge.GetYesterdayLostChallenge(qqid, clanGroup.GvgId, clanGroup.GameServer)
	} else {
		lostChallenge, err = gvg_challenge.GetDayLostChallenge(qqid, clanGroup.GvgId, clanGroup.GameServer)
	}
	if err != nil {
		bot.Send(mess.FromGroupID, 2, err.Error())
		return
	}
	if lostChallenge == nil {
		if isYesterday {
			bot.Send(mess.FromGroupID, 2, fmt.Sprintf("%s 撤刀失败,昨日无可撤销刀", clanMember.GameName))
		} else {
			bot.Send(mess.FromGroupID, 2, fmt.Sprintf("%s 撤刀失败,今日无可撤销刀，如果想撤销昨日刀，请输入“撤刀 昨天 [@某人]“ 例：”撤刀 昨天“", clanMember.GameName))
		}
		return
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
		clanMember.GameName, time2.GetTimeAtUnixToZone(clanGroup.GameServer,
			lostChallenge.ChallengeTime).Format("2006-01-02 15:04:05"),
		lostChallenge.BossCycle, lostChallenge.BossNum,
		Tools.NumberFormat(lostChallenge.ChallengeDamage), strEnd))
	bot.Send(mess.FromGroupID, 2, pcr.GetBossStateStr(mess.FromGroupID))
}
