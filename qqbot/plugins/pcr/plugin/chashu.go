package plugin

import (
	"fmt"
	"github.com/vivid-vvo/vvbot/app/model/clan_member"
	"github.com/vivid-vvo/vvbot/app/model/gvg_member_extra"
	"time"

	"github.com/vivid-vvo/vvbot/qqbot/model/bot"
	"github.com/vivid-vvo/vvbot/qqbot/model/pcr"
	"github.com/vivid-vvo/vvbot/qqbot/plugins"
)

type Chasu struct {
}

func (l Chasu) CheckOrder(cm string) bool {
	return true
}

func (l Chasu) IsNotCheckOrder() bool {
	return false
}

func (l Chasu) GetOrders() []string {
	return []string{
		"^查树$",
		"^查1$",

		"^查询挂树$",
		"^挂树查询$",
		"^树上成员$",
		"^树上的人$",
		"^树上有那些人人$",
	}
}

func (l Chasu) Run(mess plugins.MeassageData, cm string, atqq int64) {
	clanGroup, err := pcr.GetClanGroupAndChack(mess.FromGroupID)
	if err != nil {
		bot.Send(mess.FromGroupID, mess.SendToType, err.Error())
		return
	}
	entityList, err := gvg_member_extra.GetAllUpTree(clanGroup.GvgId)
	if err != nil {
		bot.Send(mess.FromGroupID, mess.SendToType, err.Error())
		return
	}
	if len(entityList) == 0 {
		bot.Send(mess.FromGroupID, mess.SendToType, "当前无挂树成员")
		return
	}
	var treeStr string
	for _, entity := range entityList {
		clanMember, err := clan_member.GetClanMember(entity.Qqid, clanGroup.GroupId)
		if err != nil {
			bot.Send(mess.FromGroupID, mess.SendToType, err.Error())
			return
		}
		treeStr += fmt.Sprintf("\n%s [%s] -> %d分钟", clanMember.GameName, entity.Message, (time.Now().Unix()-entity.Time)/60)
	}
	bot.Send(mess.FromGroupID, mess.SendToType, fmt.Sprintf("挂树的成员：\n%s", treeStr))

}
