package iotqq

import (
	"encoding/json"
	"fmt"
	"github.com/gogf/gf/os/glog"
	"github.com/gogf/gf/util/gconv"
	gosocketio "github.com/graarh/golang-socketio"
	"github.com/vivid-vvo/vvbot/qqbot/getter"
	"log"
	"regexp"
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
	MsgTime       int         `json:"MsgTime"`
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
	case "PicMsg":
		err := json.Unmarshal([]byte(message.Content), &content)
		if err != nil {
			log.Fatalln(err)
			return
		}
		msg = content.Content
		if len(content.UserID) > 0 {
			atqq = content.UserID
		}
	default:
		return
	}
	var mess getter.MeassageData
	if err := gconv.Struct(args.CurrentPacket.Data, &mess); err != nil {
		glog.Error(err)
		return
	}
	var name string
	// 清除文本内的@信息
	if len(atqq) > 0 {
		member := getGourpUserData(message.FromGroupID, int(atqq[0]))
		// 群成员获取错误
		if member == nil {
			glog.Error("iotqq.GetGourpUserData")
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
		regexp1 := regexp.MustCompile(fmt.Sprintf("@%s( +)?", name))
		msg = regexp1.ReplaceAllString(msg, "")
	}
	mess.Content = msg
	mess.AtQQList = atqq
	mess.CurrentQQ = args.CurrentQQ
	OnGroupMsgsFunc(mess)
}
