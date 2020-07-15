package boot

import (
	"github.com/gogf/gf/frame/g"
	"github.com/vivid-vvo/vvbot/app/model/boss_data"
	_ "github.com/vivid-vvo/vvbot/packed"
	"github.com/vivid-vvo/vvbot/qqbot/bot/iotqq"
	"github.com/vivid-vvo/vvbot/qqbot/events"
	"github.com/vivid-vvo/vvbot/qqbot/model/bot"
)

var botName = g.Cfg().GetString("bot.botname")

func init() {
	boss_data.Init()

	var entity bot.Bot
	switch botName {
	case "iotqq":
		entity = new(iotqq.Iotqq)
		bot.Init(entity)
	}

	// 绑定事件
	entity.BindOnGroupMsgs(events.OnGroupMsgs)
}
