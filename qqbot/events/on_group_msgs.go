package events

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/util/grand"
	"github.com/vivid-vvo/vvbot/library/tools"
	"github.com/vivid-vvo/vvbot/qqbot/model/bot"
	"github.com/vivid-vvo/vvbot/qqbot/plugins"
	"github.com/vivid-vvo/vvbot/qqbot/plugins/pcr"
	public "github.com/vivid-vvo/vvbot/qqbot/plugins/public"
	"log"
	"regexp"
	"strconv"
)

// 添加需要的功能集
var executeGroupOrders = [][]plugins.ExecutePlugin{
	pcr.ExecuteGroupPlugins,
	public.ExecuteGroupPlugins,
}

// 机器人指令前缀
var prefix = g.Cfg().GetString("bot.Prefix")

func OnGroupMsgs(message plugins.MeassageData) {
	/*
		message.Content 消息内容 string
		message.FromGroupID 来源QQ群 int
		message.FromUserID 来源QQ int64
		message.iotqqType 消息类型 string
	*/
	var atqq int64
	log.Println("群聊消息: ", message.FromNickName+"<"+strconv.FormatInt(message.FromUserID, 10)+">: "+message.Content, message.AtQQList)
	msg := message.Content
	if len(message.AtQQList) > 0 {
		atqq = message.AtQQList[0]
	}
	var msg2 string
	regexp1 := regexp.MustCompile("^" + prefix)
	if !regexp1.MatchString(msg) {
		return
	}
	msg = regexp1.ReplaceAllString(msg, "")

	regexp1 = regexp.MustCompile("^( +)")
	msg2 = regexp1.ReplaceAllString(msg, "")
	msg = msg2
	if atqq == message.FromUserID {
		bot.Send(message.FromGroupID, 2, "不会吧不会吧，不会真的有人艾特自己吧")
		atqq = 0
	}
	if strconv.FormatInt(message.FromUserID, 10) == bot.GetBotQQID() {
		return
	}

	if atqq == bot.GetBotQQIDToInt64() {
		msg := []string{"憨憨", "有事？", "你在教我做事？", "不会吧不会吧，不会真的有人艾特机器人吧", "就这？", "憨憨"}[grand.N(0, 5)]
		bot.Send(message.FromGroupID, 2, msg)
		return
	}
	for _, executeGroupOrders2 := range executeGroupOrders {
		for _, executeGroupOrders := range executeGroupOrders2 {
			if !executeGroupOrders.IsNotCheckOrder() {
				if tools.CheckOrder(msg, executeGroupOrders.GetOrders()) {
					executeGroupOrders.Run(message, msg, atqq)
				}
			} else if executeGroupOrders.CheckOrder(message.Content) {
				executeGroupOrders.Run(message, msg, atqq)
			}
		}
	}
}
