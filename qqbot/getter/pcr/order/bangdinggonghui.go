package order

import (
	"fmt"
	"github.com/vivid-vvo/vvbot/app/model/clan_group"
	"github.com/vivid-vvo/vvbot/app/model/user"
	"github.com/vivid-vvo/vvbot/qqbot/getter"
	"github.com/vivid-vvo/vvbot/qqbot/model/bot"
	"strings"
)

type Bangdinggonghui struct {
}

func (l Bangdinggonghui) CheckOrder(cm string) bool {
	return true
}

func (l Bangdinggonghui) IsNotCheckOrder() bool {
	return false
}

func (l Bangdinggonghui) GetOrders() []string {
	return []string{
		"^绑定(公|行)会",
		"^(公|行)会绑定",
	}

}

func (l Bangdinggonghui) Run(mess getter.MeassageData, ms string, atqq int) {
	if !user.CheckUserAuthorityGroup(int(mess.FromUserID), user.AuthAdmin) {
		bot.Send(mess.FromGroupID, 2, "权限不足")
		return
	}
	cms := strings.Split(ms, " ")
	if len(cms) < 2 {
		bot.Send(mess.FromGroupID, 2, "绑定公会格式应为：”绑定公会 公会名“ 例：“绑定公会 拉胯会长拉胯记”")
		return
	}
	clanGroupData, err := clan_group.GetClanGroupAtQqGroupId(mess.FromGroupID)
	if err != nil {
		bot.Send(mess.FromGroupID, 2, err.Error())
		return
	}
	if clanGroupData != nil {
		if clanGroupData.GroupName == cms[1] {
			bot.Send(mess.FromGroupID, 2, fmt.Sprintf("当前已绑定到此公会，无法重复绑定"))
			return
		}
		bot.Send(mess.FromGroupID, 2, fmt.Sprintf("当前公会已绑定公会 <%s> 如需绑定到此公会，请先“解绑公会”", clanGroupData.GroupName))
		return
	}

	clanGroupData, err = clan_group.GetClanGroupAtName(cms[1])
	if err != nil {
		bot.Send(mess.FromGroupID, 2, err.Error())
		return
	}
	if clanGroupData == nil {
		bot.Send(mess.FromGroupID, 2, fmt.Sprintf("指定公会不存在，请先创建公会"))
		return
	}
	// 绑定公会
	err = clan_group.BindGroup(clanGroupData.GroupId, mess.FromGroupID)
	if err != nil {
		bot.Send(mess.FromGroupID, 2, err.Error())
		return
	}

	bot.Send(mess.FromGroupID, 2, "公会绑定成功！公会战成员请发送“加入公会”，或发送“加入全部成员”")
}
