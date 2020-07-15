package order

import (
	"fmt"
	"github.com/vivid-vvo/vvbot/app/model/gvg_group"
	"github.com/vivid-vvo/vvbot/app/model/gvg_member_extra"
	"github.com/vivid-vvo/vvbot/library/Tools"
	"github.com/vivid-vvo/vvbot/qqbot/getter"
	"github.com/vivid-vvo/vvbot/qqbot/model/bot"
	"github.com/vivid-vvo/vvbot/qqbot/model/pcr"
	"strings"
)

type Guashu struct {
}

func (l Guashu) CheckOrder(cm string) bool {
	return true
}

func (l Guashu) IsNotCheckOrder() bool {
	return false
}

func (l Guashu) GetOrders() []string {
	return []string{
		"^(挂树|上树)( +)?(\\?|？)?$",
		"^(挂树|上树)(取消|撤销)( +)?(\\?|？)?$",
		"^(取消|撤销)(挂树|上树)( +)?(\\?|？)?$",
		"^下树( +)?$",
		"^我挂树了( +)?$",
	}
}

func (l Guashu) Run(mess getter.MeassageData, cm string, atqq int) {
	var agentQqid, qqid int
	if atqq != 0 {
		agentQqid = int(mess.FromUserID)
		qqid = atqq
	} else {
		qqid = int(mess.FromUserID)
	}
	cm = strings.ToLower(cm)

	isCancel := Tools.Compare(cm, "^(取消|撤销|下树)|(取消|撤销)$")
	isQuery := Tools.Compare(cm, "[?|？]")
	_, clanMember, err := pcr.GetClanGroupAndUserGroupToCheck(mess.FromGroupID, qqid)
	if err != nil {
		bot.Send(mess.FromGroupID, 2, err.Error())
		return
	}
	_, gvgGroup, err := pcr.GetGvgGroupDataAtGroupIdToCheck(mess.FromGroupID)
	if err != nil {
		bot.Send(mess.FromGroupID, 2, err.Error())
		return
	}
	isTree, err := gvg_member_extra.GetIsUpTree(qqid)
	if err != nil {
		bot.Send(mess.FromGroupID, 2, err.Error())
		return
	}
	if isTree && !isCancel {
		bot.Send(mess.FromGroupID, 2, fmt.Sprintf("%s 您已在树上", clanMember.GameName))
		return
	}
	if isQuery || isCancel && !isTree {

		bot.Send(mess.FromGroupID, 2, fmt.Sprintf("%s 您未在树上", clanMember.GameName))
		return

	}
	if isCancel {
		err = gvg_member_extra.ReportDownTree(qqid, 0, gvgGroup.GvgId, 0)
		if err != nil {
			bot.Send(mess.FromGroupID, 2, err.Error())
			return
		}
		bot.Send(mess.FromGroupID, 2, fmt.Sprintf("%s 已下树", clanMember.GameName))
	} else {
		err = gvg_member_extra.ReportUpTree(qqid, agentQqid, gvgGroup.GvgId)
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
		bot.Send(mess.FromGroupID, 2, fmt.Sprintf("%s 已记录挂树", clanMember.GameName))
	}

}
