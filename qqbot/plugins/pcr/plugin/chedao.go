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
		"^ζ€(ε|πͺ)( +)?(ζ¨)?(ε€©|ζ₯)?( +)?$",
		"^εζΆζ₯(ε|πͺ)( +)?(ζ¨)?(ε€©|ζ₯)?( +)?$",
		"^ζ₯(ε|πͺ)εζΆ( +)?(ζ¨)?(ε€©|ζ₯)?( +)?$",
		"^θΏε( +)?",
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

	// ζ£ζ΅ε¬δΌζζ―ε¦εΌε―
	_, _, err = pcr.GetGvgGroupDataAtGroupIdToCheck(mess.FromGroupID)
	if err != nil {
		bot.Send(mess.FromGroupID, mess.SendToType, err.Error())
		return
	}
	isYesterday := tools.Compare(cm, "ζ¨")
	var lostChallenge *gvg_challenge.Entity
	// ζ―ε¦ζ€ιζ¨ζ₯ε
	if isYesterday {
		lostChallenge, err = gvg_challenge.GetYesterdayLostChallenge(reportQQ, clanGroup.GvgId, clanGroup.GameServer)
	} else {
		lostChallenge, err = gvg_challenge.GetDayLostChallenge(reportQQ, clanGroup.GvgId, clanGroup.GameServer)
	}
	if atqq != 0 && (lostChallenge == nil || lostChallenge.AgentQqid != agentQqid) {
		if !check.CheckAuthorityGroup(mess.FromUserID, check.AuthGvgAdmin, clanGroup.GroupId) {
			bot.Send(mess.FromGroupID, mess.SendToType, "ζιδΈθΆ³")
			return
		}
	}
	if err != nil {
		bot.Send(mess.FromGroupID, mess.SendToType, err.Error())
		return
	}
	if lostChallenge == nil {
		if isYesterday {
			bot.Send(mess.FromGroupID, mess.SendToType, fmt.Sprintf("%s ζ€εε€±θ΄₯,ζ¨ζ₯ζ ε―ζ€ιε", clanMember.GameName))
		} else {
			bot.Send(mess.FromGroupID, mess.SendToType, fmt.Sprintf("%s ζ€εε€±θ΄₯,δ»ζ₯ζ ε―ζ€ιεοΌε¦ζζ³ζ€ιζ¨ζ₯εοΌθ―·θΎε₯βζ€ε ζ¨ε€© [@ζδΊΊ]β δΎοΌβζ€ε ζ¨ε€©β", clanMember.GameName))
		}
		return
	}
	/*lostChallenge2, err := gvg_challenge.GetLostChallengeAndRepair(clanGroup.GvgId)
	if err != nil {
		bot.Send(mess.FromGroupID, mess.SendToType, err.Error())
		return
	}
	if lostChallenge.IsContinue == 1 && lostChallenge2.ChallengeId != lostChallenge.ChallengeId {
		bot.Send(mess.FromGroupID, mess.SendToType, "ζ€εε€±θ΄₯,ε½εε·²ζζ°εζ₯εοΌδΈθ½η΄ζ₯ζ€εε°ΎεοΌθ―·δ½Ώη¨βζ€ιβιζ­₯ζ€ε")
		return
	}*/
	err = gvg_challenge.BackChallenge(lostChallenge.ChallengeId)
	if err != nil {
		bot.Send(mess.FromGroupID, mess.SendToType, err.Error())
		return
	}
	// BOSSζ°ζ?θ?‘η?ε·ζ°
	if err := check.BossHpCount(clanGroup.GvgId, clanGroup.GameServer); err != nil {
		bot.Send(mess.FromGroupID, mess.SendToType, err.Error())
		return
	}
	var strEnd string
	if lostChallenge.IsSurplus == 1 {
		if lostChallenge.IsContinue == 1 {
			strEnd = "δ½ε°Ύε"
		} else {
			strEnd = "δ½ε"
		}
	} else if lostChallenge.IsContinue == 1 {
		strEnd = "ζΆε°Ύε"
	} else {
		strEnd = "ε?ζ΄ε"
	}
	bot.Send(mess.FromGroupID, mess.SendToType, fmt.Sprintf("%sε¨%s\nε―Ή%dε¨η?%dε·boss\nι ζη%s%sδΌ€ε?³ε·²ζ€ι",
		clanMember.GameName, time2.GetTimeAtUnixToZone(clanGroup.GameServer,
			lostChallenge.ChallengeTime).Format("2006-01-02 15:04:05"),
		lostChallenge.BossCycle, lostChallenge.BossNum,
		tools.NumberFormat(lostChallenge.ChallengeDamage), strEnd))
	bot.Send(mess.FromGroupID, mess.SendToType, pcr.GetBossStateStr(mess.FromGroupID))
}
