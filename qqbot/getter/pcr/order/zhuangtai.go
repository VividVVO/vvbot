package order

import (
	"github.com/vivid-vvo/vvbot/qqbot/getter"
	"github.com/vivid-vvo/vvbot/qqbot/model/bot"
	"github.com/vivid-vvo/vvbot/qqbot/model/pcr"
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

func (l Zhuangtai) Run(mess getter.MeassageData, ms string, atqq int) {
	bot.Send(mess.FromGroupID, 2, pcr.GetBossStateStr(mess.FromGroupID))
}
