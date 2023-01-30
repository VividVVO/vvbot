package plugin

import (
	"fmt"
	"github.com/vivid-vvo/vvbot/app/model/clan_group"
	"github.com/vivid-vvo/vvbot/qqbot/model/bot"
	"github.com/vivid-vvo/vvbot/qqbot/plugins"
)

type Gonghuiliebiao struct {
}

func (l Gonghuiliebiao) CheckOrder(cm string) bool {
	return true
}

func (l Gonghuiliebiao) IsNotCheckOrder() bool {
	return false
}

func (l Gonghuiliebiao) GetOrders() []string {
	return []string{
		"^(公|工|行)会列表$",
		"^查询(公|工|行)会(列表)?$",
		"^查看(公|工|行)会(列表)?$",
	}
}

func (l Gonghuiliebiao) Run(mess plugins.MeassageData, ms string, atqq int64) {
	clanGroupList, err := clan_group.GetAllClanGroup()
	if err != nil {
		bot.Send(mess.FromGroupID, mess.SendToType, err.Error())
		return
	}
	if clanGroupList == nil {
		bot.Send(mess.FromGroupID, mess.SendToType, "当前未创建公会，请输入”创建[日台韩国]服公会 公会名“ 例：“创建国服公会 拉胯会长拉胯记”")
		return
	}
	msg := "当前公会列表如下：\n"
	for i, clanGroup := range clanGroupList {
		msg += fmt.Sprintf(" %d.%s\n", i+1, clanGroup.GroupName)
	}
	bot.Send(mess.FromGroupID, mess.SendToType, msg)
}
