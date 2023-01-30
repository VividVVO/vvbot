package plugin

import (
	"fmt"
	"github.com/vivid-vvo/vvbot/app/model/gvg_group"
	"github.com/vivid-vvo/vvbot/app/model/gvg_member_extra"
	"github.com/vivid-vvo/vvbot/library/tools"
	"github.com/vivid-vvo/vvbot/qqbot/model/bot"
	"github.com/vivid-vvo/vvbot/qqbot/model/pcr"
	"github.com/vivid-vvo/vvbot/qqbot/plugins"
	"regexp"
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
		"^(挂(树|🌲)|上(树|🌲))( +)?(\\?|？)?$",
		"^(挂(树|🌲)|上(树|🌲))( +)?(\\:|：)",

		"^(挂(树|🌲)|上(树|🌲))(取消|撤销)( +)?(\\?|？)?$",
		"^(取消|撤销)(挂(树|🌲)|上(树|🌲))( +)?(\\?|？)?$",
		"^下(树|🌲)( +)?$",
	}
}

func (l Guashu) Run(mess plugins.MeassageData, cm string, atqq int64) {
	var agentQqid, qqid int64
	if atqq != 0 {
		agentQqid = mess.FromUserID
		qqid = atqq
	} else {
		qqid = mess.FromUserID
	}
	cm = strings.ToLower(cm)
	isCancel := tools.Compare(cm, "^(取消|撤销|下树)|(取消|撤销)$")
	isQuery := tools.Compare(cm, "[?|？]")

	reg := regexp.MustCompile(`[:|：](.+)`)
	message := reg.FindString(cm)
	if message != "" {
		message = string([]rune(message)[1:])
	}
	_, clanMember, err := pcr.GetClanGroupAndUserGroupToCheck(mess.FromGroupID, qqid)
	if err != nil {
		bot.Send(mess.FromGroupID, mess.SendToType, err.Error())
		return
	}
	_, gvgGroup, err := pcr.GetGvgGroupDataAtGroupIdToCheck(mess.FromGroupID)
	if err != nil {
		bot.Send(mess.FromGroupID, mess.SendToType, err.Error())
		return
	}
	treeData, err := gvg_member_extra.GetIsUpTree(qqid, gvgGroup.GvgId)
	if err != nil {
		bot.Send(mess.FromGroupID, mess.SendToType, err.Error())
		return
	}
	if treeData != nil && !isCancel {
		bot.Send(mess.FromGroupID, mess.SendToType, fmt.Sprintf("%s 您已在树上 [%s]", clanMember.GameName, treeData.Message))
		return
	}
	if isQuery || isCancel && treeData == nil {
		bot.Send(mess.FromGroupID, mess.SendToType, fmt.Sprintf("%s 您未在树上", clanMember.GameName))
		return
	}
	if isCancel {
		// 删树
		err = gvg_member_extra.ReportDownTree(qqid, gvgGroup.GvgId, 0)
		if err != nil {
			bot.Send(mess.FromGroupID, mess.SendToType, err.Error())
			return
		}
		bot.Send(mess.FromGroupID, mess.SendToType, fmt.Sprintf("%s 已下树", clanMember.GameName))
	} else {
		// 报告上树
		err = gvg_member_extra.ReportUpTree(qqid, agentQqid, gvgGroup.GvgId, message)
		if err != nil {
			bot.Send(mess.FromGroupID, mess.SendToType, err.Error())
			return
		}
		bot.Send(mess.FromGroupID, mess.SendToType, fmt.Sprintf("%s 已记录挂树", clanMember.GameName))
		if qqid == gvgGroup.ChallengeStratQqid {
			// 出刀锁定解锁
			if err = gvg_group.CancelChallenge(gvgGroup.GvgId); err != nil {
				bot.Send(mess.FromGroupID, mess.SendToType, err.Error())
				return
			}
			msg := fmt.Sprintf("boss挑战已可申请\n现在%d周目，%d号boss\n生命值%s", gvgGroup.BossCycle, gvgGroup.BossNum, tools.NumberFormat(gvgGroup.BossHp))
			bot.Send(mess.FromGroupID, mess.SendToType, msg)
		}
	}

}
