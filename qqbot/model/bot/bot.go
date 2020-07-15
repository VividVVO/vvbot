package bot

import (
	"github.com/vivid-vvo/vvbot/qqbot/getter"
)

type Bot interface {
	Init()
	GetBotQQID() string
	GetGourpUserData(qqGroupId int, qqid int) *getter.Member
	GetGroupUserList(qqGroupId int, lastUin int) (*getter.MemberData, error)
	Send(ToUser int, SendToType int, Content string)
	SendPic(ToUser int, SendToType int, Content string, PicBase64Buf string, PicUrl string)
	SendAtQq(ToUser int, SendToType int, Content string, Qqid int)
	BindOnGroupMsgs(func1 func(getter.MeassageData))
}

var botEntity Bot

func Init(bot Bot) {
	botEntity = bot
	botEntity.Init()
}

func GetGourpUserData(qqGroupId int, qqid int) *getter.Member {
	return botEntity.GetGourpUserData(qqGroupId, qqid)
}
func GetGroupUserList(qqGroupId int, lastUin int) (*getter.MemberData, error) {
	return botEntity.GetGroupUserList(qqGroupId, lastUin)
}
func GetBotQQID() string {
	return botEntity.GetBotQQID()
}
func Send(ToUser int, SendToType int, Content string) {
	botEntity.Send(ToUser, SendToType, Content)
}
func SendPic(ToUser int, SendToType int, Content string, PicBase64Buf string, PicUrl string) {
	botEntity.SendPic(ToUser, SendToType, Content, PicBase64Buf, PicUrl)
}
func SendAtQq(ToUser int, SendToType int, Content string, Qqid int) {
	botEntity.SendAtQq(ToUser, SendToType, Content, Qqid)
}
