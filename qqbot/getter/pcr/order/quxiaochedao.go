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

type Quxiaochedao struct {
}

func (l Quxiaochedao) CheckOrder(cm string) bool {
	return true
}

func (l Quxiaochedao) IsNotCheckOrder() bool {
	return false
}

func (l Quxiaochedao) GetOrders() []string {
	return []string{
		"^取消撤刀( +)?(昨)?(天|日)?( +)?$",
		"^撤刀取消( +)?(昨)?(天|日)?( +)?$",
	}
}

func (l Quxiaochedao) Run(mess getter.MeassageData, cm string, atqq int) {
	qqid := int(mess.FromUserID)
	if atqq != 0 {
		qqid = atqq
		if !user.CheckUserAuthorityGroup(int(mess.FromUserID), user.AuthAdmin) {
			bot.Send(mess.FromGroupID, 2, "权限不足")
			return
		}
	}
	isYesterday := Tools.Compare(cm, "昨")
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
	var lostChallenge *gvg_challenge.Entity
	if isYesterday {
		lostChallenge, err = gvg_challenge.GetYesterdayLostChallengeToBack(qqid, clanGroup.GvgId, clanGroup.GameServer)
	} else {
		lostChallenge, err = gvg_challenge.GetDayLostChallengeToBack(qqid, clanGroup.GvgId, clanGroup.GameServer)
	}
	if err != nil {
		bot.Send(mess.FromGroupID, 2, err.Error())
		return
	}
	if lostChallenge == nil {
		if isYesterday {
			bot.Send(mess.FromGroupID, 2, "取消撤刀失败,昨日无出刀记录")
		} else {
			bot.Send(mess.FromGroupID, 2, "取消撤刀失败,今日无出刀记录")
		}
		return
	}
	if lostChallenge.IsDelete == 0 {
		bot.Send(mess.FromGroupID, 2, "取消撤刀失败,最近一刀未撤刀")
		return
	}
	err = gvg_challenge.CancelBackChallenge(lostChallenge.ChallengeId)
	if err != nil {
		bot.Send(mess.FromGroupID, 2, err.Error())
		return
	}
	// BOSS数据计算刷新
	if err := gvg.BossHpCount(clanGroup.GvgId, clanGroup.GameServer); err != nil {
		bot.Send(mess.FromGroupID, 2, err.Error())
		return
	}
	bot.Send(mess.FromGroupID, 2, fmt.Sprintf("%s在%s\n对%d周目%d号boss\n造成的%s伤害已恢复",
		clanMember.GameName, time2.GetTimeAtUnixToZone(clanGroup.GameServer,
			lostChallenge.ChallengeTime).Format("2006-01-02 15:04:05"),
		lostChallenge.BossCycle, lostChallenge.BossNum, Tools.NumberFormat(lostChallenge.ChallengeDamage)))
	bot.Send(mess.FromGroupID, 2, pcr.GetBossStateStr(mess.FromGroupID))

}
