package plugin

import (
	"fmt"
	"github.com/gogf/gf/frame/g"
	"github.com/vivid-vvo/vvbot/qqbot/model/bot"
	"github.com/vivid-vvo/vvbot/qqbot/model/pcr"
	"github.com/vivid-vvo/vvbot/qqbot/plugins"
)

type Mianban struct {
}

func (l Mianban) CheckOrder(cm string) bool {
	return true
}

func (l Mianban) IsNotCheckOrder() bool {
	return false
}

func (l Mianban) GetOrders() []string {
	return []string{
		"^面板$",
	}
}
func (l Mianban) Run(mess plugins.MeassageData, cm string, atqq int64) {
	var WebUrl = g.Cfg().GetString("server.WebUrl")
	clanGroup, err := pcr.GetClanGroupAndChack(mess.FromGroupID)
	if err != nil {
		bot.Send(mess.FromGroupID, mess.SendToType, err.Error())
		return
	}

	bot.Send(mess.FromGroupID, mess.SendToType, fmt.Sprintf("面板 %s %s/clan/%d/home", clanGroup.GroupName, WebUrl, clanGroup.GroupId))
}
