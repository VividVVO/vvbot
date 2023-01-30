package plugin

import (
	"fmt"
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/util/grand"
	"github.com/vivid-vvo/vvbot/app/model/clan_group"
	"github.com/vivid-vvo/vvbot/app/service/check"
	"github.com/vivid-vvo/vvbot/library/tools"
	"github.com/vivid-vvo/vvbot/qqbot/model/bot"
	"github.com/vivid-vvo/vvbot/qqbot/plugins"
	"strings"
)

type Chaungjiangonghui struct {
}

func (l Chaungjiangonghui) CheckOrder(cm string) bool {
	return true
}

func (l Chaungjiangonghui) IsNotCheckOrder() bool {
	return false
}

func (l Chaungjiangonghui) GetOrders() []string {
	return []string{
		"^创建(新)?(日|台|韩|国)?(服)?(公|工|行)会",
		"^建立(新)?(日|台|韩|国)?(服)?(公|工|行)会",
		"^成立(新)?(日|台|韩|国)?(服)?(公|工|行)会",
		"^创立(新)?(日|台|韩|国)?(服)?(公|工|行)会",
	}

}

// !check.CheckAuthorityGroup(mess.
func (l Chaungjiangonghui) Run(mess plugins.MeassageData, ms string, atqq int64) {
	qqid := mess.FromUserID
	var gameServer, groupName string
	if !check.CheckAuthorityGroup(mess.FromUserID, check.AuthAdmin, 0) {
		bot.Send(mess.FromGroupID, mess.SendToType, "权限不足")
		return
	}
	cms := strings.Split(ms, " ")

	if len(cms) < 2 {
		bot.Send(mess.FromGroupID, mess.SendToType, "创建公会格式应为：”创建[日台韩国]服公会 公会名“ 例：“创建国服公会 拉胯会长拉胯记”")
		return

	}
	gameServer = string([]rune(cms[0])[2:3])

	if tools.Compare(gameServer, "国") {
		gameServer = "CN"
	} else if tools.Compare(gameServer, "韩") {
		gameServer = "KR"
	} else if tools.Compare(gameServer, "日") {
		gameServer = "JP"
	} else if tools.Compare(gameServer, "台") {
		gameServer = "TW"
	} else {
		bot.Send(mess.FromGroupID, mess.SendToType, "创建公会格式应为：”创建[日台韩国]服公会 公会名“")
		return
	}
	groupName = cms[1]
	if groupName == "" {
		bot.Send(mess.FromGroupID, mess.SendToType, "创建公会格式应为：”创建[日台韩国]服公会 公会名“")
		return
	}

	clanGroupData, err := clan_group.GetClanGroupAtQqGroupId(mess.FromGroupID)
	if err != nil {
		bot.Send(mess.FromGroupID, mess.SendToType, err.Error())
		return
	}
	if clanGroupData != nil {
		bot.Send(mess.FromGroupID, mess.SendToType, fmt.Sprintf("当前公会已绑定公会 <%s> 如需创建新公会，请先“解绑公会”", clanGroupData.GroupName))
		return
	}

	clanGroupData, err = clan_group.GetClanGroupAtName(groupName)
	if err != nil {
		bot.Send(mess.FromGroupID, mess.SendToType, err.Error())
		return
	}
	if clanGroupData != nil {
		bot.Send(mess.FromGroupID, mess.SendToType, fmt.Sprintf("公会已存在，请使用“绑定公会 %s”来绑定到此公会", clanGroupData.GroupName))
		return
	}

	entity := new(clan_group.Entity)
	entity.GameServer = gameServer
	entity.GroupName = groupName
	entity.CreatorQqid = qqid
	entity.CreateTime = gtime.Now().Unix()
	entity.GroupId = grand.N(100000000, 1000000000)
	entity.BindQqGroup = mess.FromGroupID
	if err := clan_group.ClanGroupCreate(entity); err != nil {
		bot.Send(mess.FromGroupID, mess.SendToType, err.Error())
		return
	}

	bot.Send(mess.FromGroupID, mess.SendToType, "公会创建成功，请登录后台查看，公会战成员请发送“加入公会”，或发送“加入全部成员”")
}
