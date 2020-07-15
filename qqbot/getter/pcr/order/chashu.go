package order

import (
	"fmt"
	"github.com/vivid-vvo/vvbot/app/model/clan_member"

	"github.com/vivid-vvo/vvbot/app/model/gvg_member_extra"

	"github.com/vivid-vvo/vvbot/qqbot/getter"
	"github.com/vivid-vvo/vvbot/qqbot/model/bot"
	"github.com/vivid-vvo/vvbot/qqbot/model/pcr"
)

type Chadao struct {
}

func (l Chadao) CheckOrder(cm string) bool {
	return true
}

func (l Chadao) IsNotCheckOrder() bool {
	return false
}

func (l Chadao) GetOrders() []string {
	return []string{
		"^查树$",
		"^查询挂树$",
		"^挂树查询$",
		"^树上成员$",
		"^树上的人$",
		"^树上有那些人人$",
	}
}

func (l Chadao) Run(mess getter.MeassageData, cm string, atqq int) {
	clanGroup, err := pcr.GetClanGroupAndChack(mess.FromGroupID)
	if err != nil {
		bot.Send(mess.FromGroupID, 2, err.Error())
		return
	}
	entityList, err := gvg_member_extra.GetAllUpTree()
	if err != nil {
		bot.Send(mess.FromGroupID, 2, err.Error())
		return
	}
	var treeStr string
	for _, entity := range entityList {
		clanMember, err := clan_member.GetClanMember(entity.Qqid, clanGroup.GroupId)
		if err != nil {
			bot.Send(mess.FromGroupID, 2, err.Error())
			return
		}
		treeStr += "\n" + clanMember.GameName
	}
	if treeStr == "" {
		bot.Send(mess.FromGroupID, 2, "当前无挂树成员")
	} else {
		bot.Send(mess.FromGroupID, 2, fmt.Sprintf("挂树的成员：\n%s", treeStr))
	}
}
