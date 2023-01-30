package iotqq

import (
	"errors"
	"fmt"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/glog"
	"github.com/gogf/gf/util/gconv"
	gosocketio "github.com/graarh/golang-socketio"
	"github.com/graarh/golang-socketio/transport"
	"github.com/vivid-vvo/vvbot/qqbot/plugins"
	"strconv"
	"strings"
	"time"
)

var (
	site   = g.Cfg().GetString("iotqq.Site")
	port   = g.Cfg().GetInt("iotqq.Port")
	qq     = g.Cfg().GetString("iotqq.QQ")
	domain string
	c      *gosocketio.Client
)
var mSendChan chan MeassageSend

type MeassageSend struct {
	FromGroupID  int64
	SendToType   int
	PicBase64Buf string
	PicUrl       string
	Content      string
	Qqid         int64
	PriQQid      int64
}

type Iotqq struct{}

var (
	OnGroupMsgsFunc  func(plugins.MeassageData)
	OnFriendMsgsFunc func(plugins.MeassageData)
)

// 自动连接/重连QQ机器人连接
func autoJoin(qq string) {
	for {
		glog.Println("获取QQ号连接")
		c, err := gosocketio.Dial(
			gosocketio.GetUrl(site, port, false),
			transport.GetDefaultWebsocketTransport())
		if err != nil {
			glog.Error(err)
			time.Sleep(3 * time.Second)
			continue
		}
		err = c.On("OnGroupMsgs", OnGroupMsgs)
		if err != nil {
			glog.Error(err)
			time.Sleep(3 * time.Second)
			continue
		}
		err = c.On("OnFriendMsgs", OnFriendMsgs)
		if err != nil {
			glog.Error(err)
			time.Sleep(3 * time.Second)
			continue
		}

		for {
			result, err := c.Ack("GetWebConn", qq, time.Second*5)
			if err != nil {
				glog.Error("qq连接失败", err)
				break
			} else {
				glog.Println("emit", result)
				if strings.Contains(result, "OK") {
					break
				}
			}
			time.Sleep(3 * time.Second)
		}
		for c.IsAlive() {
			time.Sleep(3 * time.Second)
		}
		time.Sleep(10 * time.Second)
	}
}
func msgSendWorker() {
	mSendChan = make(chan MeassageSend, 100)
	for mSend := range mSendChan {
		if mSend.PriQQid != 0 {
			go SendPrivate(mSend.FromGroupID, mSend.Content, mSend.PriQQid)
		} else if len(mSend.PicBase64Buf) != 0 {
			go sendPic(mSend.FromGroupID, mSend.SendToType, mSend.Content, mSend.PicBase64Buf, mSend.PicUrl)
		} else {
			go send(mSend.FromGroupID, mSend.SendToType, mSend.Content, mSend.Qqid)
		}
		time.Sleep(1000 * time.Millisecond)
	}
}

// 绑定群信息接收func
func (iotqq *Iotqq) BindOnGroupMsgs(func1 func(plugins.MeassageData)) {
	OnGroupMsgsFunc = func1
}

// 绑定群信息接收func
func (iotqq *Iotqq) BindOnFriendMsgs(func1 func(plugins.MeassageData)) {
	OnFriendMsgsFunc = func1
}

func (iotqq *Iotqq) Init() {
	domain = site + ":" + strconv.Itoa(port)
	go autoJoin(qq)
	go msgSendWorker()
}

func (coolq *Iotqq) LeaveChat(qqGroupId int64, sendToType int, dismiss bool) error {
	return errors.New("iotqq无此功能")
}

func (iotqq *Iotqq) GetAtQQStr(qqid int64) string {
	return fmt.Sprintf("[ATUSER(%d)]", qqid)
}
func (iotqq *Iotqq) GetGroupUserData(qqGroupId int64, qqid int64) *plugins.Member {
	var member2 *plugins.Member
	member := getGroupUserData(qqGroupId, qqid)
	if err := gconv.Struct(member, &member2); err != nil {
		glog.Error(err)
		return nil
	}
	qqGroup := getGroupData(qqGroupId)
	if member.GroupAdmin == 1 {
		member2.Role = "admin"
	} else if qqGroup != nil && qqGroup.GroupOwner == member.Qqid {
		member2.Role = "owner"
	} else {
		member2.Role = "member"
	}
	return member2

}
func (iotqq *Iotqq) GetGroupUserList(qqGroupId int64) (*plugins.MemberData, error) {
	var qqGroupMemberData2 *plugins.MemberData
	qqGroupMemberData, err := getGroupUserList(qqGroupId, 0)
	if err != nil {
		return nil, err
	}
	if err := gconv.Struct(qqGroupMemberData, &qqGroupMemberData2); err != nil {
		glog.Error(err)
		return nil, err
	}
	qqGroup := getGroupData(qqGroupId)
	if qqGroup == nil {
		glog.Error(errors.New("获取群数据失败"))
		return nil, errors.New("获取群数据失败")
	}
	for i, member := range qqGroupMemberData.MemberList {
		if member.Qqid == qqGroup.GroupOwner {
			if member.GroupAdmin == 1 {
				qqGroupMemberData2.MemberList[i].Role = "admin"
			} else if qqGroup != nil && qqGroup.GroupOwner == member.Qqid {
				qqGroupMemberData2.MemberList[i].Role = "owner"
			} else {
				qqGroupMemberData2.MemberList[i].Role = "member"
			}
		}
	}

	return qqGroupMemberData2, nil
}
func (iotqq *Iotqq) GetBotQQID() string {
	return getQQ()
}
func (iotqq *Iotqq) Send(ToUser int64, sendToType int, Content string) {
	mSendChan <- MeassageSend{
		FromGroupID: ToUser,
		SendToType:  sendToType,
		Content:     Content,
	}
}

func (iotqq *Iotqq) SendPrivate(ToUser int64, Content string, Qqid int64) {
	mSendChan <- MeassageSend{
		FromGroupID: ToUser,
		SendToType:  1,
		Content:     Content,
		PriQQid:     Qqid,
	}
}

func (iotqq *Iotqq) SendPic(ToUser int64, sendToType int, Content string, PicBase64Buf string, PicUrl string) {
	mSendChan <- MeassageSend{
		FromGroupID:  ToUser,
		SendToType:   sendToType,
		Content:      Content,
		PicBase64Buf: PicBase64Buf,
		PicUrl:       PicUrl,
	}
}
func (iotqq *Iotqq) SendAtQq(ToUser int64, sendToType int, Content string, Qqid int64) {
	mSendChan <- MeassageSend{
		FromGroupID: ToUser,
		SendToType:  sendToType,
		Content:     Content,
		Qqid:        Qqid,
	}
}
