package order

import (
	"github.com/vivid-vvo/vvbot/app/model/gvg_group"
	"github.com/vivid-vvo/vvbot/qqbot/getter"
	"github.com/vivid-vvo/vvbot/qqbot/model/bot"
	"github.com/vivid-vvo/vvbot/qqbot/model/pcr"

	"regexp"
)

type Suoding struct {
}

func (l Suoding) CheckOrder(cm string) bool {
	return true
}

func (l Suoding) IsNotCheckOrder() bool {
	return false
}

func (l Suoding) GetOrders() []string {
	return []string{
		"^锁定(BOSS)?$|^锁定(BOSS)?( +)?(:|：)",
	}
}
func (l Suoding) Run(mess getter.MeassageData, cm string, atqq int) {
	qqid := int(mess.FromUserID)
	_, gvgGroup, err := pcr.GetGvgGroupDataAtGroupIdToCheck(mess.FromGroupID)
	if err != nil {
		bot.Send(mess.FromGroupID, 2, err.Error())
		return
	}
	regexp1 := regexp.MustCompile(`[:|：](.+)`)
	message := regexp1.FindString(cm)
	if message != "" {
		message = string([]rune(message)[1:])
	}
	if gvgGroup.BossLockType == 0 {
		err := gvg_group.BossLock(qqid, gvgGroup.GvgId, 1, message)
		if err != nil {
			bot.Send(mess.FromGroupID, 2, err.Error())
			return
		}
		bot.Send(mess.FromGroupID, 2, "已锁定BOSS")
	} else {
		bot.Send(mess.FromGroupID, 2, "BOOS已锁定")
	}
}
