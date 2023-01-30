package plugin

import (
	"github.com/vivid-vvo/vvbot/qqbot/model/bot"
	"github.com/vivid-vvo/vvbot/qqbot/model/pcr"
	"github.com/vivid-vvo/vvbot/qqbot/plugins"
)

type Zhuangtai struct {
}

func (l Zhuangtai) CheckOrder(cm string) bool {
	return true
}

func (l Zhuangtai) IsNotCheckOrder() bool {
	return false
}

func (l Zhuangtai) GetOrders() []string {
	return []string{
		"^状态$",
		"^state$",
		"^查看状态$",
		"^查询状态$",
	}
}

func (l Zhuangtai) Run(mess plugins.MeassageData, ms string, atqq int64) {
	bot.Send(mess.FromGroupID, mess.SendToType, pcr.GetBossStateStr(mess.FromGroupID))
}
