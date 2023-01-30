package plugin

import (
	"fmt"
	"github.com/gogf/gf/frame/g"
	"github.com/vivid-vvo/vvbot/qqbot/model/bot"
	"github.com/vivid-vvo/vvbot/qqbot/plugins"
)

type Help struct {
}

func (l Help) CheckOrder(cm string) bool {
	return true
}

func (l Help) IsNotCheckOrder() bool {
	return false
}

func (l Help) GetOrders() []string {
	return []string{
		"^help$",
		"^帮助$",
	}
}

var WebUrl = g.Cfg().GetString("server.WebUrl")

func (l Help) Run(mess plugins.MeassageData, ms string, atqq int64) {
	bot.Send(mess.FromGroupID, mess.SendToType, fmt.Sprintf("help %s/help.html", WebUrl))
}
