package plugin

import (
	"fmt"
	"github.com/vivid-vvo/vvbot/app/model/clan_group"
	"github.com/vivid-vvo/vvbot/qqbot/model/bot"
	"github.com/vivid-vvo/vvbot/qqbot/plugins"
)

type Jiebanggonghui struct {
}

func (l Jiebanggonghui) CheckOrder(cm string) bool {
	return true
}

func (l Jiebanggonghui) IsNotCheckOrder() bool {
	return false
}

func (l Jiebanggonghui) GetOrders() []string {
	return []string{
		"^解除(公|工|行)会$",
		"^(公|工|行)会解除$",
		"^解绑(公|工|行)会$",
		"^(公|工|行)会解绑$",
		"^解除绑定$",
		"^绑定解除$"}
}

func (l Jiebanggonghui) Run(mess plugins.MeassageData, ms string, atqq int64) {
	// qqid := mess.FromUserID
	clanGroupData, err := clan_group.GetClanGroupAtQqGroupId(mess.FromGroupID)
	if err != nil {
		bot.Send(mess.FromGroupID, mess.SendToType, err.Error())
		return
	}
	if clanGroupData == nil {
		bot.Send(mess.FromGroupID, mess.SendToType, fmt.Sprintf("当前公会未绑定公会，无需解绑"))
		return
	}

	// 解绑公会
	err = clan_group.UnBindGroup(mess.FromGroupID)
	if err != nil {
		bot.Send(mess.FromGroupID, mess.SendToType, err.Error())
		return
	}

	// 退出组内成员
	/*err = clan_member.MemberExitGroup(clanGroupData.GroupId)
	if err != nil {
		bot.Send(mess.FromGroupID, mess.SendToType, err.Error())
		return
	}*/
	bot.Send(mess.FromGroupID, mess.SendToType, "公会解绑成功！请使用“创建[日台韩国]服公会 公会名”创建一个新的公会，或者使用“绑定公会 公会名”绑定一个现有的公会")
}
