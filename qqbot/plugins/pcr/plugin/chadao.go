package plugin

import (
	"fmt"
	"github.com/gogf/gf/frame/g"
	"github.com/vivid-vvo/vvbot/qqbot/model/bot"
	"github.com/vivid-vvo/vvbot/qqbot/model/pcr"
	"github.com/vivid-vvo/vvbot/qqbot/plugins"
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
		"^查刀$",
	}
}

func (l Chadao) Run(mess plugins.MeassageData, cm string, atqq int64) {
	var WebUrl = g.Cfg().GetString("server.WebUrl")
	clanGroup, err := pcr.GetClanGroupAndChack(mess.FromGroupID)
	if err != nil {
		bot.Send(mess.FromGroupID, mess.SendToType, err.Error())
		return
	}

	bot.Send(mess.FromGroupID, mess.SendToType, fmt.Sprintf("查刀 %s %s/clan/%d/record", clanGroup.GroupName, WebUrl, clanGroup.GroupId))
}
