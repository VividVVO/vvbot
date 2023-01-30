package plugin

import (
	"fmt"
	"github.com/vivid-vvo/vvbot/app/model/gvg_challenge"
	"github.com/vivid-vvo/vvbot/app/service/check"
	time2 "github.com/vivid-vvo/vvbot/library/time"
	"github.com/vivid-vvo/vvbot/library/tools"
	"github.com/vivid-vvo/vvbot/qqbot/model/bot"
	"github.com/vivid-vvo/vvbot/qqbot/model/pcr"
	"github.com/vivid-vvo/vvbot/qqbot/plugins"
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
		"^撤(刀|🔪)( +)?(昨)?(天|日)?( +)?$",
		"^取消报(刀|🔪)( +)?(昨)?(天|日)?( +)?$",
		"^报(刀|🔪)取消( +)?(昨)?(天|日)?( +)?$",
		"^返回( +)?",
	}
}

func (l Chedao) Run(mess plugins.MeassageData, cm string, atqq int64) {
	qqid := mess.FromUserID
	// qqGroupId := mess.FromGroupID
	var agentQqid, reportQQ int64
	if atqq != 0 {
		agentQqid = qqid
		reportQQ = atqq
	} else {
		reportQQ = qqid
	}

	clanGroup, clanMember, err := pcr.GetClanGroupAndUserGroupToCheck(mess.FromGroupID, reportQQ)
	if err != nil {
		bot.Send(mess.FromGroupID, mess.SendToType, err.Error())
		return
	}

	// 检测公会战是否开启
	_, _, err = pcr.GetGvgGroupDataAtGroupIdToCheck(mess.FromGroupID)
	if err != nil {
		bot.Send(mess.FromGroupID, mess.SendToType, err.Error())
		return
	}
	isYesterday := tools.Compare(cm, "昨")
	var lostChallenge *gvg_challenge.Entity
	// 是否撤销昨日刀
	if isYesterday {
		lostChallenge, err = gvg_challenge.GetYesterdayLostChallenge(reportQQ, clanGroup.GvgId, clanGroup.GameServer)
	} else {
		lostChallenge, err = gvg_challenge.GetDayLostChallenge(reportQQ, clanGroup.GvgId, clanGroup.GameServer)
	}
	if atqq != 0 && (lostChallenge == nil || lostChallenge.AgentQqid != agentQqid) {
		if !check.CheckAuthorityGroup(mess.FromUserID, check.AuthGvgAdmin, clanGroup.GroupId) {
			bot.Send(mess.FromGroupID, mess.SendToType, "权限不足")
			return
		}
	}
	if err != nil {
		bot.Send(mess.FromGroupID, mess.SendToType, err.Error())
		return
	}
	if lostChallenge == nil {
		if isYesterday {
			bot.Send(mess.FromGroupID, mess.SendToType, fmt.Sprintf("%s 撤刀失败,昨日无可撤销刀", clanMember.GameName))
		} else {
			bot.Send(mess.FromGroupID, mess.SendToType, fmt.Sprintf("%s 撤刀失败,今日无可撤销刀，如果想撤销昨日刀，请输入“撤刀 昨天 [@某人]“ 例：”撤刀 昨天“", clanMember.GameName))
		}
		return
	}
	/*lostChallenge2, err := gvg_challenge.GetLostChallengeAndRepair(clanGroup.GvgId)
	if err != nil {
		bot.Send(mess.FromGroupID, mess.SendToType, err.Error())
		return
	}
	if lostChallenge.IsContinue == 1 && lostChallenge2.ChallengeId != lostChallenge.ChallengeId {
		bot.Send(mess.FromGroupID, mess.SendToType, "撤刀失败,当前已有新刀报告！不能直接撤回尾刀，请使用“撤销”逐步撤回")
		return
	}*/
	err = gvg_challenge.BackChallenge(lostChallenge.ChallengeId)
	if err != nil {
		bot.Send(mess.FromGroupID, mess.SendToType, err.Error())
		return
	}
	// BOSS数据计算刷新
	if err := check.BossHpCount(clanGroup.GvgId, clanGroup.GameServer); err != nil {
		bot.Send(mess.FromGroupID, mess.SendToType, err.Error())
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
	bot.Send(mess.FromGroupID, mess.SendToType, fmt.Sprintf("%s在%s\n对%d周目%d号boss\n造成的%s%s伤害已撤销",
		clanMember.GameName, time2.GetTimeAtUnixToZone(clanGroup.GameServer,
			lostChallenge.ChallengeTime).Format("2006-01-02 15:04:05"),
		lostChallenge.BossCycle, lostChallenge.BossNum,
		tools.NumberFormat(lostChallenge.ChallengeDamage), strEnd))
	bot.Send(mess.FromGroupID, mess.SendToType, pcr.GetBossStateStr(mess.FromGroupID))
}
