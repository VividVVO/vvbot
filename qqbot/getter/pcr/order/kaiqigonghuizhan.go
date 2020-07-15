package order

import (
	"fmt"
	"github.com/gogf/gf/os/gtime"
	"github.com/vivid-vvo/vvbot/app/model/gvg_group"
	"github.com/vivid-vvo/vvbot/app/model/user"
	"github.com/vivid-vvo/vvbot/app/service/gvg"
	"github.com/vivid-vvo/vvbot/qqbot/getter"
	"github.com/vivid-vvo/vvbot/qqbot/model/bot"
	"github.com/vivid-vvo/vvbot/qqbot/model/pcr"
	"strings"
)

type Kaiqigonghuizhan struct {
}

func (l Kaiqigonghuizhan) CheckOrder(cm string) bool {
	return true
}

func (l Kaiqigonghuizhan) IsNotCheckOrder() bool {
	return false
}

func (l Kaiqigonghuizhan) GetOrders() []string {
	return []string{
		"^开启(新)?(公|工|行)会站",
		"^启动(新)?(公|工|行)会站",
	}
}

func (l Kaiqigonghuizhan) Run(mess getter.MeassageData, ms string, atqq int) {
	if !user.CheckUserAuthorityGroup(int(mess.FromUserID), user.AuthAdmin) {
		bot.Send(mess.FromGroupID, 2, "权限不足")
		return
	}
	qqid := int(mess.FromUserID)
	clanGroup, err := pcr.GetClanGroupAndChack(mess.FromGroupID)
	if err != nil {
		bot.Send(mess.FromGroupID, 2, err.Error())
		return
	}
	cms := strings.Split(ms, " ")
	if len(cms) < 2 {
		bot.Send(mess.FromGroupID, 2, "开启公会战格式应为：”开启公会战 公会战名“ 例：“开启公会战 第一次公会战”")
		return
	}
	entity := new(gvg_group.Entity)
	entity.GvgName = cms[1]
	entity.GameServer = clanGroup.GameServer
	entity.CreateQqid = qqid
	entity.GroupId = clanGroup.GroupId
	entity.GvgId = clanGroup.GvgId
	entity.Time = int(gtime.Now().Unix())
	entity.GvgId = gvg_group.CreateGvgId()
	if err = gvg_group.GvgGroupCreate(entity); err != nil {
		bot.Send(mess.FromGroupID, 2, err.Error())
	}
	gvg.BossHpCount(entity.GvgId, clanGroup.GameServer)
	bot.Send(mess.FromGroupID, 2, fmt.Sprintf("公会战<%s>开启成功", entity.GvgName))

}
