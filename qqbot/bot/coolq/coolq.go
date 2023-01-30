package coolq

import (
	"fmt"
	qqbotapi "github.com/catsworld/qq-bot-api"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/glog"
	"github.com/vivid-vvo/vvbot/qqbot/plugins"
	"log"
	"net/http"
	"strconv"
)

var (
	coolqUrl      = g.Cfg().GetString("coolq.CoolqUrl")
	accessToken   = g.Cfg().GetString("coolq.AccessToken")
	secret        = g.Cfg().GetString("coolq.Secret")
	wsReverseHost = g.Cfg().GetString("coolq.WsReverseHost")
	wsReversePort = g.Cfg().GetInt("coolq.WsReversePort")
	wsReversePath = g.Cfg().GetString("coolq.WsReversePath")

	domain string
)

type MeassageSend struct {
	FromGroupID  int
	SendToType   int
	PicBase64Buf string
	PicUrl       string
	Content      string
	Qqid         int64
}

type Coolq struct{}

var (
	OnGroupMsgsFunc  func(plugins.MeassageData)
	OnFriendMsgsFunc func(plugins.MeassageData)
	coolbot          *qqbotapi.BotAPI
)

// 绑定群信息接收func
func (coolq *Coolq) BindOnGroupMsgs(func1 func(plugins.MeassageData)) {
	OnGroupMsgsFunc = func1
}

// 绑定群信息接收func
func (coolq *Coolq) BindOnFriendMsgs(func1 func(plugins.MeassageData)) {
	OnFriendMsgsFunc = func1
}
func (coolq *Coolq) Init() {
	domain = wsReverseHost + ":" + strconv.Itoa(wsReversePort)
	go qqbotConn()
}

func qqbotConn() {
	var err error
	coolbot, err = qqbotapi.NewBotAPI(accessToken, coolqUrl, secret)
	if err != nil {
		log.Fatal(err)
	}
	// bot.Debug = true
	u := qqbotapi.NewWebhook("/" + wsReversePath)
	u.PreloadUserInfo = true
	// Use WebHook as event method
	// updates := bot.ListenForWebhook(u)
	// Or if you love WebSocket Reverse
	updates := coolbot.ListenForWebSocket(u)
	go http.ListenAndServe(domain, nil)
	for update := range updates {
		if update.Message == nil {
			continue
		}
		// group
		// private
		// discuss
		log.Printf("[%s] %s [%s]", update.Message.From.String(), update.Message.Text, update.Message.Chat.Type)
		switch update.Message.Chat.Type {
		case "group":
			go OnGroupMsgs(coolbot, update)
		case "private":
			go OnFriendMsgs(coolbot, update)
		}

		// update.Message.Chat.IsGroup()
		//coolbot.SendMessage(update.Message.Chat.ID, update.Message.Chat.Type, update.Message.Text)
	}
}

func (coolq *Coolq) LeaveChat(qqGroupId int64, sendToType int, dismiss bool) error {
	var chatType string
	// "private"、"group"、"discuss"
	switch sendToType {
	case 1:
		chatType = "private"
	case 2:
		chatType = "group"
	}
	_, err := coolbot.LeaveChat(qqGroupId, chatType, dismiss)
	return err
}

func (coolq *Coolq) GetGroupUserData(qqGroupId int64, qqid int64) *plugins.Member {
	member, err := coolbot.GetGroupMemberInfo(qqGroupId, qqid, true)

	if err != nil {
		glog.Error(err)
		return nil
	}
	data2 := new(plugins.Member)
	data2.Qqid = member.ID
	data2.NickName = member.NickName
	data2.JoinTime = member.JoinTimeUnix
	data2.GroupCard = member.Card
	data2.Role = member.Role
	data2.Age = member.Age
	data2.Sex = member.Sex
	return data2

}
func (coolq *Coolq) GetGroupUserList(qqGroupId int64) (*plugins.MemberData, error) {
	qqGroupMemberData2 := new(plugins.MemberData)
	qqGroupMemberList, err := coolbot.GetGroupMemberList(qqGroupId)
	if err != nil {
		return nil, err
	}
	qqGroupMemberData2.Count = len(qqGroupMemberList)
	var data2 plugins.Member
	for _, data := range qqGroupMemberList {
		data2.Qqid = data.ID
		data2.NickName = data.NickName
		data2.JoinTime = data.JoinTimeUnix
		data2.GroupCard = data.Card
		qqGroupMemberData2.MemberList = append(qqGroupMemberData2.MemberList, data2)
	}
	return qqGroupMemberData2, nil
}
func (coolq *Coolq) GetBotQQID() string {
	return strconv.FormatInt(coolbot.Self.ID, 10)
}
func (coolq *Coolq) GetAtQQStr(qqid int64) string {
	return fmt.Sprintf("[CQ:at,qq=%d]", qqid)
}

// Send SendToType = 1 私聊 endToType = 2 群聊
func (coolq *Coolq) Send(ToUser int64, SendToType int, Content string) {
	var chatType string
	// "private"、"group"、"discuss"
	switch SendToType {
	case 1:
		chatType = "private"
	case 2:
		chatType = "group"
	}
	coolbot.SendMessage(ToUser, chatType, Content)
}
func (coolq *Coolq) SendPic(ToUser int64, SendToType int, Content string, PicBase64Buf string, PicUrl string) {
	// coolbot.Send()
}
func (coolq *Coolq) SendAtQq(ToUser int64, SendToType int, Content string, Qqid int64) {

}

// SendPrivate 发送私聊
func (coolq *Coolq) SendPrivate(ToUser int64, Content string, Qqid int64) {
	var chatType string
	// "private"、"group"、"discuss"

	chatType = "private"

	coolbot.SendMessage(ToUser, chatType, Content)
}
