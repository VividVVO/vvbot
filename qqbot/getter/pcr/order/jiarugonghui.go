package order

import (
	"fmt"
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/util/grand"
	"github.com/vivid-vvo/vvbot/app/model/clan_member"
	"github.com/vivid-vvo/vvbot/app/model/user"
	"github.com/vivid-vvo/vvbot/qqbot/getter"
	"github.com/vivid-vvo/vvbot/qqbot/model/bot"
	"github.com/vivid-vvo/vvbot/qqbot/model/pcr"
)

type Jiarugonghui struct {
}

func (l Jiarugonghui) CheckOrder(cm string) bool {
	return true
}

func (l Jiarugonghui) IsNotCheckOrder() bool {
	return false
}

func (l Jiarugonghui) GetOrders() []string {
	return []string{
		"^加入(公|工|行)会( +)?$",
		"^进入(公|工|行)会( +)?$",
	}

}

func (l Jiarugonghui) Run(mess getter.MeassageData, ms string, atqq int) {
	qqid := int(mess.FromUserID)
	if atqq != 0 {
		qqid = atqq
	}
	clanGroupData, err := pcr.GetClanGroupAndChack(mess.FromGroupID)
	if err != nil {
		bot.Send(mess.FromGroupID, 2, err.Error())
		return
	}
	userInfo, err := user.GetProfile(qqid)
	if err != nil {
		bot.Send(mess.FromGroupID, 2, err.Error())
		return
	}
	if userInfo == nil {
		// 注册
		err := user.SignUp(qqid, grand.Letters(12), mess.FromNickName, "")
		if err != nil {
			bot.Send(mess.FromGroupID, 2, err.Error())
			return
		}
	}
	clanMember, err := clan_member.GetClanMember(qqid, clanGroupData.GroupId)
	if err != nil {
		bot.Send(mess.FromGroupID, 2, err.Error())
		return
	}
	var gameName string
	if clanMember != nil {
		gameName = clanMember.GameName
	}

	if clanMember != nil {
		if clanMember.GroupId == clanGroupData.GroupId {
			bot.Send(mess.FromGroupID, 2, fmt.Sprintf("%s 已在当前公会", gameName))
			return
		}
	}
	member := bot.GetGourpUserData(mess.FromGroupID, qqid)
	if member == nil {
		bot.Send(mess.FromGroupID, 2, "内部错误")
		return
	}
	var name string
	if member.GroupCard != "" {
		name = member.GroupCard
	} else if member.AutoRemark != "" {
		name = member.AutoRemark
	} else if member.NickName != "" {
		name = member.NickName
	}
	entity := new(clan_member.Entity)
	entity.GroupId = clanGroupData.GroupId
	entity.Qqid = qqid
	entity.GameName = name
	entity.JoinTime = int(gtime.Now().Unix())
	err = clan_member.JoinClan(entity)
	if err != nil {
		bot.Send(mess.FromGroupID, 2, fmt.Sprintf("内部错误"))
		return
	}
	if err := user.ChangeClanGroupId(qqid, entity.GroupId); err != nil {
		bot.Send(mess.FromGroupID, 2, err.Error())
		return
	}
	bot.Send(mess.FromGroupID, 2, fmt.Sprintf("%s 加入公会<%s>成功！", name, clanGroupData.GroupName))
}
