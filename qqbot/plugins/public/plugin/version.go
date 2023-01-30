package plugin

import (
	"fmt"
	"github.com/gogf/gf/os/gcache"
	"github.com/vivid-vvo/vvbot/qqbot/model/bot"
	"github.com/vivid-vvo/vvbot/qqbot/plugins"
)

type Version struct {
}

func (l Version) CheckOrder(cm string) bool {
	return true
}

func (l Version) IsNotCheckOrder() bool {
	return false
}

func (l Version) GetOrders() []string {
	return []string{
		"^version$",
		"^版本$",
	}
}

func (l Version) Run(mess plugins.MeassageData, ms string, atqq int64) {
	var VERSION = gcache.Get("VERSION")
	bot.Send(mess.FromSourceID, mess.SendToType, fmt.Sprintf("vvbot[%s]", VERSION))
}
