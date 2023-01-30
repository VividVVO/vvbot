package plugin

import (
	"fmt"
	"github.com/gogf/gf/util/grand"
	"github.com/vivid-vvo/vvbot/app/model/user"
	"github.com/vivid-vvo/vvbot/qqbot/model/bot"
	"github.com/vivid-vvo/vvbot/qqbot/plugins"
)

type Suijimima struct {
}

func (l Suijimima) CheckOrder(cm string) bool {
	return true
}

func (l Suijimima) IsNotCheckOrder() bool {
	return false
}

func (l Suijimima) GetOrders() []string {
	return []string{
		"^随机密码$",
		"^重置密码$",
	}
}

func (l Suijimima) Run(mess plugins.MeassageData, ms string, atqq int64) {
	qqid := mess.FromUin
	userInfo, err := user.GetProfile(qqid)
	if err != nil {
		bot.Send(mess.FromGroupID, mess.SendToType, err.Error())
		return
	}
	if userInfo == nil {
		// 注册
		err := user.SignUp(qqid, grand.Letters(12), mess.FromNickName, "")
		if err != nil {
			bot.Send(mess.FromGroupID, mess.SendToType, err.Error())
			return
		}
	}
	password := grand.S(8, false)
	err = user.ChangePassword(qqid, password)
	if err != nil {
		bot.Send(mess.FromGroupID, mess.SendToType, err.Error())
		return
	}
	bot.Send(mess.FromGroupID, mess.SendToType, fmt.Sprintf("您的新密码：%s", password))
}
