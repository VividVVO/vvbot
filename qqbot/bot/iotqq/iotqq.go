package iotqq

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/glog"
	"github.com/gogf/gf/util/gconv"
	gosocketio "github.com/graarh/golang-socketio"
	"github.com/graarh/golang-socketio/transport"
	"github.com/vivid-vvo/vvbot/qqbot/getter"
	"log"
	"strconv"
	"time"
)

var (
	site   = g.Cfg().GetString("iotbot.site")
	port   = g.Cfg().GetInt("iotbot.port")
	qq     = g.Cfg().GetString("iotbot.qq")
	domain string
	c      *gosocketio.Client
)
var mSendChan chan MeassageSend

type MeassageSend struct {
	FromGroupID  int
	SendToType   int
	PicBase64Buf string
	PicUrl       string
	Content      string
	Qqid         int
}

type Iotqq struct{}

var (
	OnGroupMsgsFunc func(getter.MeassageData)
)

// 自动连接/重连QQ机器人连接
func autoJoin(qq string) {
	for {
		glog.Println("获取QQ号连接")
		c, err := gosocketio.Dial(
			gosocketio.GetUrl(site, port, false),
			transport.GetDefaultWebsocketTransport())
		if err != nil {
			log.Fatal(err)
		}
		err = c.On("OnGroupMsgs", OnGroupMsgs)
		if err != nil {
			log.Fatal(err)
		}
		time.Sleep(1 * time.Second)
		result, err := c.Ack("GetWebConn", qq, time.Second*5)

		if err != nil {
			glog.Fatal("qq连接失败", err)
		} else {
			glog.Println("emit", result)
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
		if len(mSend.PicBase64Buf) != 0 {
			go sendPic(mSend.FromGroupID, mSend.SendToType, mSend.Content, mSend.PicBase64Buf, mSend.PicUrl)
		} else {
			go send(mSend.FromGroupID, mSend.SendToType, mSend.Content, mSend.Qqid)
		}
		time.Sleep(1000 * time.Millisecond)
	}
}

// 绑定群信息接收func
func (iotqq *Iotqq) BindOnGroupMsgs(func1 func(getter.MeassageData)) {
	OnGroupMsgsFunc = func1
}

func (iotqq *Iotqq) Init() {
	domain = site + ":" + strconv.Itoa(port)
	go autoJoin(qq)
	go msgSendWorker()
}
func (iotqq *Iotqq) GetGourpUserData(qqGroupId int, qqid int) *getter.Member {
	var member2 *getter.Member
	member := getGourpUserData(qqGroupId, qqid)
	if err := gconv.Struct(member, &member2); err != nil {
		glog.Error(err)
		return nil
	}
	return member2

}
func (iotqq *Iotqq) GetGroupUserList(qqGroupId int, lastUin int) (*getter.MemberData, error) {
	var qqGroupMemberData2 *getter.MemberData
	qqGroupMemberData, err := getGroupUserList(qqGroupId, lastUin)
	if err != nil {
		return nil, err
	}
	if err := gconv.Struct(qqGroupMemberData, &qqGroupMemberData2); err != nil {
		glog.Error(err)
		return nil, err
	}
	return qqGroupMemberData2, nil
}
func (iotqq *Iotqq) GetBotQQID() string {
	return getQQ()
}
func (iotqq *Iotqq) Send(ToUser int, SendToType int, Content string) {
	mSendChan <- MeassageSend{
		FromGroupID: ToUser,
		SendToType:  SendToType,
		Content:     Content,
	}
}
func (iotqq *Iotqq) SendPic(ToUser int, SendToType int, Content string, PicBase64Buf string, PicUrl string) {
	mSendChan <- MeassageSend{
		FromGroupID:  ToUser,
		SendToType:   SendToType,
		Content:      Content,
		PicBase64Buf: PicBase64Buf,
		PicUrl:       PicUrl,
	}
}
func (iotqq *Iotqq) SendAtQq(ToUser int, SendToType int, Content string, Qqid int) {
	mSendChan <- MeassageSend{
		FromGroupID: ToUser,
		SendToType:  SendToType,
		Content:     Content,
		Qqid:        Qqid,
	}
}
