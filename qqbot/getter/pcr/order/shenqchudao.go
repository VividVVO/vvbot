package order

import (
	"fmt"
	"github.com/vivid-vvo/vvbot/app/model/clan_member"
	"github.com/vivid-vvo/vvbot/app/model/gvg_challenge"
	"github.com/vivid-vvo/vvbot/app/model/gvg_group"
	"github.com/vivid-vvo/vvbot/qqbot/getter"
	"github.com/vivid-vvo/vvbot/qqbot/model/bot"
	"github.com/vivid-vvo/vvbot/qqbot/model/pcr"

	time2 "github.com/vivid-vvo/vvbot/library/time"
	"time"
)

type Shenqchudao struct {
}

func (l Shenqchudao) CheckOrder(cm string) bool {
	return true
}

func (l Shenqchudao) IsNotCheckOrder() bool {
	return false
}

func (l Shenqchudao) GetOrders() []string {
	return []string{
		"^申请出刀( +)?$",
		"^出刀申请( +)?$",
		"^申请出站( +)?$",
		"^出站申请( +)?$",
		"^申请出柜( +)?$"}
}

func (l Shenqchudao) Run(mess getter.MeassageData, cm string, atqq int) {
	qqid := int(mess.FromUserID)
	if atqq != 0 {
		qqid = atqq
	}
	_, clanMember, err := pcr.GetClanGroupAndUserGroupToCheck(mess.FromGroupID, qqid)
	if err != nil {
		bot.Send(mess.FromGroupID, 2, err.Error())
		return
	}
	clanGroup, gvgGroup, err := pcr.GetGvgGroupDataAtGroupIdToCheck(mess.FromGroupID)
	if err != nil {
		bot.Send(mess.FromGroupID, 2, err.Error())
		return
	}

	// 跨日，出刀申请自动清除
	if int(time2.GetPcrDayStartTimeToUnix(clanGroup.GameServer)) > gvgGroup.ChallengeStratTime {
		gvgGroup.ChallengeStratQqid = 0
	}
	if gvgGroup.ChallengeStratQqid == qqid {
		bot.Send(mess.FromGroupID, 2, fmt.Sprintf("%s 您正在挑战BOSS", clanMember.GameName))
		return
	}
	// 获取今日刀，过滤尾刀
	num, err := gvg_challenge.GetDayChallenge(qqid, gvgGroup.GvgId, clanGroup.GameServer)
	if err != nil {
		bot.Send(mess.FromGroupID, 2, err.Error())
		return
	}
	if num >= 3 {
		bot.Send(mess.FromGroupID, 2, fmt.Sprintf("今日 %s 的出刀次数已达到3次", clanMember.GameName))
		return
	}

	if gvgGroup.BossLockType == 1 {
		clanMember, err := clan_member.GetClanMember(gvgGroup.BossLockQqid, clanGroup.GroupId)
		if err != nil {
			bot.Send(mess.FromGroupID, 2, err.Error())
			return
		}
		var gameName string
		if clanMember != nil {
			gameName = clanMember.GameName
		}

		t := int(time.Now().Unix()) - gvgGroup.BossLockTime
		msg := fmt.Sprintf("申请失败，%s在%d分%d秒前锁定了boss", gameName, t/60, t%60)
		if gvgGroup.BossLockMsg != "" {
			msg += "\n留言：" + gvgGroup.BossLockMsg
		}
		bot.Send(mess.FromGroupID, 2, msg)
		return
	}
	if gvgGroup.ChallengeStratQqid != 0 {
		clanMember, err := clan_member.GetClanMember(gvgGroup.ChallengeStratQqid, clanGroup.GroupId)
		if err != nil {
			bot.Send(mess.FromGroupID, 2, err.Error())
			return
		}
		var gameName string
		if clanMember != nil {
			gameName = clanMember.GameName
		}
		bot.Send(mess.FromGroupID, 2, fmt.Sprintf("申请失败，%s 正在挑战boss", gameName))
		return
	}
	err = gvg_group.ApplyChallenge(qqid, gvgGroup.GvgId)
	if err != nil {
		bot.Send(mess.FromGroupID, 2, err.Error())
		return
	}
	clanMember, err = clan_member.GetClanMember(qqid, clanGroup.GroupId)
	if err != nil {
		bot.Send(mess.FromGroupID, 2, err.Error())
		return
	}
	var gameName string
	if clanMember != nil {
		gameName = clanMember.GameName
	}
	if err != nil {
		bot.Send(mess.FromGroupID, 2, err.Error())
		return
	}

	msg := fmt.Sprintf("%s 已开始挑战boss\n%s", gameName,
		pcr.GetBossStateStr(mess.FromGroupID))
	bot.Send(mess.FromGroupID, 2, msg)
}
