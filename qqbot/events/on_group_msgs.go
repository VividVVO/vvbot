package events

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/util/grand"
	"github.com/vivid-vvo/vvbot/library/Tools"
	"github.com/vivid-vvo/vvbot/qqbot/getter"
	"github.com/vivid-vvo/vvbot/qqbot/getter/pcr"
	"github.com/vivid-vvo/vvbot/qqbot/model/bot"
	"log"
	"regexp"
	"strconv"
)

// 添加需要的功能集
var executeOrders = [][]getter.ExecuteOrder{
	pcr.ExecuteOrders,
}

// 机器人指令前缀
var prefix = g.Cfg().GetString("iotbot.prefix")

func OnGroupMsgs(message getter.MeassageData) {
	/*
		message.Content 消息内容 string
		message.FromGroupID 来源QQ群 int
		message.FromUserID 来源QQ int64
		message.iotqqType 消息类型 string
	*/
	var atqq int
	log.Println("群聊消息: ", message.FromNickName+"<"+strconv.FormatInt(message.FromUserID, 10)+">: "+message.Content, message.AtQQList)
	msg := message.Content
	if len(message.AtQQList) > 0 {
		atqq = int(message.AtQQList[0])
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
	if atqq == int(message.FromUserID) {
		bot.Send(message.FromGroupID, 2, "不会吧不会吧，不会真的有人艾特自己吧")
		atqq = 0
	}
	if strconv.FormatInt(message.FromUserID, 10) == bot.GetBotQQID() {
		return
	}
	if strconv.Itoa(atqq) == bot.GetBotQQID() {
		msg := []string{"憨憨", "有事？", "你在教我做事？", "不会吧不会吧，不会真的有人艾特机器人吧", "就这？", "憨憨"}[grand.N(0, 5)]
		bot.Send(message.FromGroupID, 2, msg)
		return
	}
	for _, executeOrders2 := range executeOrders {
		for _, executeOrder := range executeOrders2 {
			if !executeOrder.IsNotCheckOrder() {
				if Tools.CheckOrder(msg, executeOrder.GetOrders()) {
					executeOrder.Run(message, msg, atqq)
				}
			} else if executeOrder.CheckOrder(message.Content) {
				executeOrder.Run(message, msg, atqq)
			}
		}
	}
}
