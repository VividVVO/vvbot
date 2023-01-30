package plugin

import (
	"fmt"
	"github.com/vivid-vvo/vvbot/app/model/gvg_member_extra"
	"github.com/vivid-vvo/vvbot/qqbot/model/bot"
	"github.com/vivid-vvo/vvbot/qqbot/model/pcr"
	"github.com/vivid-vvo/vvbot/qqbot/plugins"

	"github.com/vivid-vvo/vvbot/library/tools"
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
		"^sl( +)?(\\:|：)",

		"^查(询)?sl( +)?(\\?|？)?$",
		"^sl查(询)?( +)?(\\?|？)?$",
		"^(撤销|取消)sl$",
		"^sl(取消|撤销)$",
		"^sl(取消|撤销)$",
		"^sl( )?(no|cancel)$",
		"^(no|cancel)( )?sl$",
	}
}

func (l Sl) Run(mess plugins.MeassageData, cm string, atqq int64) {
	var agentQqid, qqid int64
	if atqq != 0 {
		agentQqid = mess.FromUserID
		qqid = atqq
	} else {
		qqid = mess.FromUserID
	}
	cm = strings.ToLower(cm)
	regexp1 := regexp.MustCompile(`@.*? `)
	cm = regexp1.ReplaceAllString(cm, "")
	isCancel := tools.Compare(cm, "^(取消|撤销|no|cancel)|(取消|撤销|no|cancel)$")

	isQuery := tools.Compare(cm, "[?|？]|查")

	reg := regexp.MustCompile(`[:|：](.+)`)
	message := reg.FindString(cm)
	if message != "" {
		message = string([]rune(message)[1:])
	}
	clanGroup, clanMember, err := pcr.GetClanGroupAndUserGroupToCheck(mess.FromGroupID, qqid)
	if err != nil {
		bot.Send(mess.FromGroupID, mess.SendToType, err.Error())
		return
	}
	_, gvgGroup, err := pcr.GetGvgGroupDataAtGroupIdToCheck(mess.FromGroupID)
	if err != nil {
		bot.Send(mess.FromGroupID, mess.SendToType, err.Error())
		return
	}
	slData, err := gvg_member_extra.GetDaySL(qqid, gvgGroup.GvgId, clanGroup.GameServer)
	if err != nil {
		bot.Send(mess.FromGroupID, mess.SendToType, err.Error())
		return
	}
	if slData != nil && !isCancel {
		bot.Send(mess.FromGroupID, mess.SendToType, fmt.Sprintf("%s 今日已SL [%s]", clanMember.GameName, slData.Message))
		return
	}
	if isQuery || isCancel && slData == nil {
		if slData == nil {
			bot.Send(mess.FromGroupID, mess.SendToType, fmt.Sprintf("%s 今日未SL", clanMember.GameName))
			return
		}
	}
	if isCancel {
		err = gvg_member_extra.CancelDaySL(qqid, clanGroup.GvgId, clanGroup.GameServer)
		if err != nil {
			bot.Send(mess.FromGroupID, mess.SendToType, err.Error())
			return
		}
		bot.Send(mess.FromGroupID, mess.SendToType, fmt.Sprintf("%s 已取消SL", clanMember.GameName))
	} else {
		err = gvg_member_extra.ReportDaySL(qqid, agentQqid, gvgGroup.GvgId, message)
		if err != nil {
			bot.Send(mess.FromGroupID, mess.SendToType, err.Error())
			return
		}
		bot.Send(mess.FromGroupID, mess.SendToType, fmt.Sprintf("%s 已记录SL", clanMember.GameName))
		if qqid == gvgGroup.ChallengeStratQqid {
			// 出刀锁定解锁
			/*if err = gvg_group.CancelChallenge(gvgGroup.GvgId); err != nil {
				bot.Send(mess.FromGroupID, mess.SendToType, err.Error())
				return
			}*/

			msg := fmt.Sprintf("boss挑战已可申请\n现在%d周目，%d号boss\n生命值%s", gvgGroup.BossCycle, gvgGroup.BossNum, tools.NumberFormat(gvgGroup.BossHp))
			bot.Send(mess.FromGroupID, mess.SendToType, msg)
		}

		// 下树
		treeData, err := gvg_member_extra.GetIsUpTree(qqid, gvgGroup.GvgId)
		if err != nil {
			bot.Send(mess.FromGroupID, mess.SendToType, err.Error())
			return
		}
		if isQuery || isCancel && treeData == nil {
			return
		}
		err = gvg_member_extra.ReportDownTree(qqid, gvgGroup.GvgId, 0)
		if err != nil {
			bot.Send(mess.FromGroupID, mess.SendToType, err.Error())
			return
		}
	}

}
