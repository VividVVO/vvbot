package order

import (
	"fmt"
	"github.com/vivid-vvo/vvbot/app/model/gvg_challenge"
	"github.com/vivid-vvo/vvbot/app/model/gvg_group"
	"github.com/vivid-vvo/vvbot/app/service/gvg"
	"github.com/vivid-vvo/vvbot/qqbot/getter"
	"github.com/vivid-vvo/vvbot/qqbot/model/bot"
	"github.com/vivid-vvo/vvbot/qqbot/model/pcr"

	"github.com/vivid-vvo/vvbot/library/Tools"
	time2 "github.com/vivid-vvo/vvbot/library/time"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Budao struct {
}

func (l Budao) CheckOrder(cm string) bool {
	return true
}

func (l Budao) IsNotCheckOrder() bool {
	return false
}

func (l Budao) GetOrders() []string {
	return []string{
		"^补报( +)?",
		"^补刀( +)?",
		"^补报刀( +)?",
		"^补(尾刀|收尾)( +)?",
		"^补掉刀( +)?",
	}
}

func (l Budao) Run(mess getter.MeassageData, cm string, atqq int) {
	var reg *regexp.Regexp
	isYesterday := Tools.Compare(cm, "昨")
	isContinue := Tools.Compare(cm, "尾")
	isFall := Tools.Compare(cm, "掉")
	cms := strings.Split(cm, " ")
	if len(cms) < 4 {
		bot.Send(mess.FromGroupID, 2, "补刀格式应为：”[补刀|补尾刀|补掉刀] [昨日] 伤害 周目 Boss编号 [@某人] :留言“ 例：“补刀 100w 1 5 :一周目五王的刀”")
		return
	}
	var err error
	var damage, bossCycle, bossNum int
	bossCycle, err = strconv.Atoi(cms[2])
	bossNum, err = strconv.Atoi(cms[3])
	if err != nil || bossCycle < 1 || bossCycle > 100 {
		bot.Send(mess.FromGroupID, 2, "补刀格式应为：”[补刀|补尾刀|补掉刀] 伤害 周目 Boss编号 [昨日] [@某人] :留言“ 例：“补刀 100w 1 5 :一周目五王的刀”")
		return
	}
	if err != nil || bossNum < 1 || bossNum > 5 {
		bot.Send(mess.FromGroupID, 2, "补刀格式应为：”[补刀|补尾刀|补掉刀] 伤害 周目 Boss编号 [昨日] [@某人] :留言“ 例：“补刀 100w 1 5 :一周目五王的刀”")
		return
	}
	if !isContinue && !isFall {
		str2 := strings.Replace(cms[1], " ", "", -1)
		str2 = strings.Replace(str2, ":", "", -1)
		str2 = strings.Replace(str2, "k", "000", -1)
		str2 = strings.Replace(str2, "K", "000", -1)
		str2 = strings.Replace(str2, "千", "000", -1)
		str2 = strings.Replace(str2, "w", "0000", -1)
		str2 = strings.Replace(str2, "W", "0000", -1)
		str2 = strings.Replace(str2, "万", "0000", -1)
		damage, err = strconv.Atoi(str2)
		if err != nil {
			bot.Send(mess.FromGroupID, 2, "补刀格式应为：”[补刀|补尾刀|补掉刀] 伤害 周目 Boss编号 [@某人] :留言“ 例：“补刀 100w 1 5 :一周目五王的刀”")
			return
		}
	}
	reg = regexp.MustCompile(`[:|：](.+)`)
	message := reg.FindString(cm)
	if message != "" {
		message = string([]rune(message)[1:])
	}
	qqid := int(mess.FromUserID)
	qqGroupId := mess.FromGroupID
	var agentQqid, reportQQ int
	if atqq != 0 {
		agentQqid = qqid
		reportQQ = atqq
	} else {
		reportQQ = qqid
	}
	// 检查用户是否已加入当前公会
	_, clanMember, err := pcr.GetClanGroupAndUserGroupToCheck(qqGroupId, reportQQ)
	if err != nil {
		bot.Send(mess.FromGroupID, 2, err.Error())
		return
	}
	clanGroup, gvgData, err := pcr.GetGvgGroupDataAtGroupIdToCheck(qqGroupId)
	if err != nil {
		bot.Send(mess.FromGroupID, 2, err.Error())
		return
	}
	if isContinue {
		damage = gvgData.BossHp
	} else if isFall {
		damage = 0
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
		bot.Send(mess.FromGroupID, 2, err.Error())
		return
	}
	if num >= 3 {
		if isYesterday {
			bot.Send(mess.FromGroupID, 2, fmt.Sprintf("今日 %s 的上报次数已达到3次", clanMember.GameName))
		} else {
			bot.Send(mess.FromGroupID, 2, fmt.Sprintf("昨日 %s 的上报次数已达到3次", clanMember.GameName))
		}
		return
	}
	var challengeTime int
	// 今日开始时间
	dayTimeS := time2.GetPcrDayStartTimeToUnix(clanGroup.GameServer)
	var gvgChallenge2 *gvg_challenge.Entity
	// 刀是否加入到昨日
	if isYesterday {
		// 添加到昨天最后一秒
		challengeTime = int(dayTimeS - 1)
		gvgChallenge2, err = gvg_challenge.GetYesterdayLostChallenge(reportQQ, gvgData.GvgId, clanGroup.GameServer)
	} else {
		challengeTime = int(time.Now().Unix())
		gvgChallenge2, err = gvg_challenge.GetDayLostChallenge(reportQQ, gvgData.GvgId, clanGroup.GameServer)
	}
	if err != nil {
		bot.Send(mess.FromGroupID, 2, err.Error())
		return
	}
	gvgGroup, err := gvg_group.GetGvgGroupData(clanGroup.GvgId)
	if err != nil {
		bot.Send(mess.FromGroupID, 2, err.Error())
		return
	}
	var IsSurplus int
	var strEnd string
	if isFall {
		strEnd = "补掉刀"
	} else if gvgChallenge2 != nil {
		if gvgChallenge2.IsContinue == 0 {
			if isContinue {
				strEnd = "补收尾刀"
			} else {
				strEnd = "补完整刀"
			}
		} else {
			if gvgChallenge2.IsSurplus == 0 {
				if isContinue {
					strEnd = "补余尾刀"
				} else {
					strEnd = "补余刀"
				}
				IsSurplus = 1
			} else {
				if isContinue {
					strEnd = "补余尾刀"
				} else {
					strEnd = "补余刀"
				}
			}
		}
	} else {
		if isContinue {
			strEnd = "补收尾刀"
		} else {
			strEnd = "补完整刀"
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
		BossCycle:       bossCycle,
		BossNum:         bossNum,
	}
	if isContinue {
		gvgChallenge.IsContinue = 1
	}
	err = gvg_challenge.ReportChallenge(gvgChallenge)
	if err != nil {
		bot.Send(mess.FromGroupID, 2, err.Error())
		return
	}
	if reportQQ == gvgGroup.ChallengeStratQqid {
		// 出刀锁定解锁
		if err = gvg_group.CancelChallenge(gvgGroup.GvgId); err != nil {
			bot.Send(mess.FromGroupID, 2, err.Error())
			return
		}
	}
	// BOSS数据计算刷新
	if err := gvg.BossHpCount(clanGroup.GvgId, clanGroup.GameServer); err != nil {
		bot.Send(mess.FromGroupID, 2, err.Error())
		return
	}
	gvgGroup, err = gvg_group.GetGvgGroupData(clanGroup.GvgId)
	if err != nil {
		bot.Send(mess.FromGroupID, 2, err.Error())
		return
	}
	var dayS string
	if isYesterday {
		dayS = "昨日"
	} else {
		dayS = "今日"
	}
	msg := fmt.Sprintf("%s对%d周目，%d号boss\n造成%s点补刀伤害\n(%s第%d刀，%s)\n现在%d周目，%d号boss\n生命值%s",
		clanMember.GameName, bossCycle, bossNum,
		Tools.NumberFormat(damage), dayS, num+1, strEnd, gvgGroup.BossCycle, gvgGroup.BossNum,
		Tools.NumberFormat(gvgGroup.BossHp))
	bot.Send(mess.FromGroupID, 2, msg)

}
