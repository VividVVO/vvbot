package plugin

import (
	"fmt"
	"github.com/vivid-vvo/vvbot/app/model/clan_group"
	"github.com/vivid-vvo/vvbot/app/model/clan_member"
	"github.com/vivid-vvo/vvbot/app/model/user"
	"github.com/vivid-vvo/vvbot/qqbot/model/bot"
	"github.com/vivid-vvo/vvbot/qqbot/plugins"
)

type Tuichugonghui struct {
}

func (l Tuichugonghui) CheckOrder(cm string) bool {
	return true
}

func (l Tuichugonghui) IsNotCheckOrder() bool {
	return false
}

func (l Tuichugonghui) GetOrders() []string {
	return []string{
		"^退出(公|行)会( +)?$",
		"^(公|行)会退出( +)?$",
	}

}
func (l Tuichugonghui) Run(mess plugins.MeassageData, cm string, atqq int64) {
	qqid := mess.FromUserID
	if atqq != 0 {
		qqid = atqq
	}
	clanGroup, err := clan_group.GetClanGroupAtQqGroupId(mess.FromGroupID)
	if err != nil {
		bot.Send(mess.FromGroupID, mess.SendToType, err.Error())
		return
	}
	if clanGroup == nil {
		bot.Send(mess.FromGroupID, mess.SendToType, "退出失败，本群当前无绑定公会")
		return
	}

	userClanGroupData, err := clan_member.GetClanMember(qqid, clanGroup.GroupId)
	if err != nil {
		bot.Send(mess.FromGroupID, mess.SendToType, err.Error())
		return
	}
	if userClanGroupData == nil {
		bot.Send(mess.FromGroupID, mess.SendToType, "您未加入本公会，无需退出")
		return
	}
	err = clan_member.MemberExitGroupAtQqid(qqid, userClanGroupData.GroupId)
	if err != nil {
		bot.Send(mess.FromGroupID, mess.SendToType, fmt.Sprintf("内部错误"))
		return
	}
	if err := user.ChangeClanGroupId(qqid, 0); err != nil {
		bot.Send(mess.FromGroupID, mess.SendToType, err.Error())
		return
	}
	bot.Send(mess.FromGroupID, mess.SendToType, fmt.Sprintf("%s 退出公会成功", userClanGroupData.GameName))
}
