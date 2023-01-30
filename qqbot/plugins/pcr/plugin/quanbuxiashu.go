package plugin

import (
	"fmt"
	"github.com/vivid-vvo/vvbot/app/model/clan_member"
	"github.com/vivid-vvo/vvbot/app/model/gvg_member_extra"
	"github.com/vivid-vvo/vvbot/qqbot/model/bot"
	"github.com/vivid-vvo/vvbot/qqbot/model/pcr"
	"github.com/vivid-vvo/vvbot/qqbot/plugins"
	"strings"
)

type Quanbuxiashu struct {
}

func (l Quanbuxiashu) CheckOrder(cm string) bool {
	return true
}
func (l Quanbuxiashu) IsNotCheckOrder() bool {
	return false
}
func (l Quanbuxiashu) GetOrders() []string {
	return []string{
		"^全.+下(树|🌲)( +)?$",
	}
}
func (l Quanbuxiashu) Run(mess plugins.MeassageData, cm string, atqq int64) {
	cm = strings.ToLower(cm)
	clanGroup, _, err := pcr.GetGvgGroupDataAtGroupIdToCheck(mess.FromGroupID)
	if err != nil {
		bot.Send(mess.FromGroupID, mess.SendToType, err.Error())
		return
	}
	upTreeList, err := gvg_member_extra.GetAllUpTree(clanGroup.GvgId)
	if err != nil {
		bot.Send(mess.FromGroupID, mess.SendToType, err.Error())
		return
	}
	if len(upTreeList) == 0 {
		bot.Send(mess.FromGroupID, mess.SendToType, "没有人在挂树")
		return
	}
	var msg, errMsg string
	for _, upTree := range upTreeList {
		clanMember, err := clan_member.GetClanMember(upTree.Qqid, clanGroup.GroupId)
		if err != nil {
			errMsg += fmt.Sprintf("\nqq: %d（%s）, 下树失败 -> %s", upTree.Qqid, clanMember.GameName, err.Error())
			continue
		}
		err = gvg_member_extra.ReportDownTree(upTree.Qqid, clanGroup.GvgId, 0)
		if err != nil {
			errMsg += fmt.Sprintf("\nqq: %d（%s）, 下树失败 -> %s", upTree.Qqid, clanMember.GameName, err.Error())
			continue
		}
		msg += fmt.Sprintf("\n%s %s", clanMember.GameName, bot.GetAtQQStr(int64(upTree.Qqid)))
	}
	bot.Send(mess.FromGroupID, mess.SendToType, fmt.Sprintf("以下成员已下树：%s", errMsg+msg))
}
