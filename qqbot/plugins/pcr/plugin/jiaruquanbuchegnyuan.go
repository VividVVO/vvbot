package plugin

import (
	"fmt"
	"github.com/vivid-vvo/vvbot/app/model/clan_group"
	"github.com/vivid-vvo/vvbot/qqbot/model/bot"
	"github.com/vivid-vvo/vvbot/qqbot/model/pcr"
	"github.com/vivid-vvo/vvbot/qqbot/plugins"
	"strconv"
)

type Jiaruquanbuchengyuan struct {
}

func (l Jiaruquanbuchengyuan) CheckOrder(cm string) bool {
	return true
}

func (l Jiaruquanbuchengyuan) IsNotCheckOrder() bool {
	return false
}

func (l Jiaruquanbuchengyuan) GetOrders() []string {
	return []string{
		"^加入全部成员$",
		"^全部成员加入公会$",
		"^所有人加入公会$",
		"^全部加入公会$",
	}

}

func (l Jiaruquanbuchengyuan) Run(mess plugins.MeassageData, ms string, atqq int64) {
	groupUser, err := bot.GetGroupUserList(mess.FromGroupID)
	if err != nil {
		bot.Send(mess.FromGroupID, mess.SendToType, "群成员获取失败，请重新尝试")
		return
	}
	if groupUser.Count > 40 {
		bot.Send(mess.FromGroupID, mess.SendToType, "群成员数量大于40人，无法加入全部成员！")
		return
	}

	var clanGroup *clan_group.Entity
	clanGroup, err = pcr.GetClanGroupAndChack(mess.FromGroupID)
	if err != nil {
		bot.Send(mess.FromGroupID, mess.SendToType, err.Error())
		return
	}
	var errMsg string
	var addNum int
	mainQQ, _ := strconv.ParseInt(bot.GetBotQQID(), 10, 64)
	for _, data := range groupUser.MemberList {
		if data.Qqid != mainQQ {
			var name string
			if data.GroupCard != "" {
				name = data.GroupCard
			} else if data.AutoRemark != "" {
				name = data.AutoRemark
			} else if data.NickName != "" {
				name = data.NickName
			}

			err = pcr.UserJoinGroup(mess.FromGroupID, data.Qqid, name)
			if err != nil {
				errMsg += fmt.Sprintf("qq: %d（%s）, 加入失败 -> %s\n", data.Qqid, name, err.Error())
			} else {
				addNum++
			}
		}
	}
	if errMsg != "" {
		bot.Send(mess.FromGroupID, mess.SendToType, errMsg)
		bot.Send(mess.FromGroupID, mess.SendToType, fmt.Sprintf("成功加入公会成员数： %d", addNum))
	} else {
		bot.Send(mess.FromGroupID, mess.SendToType, fmt.Sprintf("本群所有成员已加入公会 <%s>", clanGroup.GroupName))
	}

}
