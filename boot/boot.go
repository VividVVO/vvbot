package boot

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/glog"
	"github.com/vivid-vvo/vvbot/app/model/boss_data"
	_ "github.com/vivid-vvo/vvbot/packed"
	"github.com/vivid-vvo/vvbot/qqbot/bot/coolq"
	"github.com/vivid-vvo/vvbot/qqbot/bot/iotqq"

	"github.com/vivid-vvo/vvbot/qqbot/events"
	"github.com/vivid-vvo/vvbot/qqbot/model/bot"
	"os"
)

var botName = g.Cfg().GetString("bot.botname")

func init() {
	os.Setenv("ZONEINFO", "./static/data.zip")
	boss_data.Init()
	if botName != "" {
		var entity bot.Bot
		switch botName {
		case "iotqq":
			entity = new(iotqq.Iotqq)
			bot.Init(entity)
		case "coolq":
			entity = new(coolq.Coolq)
			bot.Init(entity)
		default:
			glog.Warning("指定机器人模块不存在！")
			return
		}
		// 绑定事件
		entity.BindOnGroupMsgs(events.OnGroupMsgs)
		entity.BindOnFriendMsgs(events.OnFriendMsgs)
	}

}
