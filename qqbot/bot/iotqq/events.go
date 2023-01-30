package iotqq

import (
	"encoding/json"
	"fmt"
	"github.com/gogf/gf/os/glog"
	"github.com/gogf/gf/util/gconv"
	gosocketio "github.com/graarh/golang-socketio"
	"github.com/vivid-vvo/vvbot/qqbot/plugins"
	"log"
	"regexp"
)

type Message struct {
	CurrentPacket CurrentPacket `json:"CurrentPacket"`
	CurrentQQ     int64         `json:"CurrentQQ"`
}
type CurrentPacket struct {
	Data  Data  `json:"Data"`
	ToUin int64 `json:"ToUin"`
}
type Data struct {
	ToUin         int64       `json:"ToUin"`
	FromUin       int64       `json:"FromUin"`
	Content       string      `json:"Content"`
	FromGroupID   int64       `json:"FromGroupId"`
	FromGroupName string      `json:"FromGroupName"`
	FromNickName  string      `json:"FromNickName"`
	FromUserID    int64       `json:"FromUserId"`
	MsgRandom     int         `json:"MsgRandom"`
	MsgSeq        int         `json:"MsgSeq"`
	MsgTime       int64       `json:"MsgTime"`
	MsgType       string      `json:"MsgType"`
	RedBaginfo    interface{} `json:"RedBaginfo"`
}
type content struct {
	Content string  `json:"Content"`
	UserID  []int64 `json:"UserID"`
}

func OnGroupMsgs(h *gosocketio.Channel, args Message) {
	var msg string
	var atqq []int64
	message := args.CurrentPacket.Data
	content := content{}
	switch message.MsgType {
	case "TextMsg":
		msg = message.Content
	case "AtMsg":
		err := json.Unmarshal([]byte(message.Content), &content)
		if err != nil {
			log.Fatalln(err)
			return
		}
		msg = content.Content
		if len(content.UserID) > 0 {
			atqq = content.UserID
		}
	}
	var mess plugins.MeassageData
	if err := gconv.Struct(args.CurrentPacket.Data, &mess); err != nil {
		glog.Error(err)
		return
	}
	var name string
	// 清除文本内的@信息
	if len(atqq) > 0 {
		// @艾特全体过滤
		if atqq[0] == 0 {
			return
		}
		member := getGroupUserData(message.FromGroupID, atqq[0])
		// 群成员获取错误
		if member == nil {
			glog.Error("iotqq.GetGroupUserData")
			send(message.FromGroupID, 2, "内部错误", 0)
			return
		}
		if member.GroupCard != "" {
			name = member.GroupCard
		} else if member.AutoRemark != "" {
			name = member.AutoRemark
		} else if member.NickName != "" {
			name = member.NickName
		}
		regexp1, err := regexp.Compile(fmt.Sprintf("@%s( +)?", name))
		if err != nil {
			return
		}
		msg = regexp1.ReplaceAllString(msg, "")
		regexp1, err = regexp.Compile("@( +)?")
		msg = regexp1.ReplaceAllString(msg, "")
		if err != nil {
			return
		}

	}
	mess.Content = msg
	mess.SendToType = 2
	mess.AtQQList = atqq
	mess.CurrentQQ = args.CurrentQQ
	mess.FromSourceID = message.FromGroupID
	OnGroupMsgsFunc(mess)
}

func OnFriendMsgs(h *gosocketio.Channel, args Message) {
	var msg string
	message := args.CurrentPacket.Data
	content := content{}
	msg = message.Content
	err := json.Unmarshal([]byte(message.Content), &content)
	if err == nil {
		msg = content.Content
	}
	var mess plugins.MeassageData
	if err := gconv.Struct(args.CurrentPacket.Data, &mess); err != nil {
		glog.Error(err)
		return
	}
	mess.SendToType = 1
	mess.Content = msg
	mess.FromGroupID = message.FromGroupID
	mess.CurrentQQ = args.CurrentQQ
	mess.FromSourceID = message.FromUin
	OnFriendMsgsFunc(mess)
}
