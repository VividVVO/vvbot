package bot

import (
	"fmt"
	"github.com/gogf/gf/os/glog"
	"github.com/vivid-vvo/vvbot/qqbot/plugins"
	"strconv"
)

type Bot interface {
	Init()
	GetBotQQID() string
	LeaveChat(qqGroupId int64, sendToType int, dismiss bool) error
	GetGroupUserData(qqGroupId int64, qqid int64) *plugins.Member
	GetGroupUserList(qqGroupId int64) (*plugins.MemberData, error)
	GetAtQQStr(qqID int64) string
	// Send SendToType = 1 私聊 SendToType = 2 群聊
	Send(ToUser int64, SendToType int, Content string)
	SendPrivate(ToUser int64, Content string, Qqid int64)
	SendPic(ToUser int64, SendToType int, Content string, PicBase64Buf string, PicUrl string)
	SendAtQq(ToUser int64, SendToType int, Content string, Qqid int64)
	// BindOnGroupMsgs 绑定群聊事件
	BindOnGroupMsgs(func1 func(plugins.MeassageData))
	// BindOnFriendMsgs 绑定私聊事件
	BindOnFriendMsgs(func1 func(plugins.MeassageData))
}

var botEntity Bot

func Init(bot Bot) {
	botEntity = bot
	botEntity.Init()
}

func LeaveChat(qqGroupId int64, sendToType int, dismiss bool) error {
	return botEntity.LeaveChat(qqGroupId, sendToType, dismiss)
}
func GetGroupUserData(qqGroupId int64, qqid int64) *plugins.Member {
	return botEntity.GetGroupUserData(qqGroupId, qqid)
}
func GetGroupUserList(qqGroupId int64) (*plugins.MemberData, error) {
	return botEntity.GetGroupUserList(qqGroupId)
}
func GetBotQQID() string {
	return botEntity.GetBotQQID()
}
func GetBotQQIDToInt64() int64 {
	qqID, _ := strconv.ParseInt(botEntity.GetBotQQID(), 10, 64)
	return qqID
}

func GetAtQQStr(qqID int64) string {
	return botEntity.GetAtQQStr(qqID)
}
func Send(ToUser int64, SendToType int, Content string) {
	glog.Println(fmt.Sprintf("bot.SendMessage(%d, %d, %s)", ToUser, SendToType, Content))
	botEntity.Send(ToUser, SendToType, Content)
}

func SendPrivate(ToUser int64, Content string, Qqid int64) {
	glog.Println(fmt.Sprintf("bot.SendMessage(%d, %s)", ToUser, Content))
	botEntity.SendPrivate(ToUser, Content, Qqid)
}

func SendPic(ToUser int64, SendToType int, Content string, PicBase64Buf string, PicUrl string) {
	botEntity.SendPic(ToUser, SendToType, Content, PicBase64Buf, PicUrl)
}
func SendAtQq(ToUser int64, SendToType int, Content string, Qqid int64) {
	botEntity.SendAtQq(ToUser, SendToType, Content, Qqid)
}
