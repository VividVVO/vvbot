package plugin

import (
	"github.com/vivid-vvo/vvbot/qqbot/plugins"
)

type Chongqi struct {
}

func (l Chongqi) CheckOrder(cm string) bool {
	return true
}

func (l Chongqi) IsNotCheckOrder() bool {
	return false
}

func (l Chongqi) GetOrders() []string {
	return []string{
		"^重启$",
	}
}

func (l Chongqi) Run(mess plugins.MeassageData, ms string, atqq int64) {

}
