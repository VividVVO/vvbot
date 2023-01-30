package events

import (
	"github.com/vivid-vvo/vvbot/library/tools"
	"github.com/vivid-vvo/vvbot/qqbot/plugins"
	public "github.com/vivid-vvo/vvbot/qqbot/plugins/public"
	user "github.com/vivid-vvo/vvbot/qqbot/plugins/user"
	"log"
	"regexp"
	"strconv"
)

// 添加需要的功能集
var executeFriendOrders = [][]plugins.ExecutePlugin{
	user.ExecuteFriendPlugins,
	public.ExecuteFriendPlugins,
}

func OnFriendMsgs(message plugins.MeassageData) {
	/*
		message.Content 消息内容 string
		message.FromGroupID 来源QQ群 int
		message.FromUserID 来源QQ int64
		message.iotqqType 消息类型 string
	*/
	var atqq int64
	log.Println("私聊消息: ", message.FromNickName+"<"+strconv.FormatInt(message.FromUserID, 10)+">: "+message.Content, message.AtQQList)
	msg := message.Content
	var msg2 string
	regexp1 := regexp.MustCompile("^( +)")
	msg2 = regexp1.ReplaceAllString(msg, "")
	msg = msg2
	for _, executeFriendOrders2 := range executeFriendOrders {
		for _, executeFriendOrders := range executeFriendOrders2 {
			if !executeFriendOrders.IsNotCheckOrder() {
				if tools.CheckOrder(msg, executeFriendOrders.GetOrders()) {
					executeFriendOrders.Run(message, msg, atqq)
				}
			} else if executeFriendOrders.CheckOrder(message.Content) {
				executeFriendOrders.Run(message, msg, atqq)
			}
		}
	}
}
