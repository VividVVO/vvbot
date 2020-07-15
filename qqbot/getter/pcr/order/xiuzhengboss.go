package order

import (
	"github.com/vivid-vvo/vvbot/app/model/gvg_challenge"
	"github.com/vivid-vvo/vvbot/app/service/gvg"
	"github.com/vivid-vvo/vvbot/qqbot/getter"
	"github.com/vivid-vvo/vvbot/qqbot/model/bot"
	"github.com/vivid-vvo/vvbot/qqbot/model/pcr"

	"github.com/vivid-vvo/vvbot/library/Tools"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Xiuzhengboss struct {
}

func (l Xiuzhengboss) CheckOrder(cm string) bool {
	return true
}

func (l Xiuzhengboss) IsNotCheckOrder() bool {
	return false
}

func (l Xiuzhengboss) GetOrders() []string {
	return []string{
		"^修正血量",
		"^修正周目",
		"^修正boss",
	}
}

func (l Xiuzhengboss) Run(mess getter.MeassageData, cm string, atqq int) {
	qqid := int(mess.FromUserID)
	clanData, gvgData, err := pcr.GetGvgGroupDataAtGroupIdToCheck(mess.FromGroupID)
	if err != nil {
		bot.Send(mess.FromGroupID, 2, err.Error())
		return
	}
	isBossHp := Tools.Compare(cm, "修正血量")
	isBoosCycle := Tools.Compare(cm, "修正周目")
	isBoosNum := Tools.Compare(cm, "修正boss")
	var repairType, damage, bossCycle, bossNum, bossHp int
	if isBossHp {
		reg := regexp.MustCompile(`(\d+)(k|w|K|W|千|万)?`)
		str2 := reg.FindString(cm)
		str2 = strings.Replace(str2, " ", "", -1)
		str2 = strings.Replace(str2, ":", "", -1)
		str2 = strings.Replace(str2, "k", "000", -1)
		str2 = strings.Replace(str2, "K", "000", -1)
		str2 = strings.Replace(str2, "千", "000", -1)
		str2 = strings.Replace(str2, "w", "0000", -1)
		str2 = strings.Replace(str2, "W", "0000", -1)
		str2 = strings.Replace(str2, "万", "0000", -1)

		bossHp, err = strconv.Atoi(str2)
		if err != nil {
			bot.Send(mess.FromGroupID, 2, "修正血量格式错误，正确格式为“修正血量 血量” 例：”修正血量 1000000“")
			return
		}
		damage = -(bossHp - gvgData.BossHp)
		repairType = 1
	} else if isBoosCycle {
		reg := regexp.MustCompile(`[0-9]{1,4}`)
		str2 := reg.FindString(cm)
		bossCycle, err = strconv.Atoi(str2)
		if err != nil || bossCycle < 1 || bossCycle > 100 {
			bot.Send(mess.FromGroupID, 2, "修正周目格式错误，正确格式为“修正周目 周目” 例：”修正周目 2“")
			return
		}
		repairType = 2
	} else if isBoosNum {
		reg := regexp.MustCompile(`[0-9]{1,2}`)
		str2 := reg.FindString(cm)
		bossNum, err = strconv.Atoi(str2)
		if err != nil || bossNum < 1 || bossNum > 5 {
			bot.Send(mess.FromGroupID, 2, "修正boss格式错误，正确格式为“修正boss 第几个” 例：”修正boss 5“")
			return
		}
		repairType = 3
	} else {
		return
	}
	gvgChallenge := gvg_challenge.Entity{
		GvgId:           clanData.GvgId,
		ClanGroupId:     clanData.GroupId,
		Qqid:            qqid,
		ChallengeDamage: damage,
		RepairType:      repairType,
		RepairCycle:     bossCycle,
		RepairHp:        bossHp,
		RepairNum:       bossNum,
		IsContinue:      0,
		Message:         "修正",
		ChallengeTime:   int(time.Now().Unix()),
	}
	if err := gvg_challenge.ReportChallenge(gvgChallenge); err != nil {
		bot.Send(mess.FromGroupID, 2, err.Error())
		return
	}
	// BOSS数据计算刷新
	if err := gvg.BossHpCount(clanData.GvgId, clanData.GameServer); err != nil {
		bot.Send(mess.FromGroupID, 2, err.Error())
		return
	}
	bot.Send(mess.FromGroupID, 2, "修正完毕")
	bot.Send(mess.FromGroupID, 2, pcr.GetBossStateStr(mess.FromGroupID))
}
