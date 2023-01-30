package plugin

import (
	"fmt"
	"github.com/vivid-vvo/vvbot/app/model/clan_member"
	"github.com/vivid-vvo/vvbot/app/model/gvg_challenge"
	"github.com/vivid-vvo/vvbot/app/model/gvg_group"
	"github.com/vivid-vvo/vvbot/app/model/gvg_member_extra"
	"github.com/vivid-vvo/vvbot/app/service/check"
	"github.com/vivid-vvo/vvbot/qqbot/model/bot"
	"github.com/vivid-vvo/vvbot/qqbot/model/pcr"
	"github.com/vivid-vvo/vvbot/qqbot/plugins"

	time2 "github.com/vivid-vvo/vvbot/library/time"
	"github.com/vivid-vvo/vvbot/library/tools"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Baodao struct {
}

func (l Baodao) CheckOrder(cm string) bool {
	return true
}

func (l Baodao) IsNotCheckOrder() bool {
	return false
}

func (l Baodao) GetOrders() []string {
	return []string{
		"^报(刀|🔪)",
		"^(尾(刀|🔪)|收尾|出补时(刀|🔪))( +)?$",
		"^掉(刀|🔪)( +)?$\"",
		"^出(刀|🔪)",
	}
}

func (l Baodao) Run(mess plugins.MeassageData, cm string, atqq int64) {
	var reg *regexp.Regexp
	isYesterday := tools.Compare(cm, "昨")
	isContinue := tools.Compare(cm, "(尾|补)")
	isFall := tools.Compare(cm, "掉")
	var damage int
	var err error
	if !isContinue && !isFall {
		reg = regexp.MustCompile(`(\d+)(.+)[:|：]`)
		str2 := reg.FindString(cm)
		if str2 == "" {
			reg = regexp.MustCompile(`(\d+)(k|w|K|W|千|万)?`)
			str2 = reg.FindString(cm)
		}
		str2 = strings.Replace(str2, " ", "", -1)
		str2 = strings.Replace(str2, ":", "", -1)
		str2 = strings.Replace(str2, "：", "", -1)
		str2 = strings.Replace(str2, "k", "000", -1)
		str2 = strings.Replace(str2, "K", "000", -1)
		str2 = strings.Replace(str2, "千", "000", -1)
		str2 = strings.Replace(str2, "w", "0000", -1)
		str2 = strings.Replace(str2, "W", "0000", -1)
		str2 = strings.Replace(str2, "万", "0000", -1)
		damage, err = strconv.Atoi(str2)
		if err != nil {
			bot.Send(mess.FromGroupID, mess.SendToType, "报刀格式错误，正确格式为“报刀 伤害 [昨天] [@某人] [:留言]” 例：“报刀100 @会长 :甜心一刀”")
			return
		}
	}
	reg = regexp.MustCompile(`[:|：](.+)`)
	message := reg.FindString(cm)
	if message != "" {
		message = string([]rune(message)[1:])
	}
	qqid := mess.FromUserID
	qqGroupId := mess.FromGroupID

	var agentQqid, reportQQ int64
	if atqq != 0 {
		agentQqid = qqid
		reportQQ = atqq
	} else {
		reportQQ = qqid
	}
	// 检查用户是否已加入当前公会
	_, clanMember, err := pcr.GetClanGroupAndUserGroupToCheck(qqGroupId, reportQQ)
	if err != nil {
		bot.Send(mess.FromGroupID, mess.SendToType, err.Error())
		return
	}
	clanGroup, gvgData, err := pcr.GetGvgGroupDataAtGroupIdToCheck(qqGroupId)
	if err != nil {
		bot.Send(mess.FromGroupID, mess.SendToType, err.Error())
		return
	}
	if isContinue {
		damage = gvgData.BossHp
	} else if isFall {
		damage = 0
	} else {
		if damage >= gvgData.BossHp {
			bot.Send(mess.FromGroupID, mess.SendToType, "报刀无效，伤害量必须小于剩余血量，如果击败boss请发送“尾刀”")
			return
		}
	}
	var timeS, timeE int64

	// 是否是昨天的刀
	if isYesterday {
		// 每天5点刷新
		timeS = time2.GetPcrYesterdayStartTimeToUnix(clanGroup.GameServer)
		timeE = time2.GetPcrYesterdayEndTimeToUnix(clanGroup.GameServer)
	} else {
		timeS = time2.GetPcrDayStartTimeToUnix(clanGroup.GameServer)
		timeE = time2.GetPcrDayEndTimeToUnix(clanGroup.GameServer)
	}
	// 获取指定时间内的刀，过滤尾刀
	num, err := gvg_challenge.GetChallengeAtTime(reportQQ, gvgData.GvgId, timeS, timeE)
	if err != nil {
		bot.Send(mess.FromGroupID, mess.SendToType, err.Error())
		return
	}
	if num >= 3 {
		if isYesterday {
			bot.Send(mess.FromGroupID, mess.SendToType, fmt.Sprintf(" %s 昨日的上报次数已达到3次", clanMember.GameName))
		} else {
			bot.Send(mess.FromGroupID, mess.SendToType, fmt.Sprintf(" %s 今日的上报次数已达到3次", clanMember.GameName))
		}
		return
	}
	var challengeTime int64
	// 今日开始时间
	dayTimeS := time2.GetPcrDayStartTimeToUnix(clanGroup.GameServer)
	var gvgChallenge2 *gvg_challenge.Entity
	// 刀是否加入到昨日
	if isYesterday {
		// 添加到昨天最后一秒
		challengeTime = dayTimeS - 1
		gvgChallenge2, err = gvg_challenge.GetYesterdayLostChallenge(reportQQ, gvgData.GvgId, clanGroup.GameServer)
	} else {
		challengeTime = time.Now().Unix()
		gvgChallenge2, err = gvg_challenge.GetDayLostChallenge(reportQQ, gvgData.GvgId, clanGroup.GameServer)
	}
	if err != nil {
		bot.Send(mess.FromGroupID, mess.SendToType, err.Error())
		return
	}
	gvgGroup, err := gvg_group.GetGvgGroupData(clanGroup.GvgId)
	if err != nil {
		bot.Send(mess.FromGroupID, mess.SendToType, err.Error())
		return
	}
	var IsSurplus int
	var strEnd string
	if isFall {
		strEnd = "掉刀"
	} else if gvgChallenge2 != nil {
		if gvgChallenge2.IsContinue == 0 {
			if isContinue {
				strEnd = "收尾刀"
			} else {
				strEnd = "完整刀"
			}
		} else {
			if gvgChallenge2.IsSurplus == 0 {
				if isContinue {
					strEnd = "余尾刀"
				} else {
					strEnd = "余刀"
				}
				IsSurplus = 1
			} else {
				if isContinue {
					strEnd = "收尾刀"
				} else {
					strEnd = "完整刀"
				}
			}
		}
	} else {
		if isContinue {
			strEnd = "收尾刀"
		} else {
			strEnd = "完整刀"
		}
	}
	gvgChallenge := gvg_challenge.Entity{
		Qqid:            reportQQ,
		GvgId:           clanGroup.GvgId,
		ClanGroupId:     clanGroup.GroupId,
		AgentQqid:       agentQqid,
		ChallengeDamage: damage,
		Message:         message,
		IsSurplus:       IsSurplus,
		ChallengeTime:   challengeTime,
		BossCycle:       gvgGroup.BossCycle,
		BossNum:         gvgGroup.BossNum,
		BossHp:          gvgGroup.BossHp,
	}
	if isContinue {
		gvgChallenge.IsContinue = 1
	}
	err = gvg_challenge.ReportChallenge(gvgChallenge)
	if err != nil {
		bot.Send(mess.FromGroupID, mess.SendToType, err.Error())
		return
	}
	if reportQQ == gvgGroup.ChallengeStratQqid {
		// 出刀锁定解锁
		if err = gvg_group.CancelChallenge(gvgGroup.GvgId); err != nil {
			bot.Send(mess.FromGroupID, mess.SendToType, err.Error())
			return
		}
	}
	// BOSS数据计算刷新
	if err := check.BossHpCount(clanGroup.GvgId, clanGroup.GameServer); err != nil {
		bot.Send(mess.FromGroupID, mess.SendToType, err.Error())
		return
	}
	gvgGroup, err = gvg_group.GetGvgGroupData(clanGroup.GvgId)
	if err != nil {
		bot.Send(mess.FromGroupID, mess.SendToType, err.Error())
		return
	}
	var dayS string
	if isYesterday {
		dayS = "昨日"
	} else {
		dayS = "今日"
	}
	msg := fmt.Sprintf("%s对boss造成%s点伤害\n(%s第%d刀，%s)\n现在%d周目，%d号boss\n生命值%s",
		clanMember.GameName, tools.NumberFormat(damage), dayS, num+1, strEnd, gvgGroup.BossCycle, gvgGroup.BossNum,
		tools.NumberFormat(gvgGroup.BossHp))
	gvg_member_extra.ReportDownTree(reportQQ, clanGroup.GvgId, 1)
	bot.Send(mess.FromGroupID, mess.SendToType, msg)
	if isContinue {
		upTreeList, err := gvg_member_extra.GetAllUpTree(clanGroup.GvgId)
		if err != nil {
			bot.Send(mess.FromGroupID, mess.SendToType, err.Error())
			return
		}
		if len(upTreeList) == 0 {
			return
		}
		var msg, errMsg string
		for _, upTree := range upTreeList {
			clanMember, err := clan_member.GetClanMember(upTree.Qqid, clanGroup.GroupId)
			if err != nil {
				errMsg += fmt.Sprintf("\nqq: %d（%s）, 下树失败 -> %s", upTree.Qqid, clanMember.GameName, err.Error())
				continue
			}
			err = gvg_member_extra.ReportDownTree(upTree.Qqid, clanGroup.GvgId, 2)
			if err != nil {
				errMsg += fmt.Sprintf("\nqq: %d（%s）, 下树失败 -> %s", upTree.Qqid, clanMember.GameName, err.Error())
				continue
			}
			msg += fmt.Sprintf("\n%s %s", clanMember.GameName, bot.GetAtQQStr(upTree.Qqid))
		}
		bot.Send(mess.FromGroupID, mess.SendToType, fmt.Sprintf("BOSS已被击败%s", errMsg+msg))
	}
}
