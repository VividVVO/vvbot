package plugin

import (
	"fmt"
	"github.com/vivid-vvo/vvbot/app/service/check"
	"github.com/vivid-vvo/vvbot/qqbot/model/bot"
	"github.com/vivid-vvo/vvbot/qqbot/plugins"
)

type LeaveGroup struct {
}

func (l LeaveGroup) CheckOrder(cm string) bool {
	return true
}

func (l LeaveGroup) IsNotCheckOrder() bool {
	return false
}

func (l LeaveGroup) GetOrders() []string {
	return []string{
		"^退出群组$",
		"^退出此群$",
		"^退出本群$",
	}
}

func (l LeaveGroup) Run(mess plugins.MeassageData, ms string, atqq int64) {
	if !check.CheckAuthorityGroup(mess.FromUserID, check.AuthAdmin, 0) {
		bot.Send(mess.FromSourceID, mess.SendToType, "权限不足")
		return
	}
	bot.Send(mess.FromSourceID, mess.SendToType, fmt.Sprintf("vvbot 即将退出本群"))
	if err := bot.LeaveChat(mess.FromGroupID, mess.SendToType, false); err != nil {
		bot.Send(mess.FromSourceID, mess.SendToType, fmt.Sprintf("退出失败 %v", err))
	}

}
