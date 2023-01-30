package plugin

import (
	"fmt"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/util/grand"
	"github.com/vivid-vvo/vvbot/app/model/user"

	"github.com/vivid-vvo/vvbot/qqbot/model/bot"
	"github.com/vivid-vvo/vvbot/qqbot/plugins"
)

type Denglu struct {
}

func (l Denglu) CheckOrder(cm string) bool {
	return true
}

func (l Denglu) IsNotCheckOrder() bool {
	return false
}

func (l Denglu) GetOrders() []string {
	return []string{
		"^登录$",
		"^login$",
	}
}

var WebUrl = g.Cfg().GetString("server.WebUrl")

func (l Denglu) Run(mess plugins.MeassageData, ms string, atqq int64) {
	qqid := mess.FromUin
	// groupID :=

	userInfo, err := user.GetProfile(qqid)
	if err != nil {
		bot.SendPrivate(605538767, err.Error(), qqid)
		return
	}

	if userInfo == nil {
		// 注册
		err := user.SignUp(qqid, grand.Letters(12), mess.FromNickName, "")
		if err != nil {
			bot.SendPrivate(605538767, err.Error(), qqid)
			return
		}
	}
	authCookie, err := user.GetLoginAuth(mess.FromUin)
	if err != nil {
		bot.SendPrivate(605538767, err.Error(), qqid)

		return
	}

	bot.SendPrivate(605538767, fmt.Sprintf("登录链接 %s/login?qqid=%d&key=%s", WebUrl, mess.FromUin, authCookie), qqid)
}
