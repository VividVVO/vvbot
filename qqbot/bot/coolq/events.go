package coolq

import (
	qqbotapi "github.com/catsworld/qq-bot-api"
	"github.com/vivid-vvo/vvbot/library/tools"
	"github.com/vivid-vvo/vvbot/qqbot/plugins"
	"regexp"
	"strconv"
	"strings"
)

type Message struct {
	CurrentPacket CurrentPacket `json:"CurrentPacket"`
	CurrentQQ     int64         `json:"CurrentQQ"`
}
type CurrentPacket struct {
	Data      Data   `json:"Data"`
	WebConnID string `json:"WebConnId"`
}
type Data struct {
	Content       string      `json:"Content"`
	FromGroupID   int         `json:"FromGroupId"`
	FromGroupName string      `json:"FromGroupName"`
	FromNickName  string      `json:"FromNickName"`
	FromUserID    int64       `json:"FromUserId"`
	MsgRandom     int         `json:"MsgRandom"`
	MsgSeq        int         `json:"MsgSeq"`
	MsgTime       int64       `json:"MsgTime"`
	MsgType       string      `json:"MsgType"`
	RedBaginfo    interface{} `json:"RedBaginfo"`
}

func OnGroupMsgs(coolbot *qqbotapi.BotAPI, args qqbotapi.Update) {
	var msg string
	var atqq []int64
	message := args.Message.Text
	reg := regexp.MustCompile(`\[CQ:(\S{0,4}),qq=(\d+)\]( )?`)
	strList := reg.FindAllString(message, -1)
	for _, str := range strList {
		message = strings.Replace(message, str, "", -1)
		qqIDStr := tools.GetBetweenStr(str, "[CQ:at,qq=", "]")
		qqID, _ := strconv.ParseInt(qqIDStr, 10, 64)
		atqq = append(atqq, qqID)
		// @艾特全体过滤
		if qqIDStr == "all" {
			return
		}
	}
	msg = message
	var mess plugins.MeassageData
	mess.FromGroupID = args.GroupID
	mess.AtQQList = atqq
	mess.FromUserID = args.Message.From.ID
	mess.FromNickName = args.Message.From.String()
	mess.SendToType = 2
	mess.Content = msg
	mess.CurrentQQ = coolbot.Self.ID
	mess.FromSourceID = args.GroupID

	OnGroupMsgsFunc(mess)
}

func OnFriendMsgs(coolbot *qqbotapi.BotAPI, args qqbotapi.Update) {
	var msg string
	message := args.Message.Text
	reg := regexp.MustCompile(`\[CQ:(\S{0,4}),qq=(\d+)\]( )?`)
	strList := reg.FindAllString(message, -1)
	for _, str := range strList {
		message = strings.Replace(message, str, "", -1)
	}
	msg = message
	var mess plugins.MeassageData
	mess.FromGroupID = args.UserID
	mess.FromUserID = args.Message.From.ID
	mess.FromUin = args.Message.From.ID
	mess.FromNickName = args.Message.From.String()
	mess.SendToType = 1
	mess.Content = msg
	mess.CurrentQQ = coolbot.Self.ID

	mess.FromSourceID = mess.FromUserID
	OnFriendMsgsFunc(mess)
}
