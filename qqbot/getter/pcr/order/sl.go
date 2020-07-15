package order

import (
	"fmt"
	"github.com/vivid-vvo/vvbot/app/model/gvg_group"
	"github.com/vivid-vvo/vvbot/app/model/gvg_member_extra"
	"github.com/vivid-vvo/vvbot/qqbot/getter"
	"github.com/vivid-vvo/vvbot/qqbot/model/bot"
	"github.com/vivid-vvo/vvbot/qqbot/model/pcr"

	"github.com/vivid-vvo/vvbot/library/Tools"
	"regexp"
	"strings"
)

type Sl struct {
}

func (l Sl) CheckOrder(cm string) bool {
	return true
}

func (l Sl) IsNotCheckOrder() bool {
	return false
}

func (l Sl) GetOrders() []string {
	return []string{
		"^sl( +)?(\\?|？)?$",
		"^查(询)?sl( +)?(\\?|？)?$",
		"^sl查(询)?( +)?(\\?|？)?$",
		"^(撤销|取消)sl$",
		"^sl(取消|撤销)$",
	}
}

func (l Sl) Run(mess getter.MeassageData, cm string, atqq int) {
	var agentQqid, qqid int
	if atqq != 0 {
		agentQqid = int(mess.FromUserID)
		qqid = atqq
	} else {
		qqid = int(mess.FromUserID)
	}
	cm = strings.ToLower(cm)
	regexp1 := regexp.MustCompile(`@.*? `)
	cm = regexp1.ReplaceAllString(cm, "")
	isCancel := Tools.Compare(cm, "^(取消|撤销)|(取消|撤销)$")
	isQuery := Tools.Compare(cm, "[?|？]|查")
	clanGroup, clanMember, err := pcr.GetClanGroupAndUserGroupToCheck(mess.FromGroupID, qqid)
	if err != nil {
		bot.Send(mess.FromGroupID, 2, err.Error())
		return
	}
	_, gvgGroup, err := pcr.GetGvgGroupDataAtGroupIdToCheck(mess.FromGroupID)
	if err != nil {
		bot.Send(mess.FromGroupID, 2, err.Error())
		return
	}
	isSl, err := gvg_member_extra.GetDaySL(qqid, clanGroup.GameServer)
	if err != nil {
		bot.Send(mess.FromGroupID, 2, err.Error())
		return
	}
	if isSl && !isCancel {
		bot.Send(mess.FromGroupID, 2, fmt.Sprintf("%s 今日已SL", clanMember.GameName))
		return
	}
	if isQuery || isCancel && !isSl {
		if !isSl {
			bot.Send(mess.FromGroupID, 2, fmt.Sprintf("%s 今日未SL", clanMember.GameName))
			return
		}
	}
	if isCancel {
		err = gvg_member_extra.CancelDaySL(qqid, clanGroup.GameServer)
		if err != nil {
			bot.Send(mess.FromGroupID, 2, err.Error())
			return
		}
		bot.Send(mess.FromGroupID, 2, fmt.Sprintf("%s 已取消SL", clanMember.GameName))
	} else {
		err = gvg_member_extra.ReportDaySL(qqid, agentQqid, gvgGroup.GvgId)
		if err != nil {
			bot.Send(mess.FromGroupID, 2, err.Error())
			return
		}
		if qqid == gvgGroup.ChallengeStratQqid {
			// 出刀锁定解锁
			if err = gvg_group.CancelChallenge(gvgGroup.GvgId); err != nil {
				bot.Send(mess.FromGroupID, 2, err.Error())
				return
			}
		}
		bot.Send(mess.FromGroupID, 2, fmt.Sprintf("%s 已记录SL", clanMember.GameName))
	}

}
