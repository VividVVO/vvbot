package plugin

import (
	"fmt"
	"github.com/vivid-vvo/vvbot/app/model/clan_member"
	"github.com/vivid-vvo/vvbot/app/model/gvg_group"
	"github.com/vivid-vvo/vvbot/app/service/check"
	"github.com/vivid-vvo/vvbot/library/tools"
	"github.com/vivid-vvo/vvbot/qqbot/model/bot"
	"github.com/vivid-vvo/vvbot/qqbot/model/pcr"
	"github.com/vivid-vvo/vvbot/qqbot/plugins"
	"time"
)

type Quxiaochudao struct {
}

func (l Quxiaochudao) CheckOrder(cm string) bool {
	return true
}

func (l Quxiaochudao) IsNotCheckOrder() bool {
	return false
}

func (l Quxiaochudao) GetOrders() []string {
	return []string{"^取消出刀$",
		"^出刀取消$",
		"^撤回出刀$",
		"^撤回出刀$",
		"^解锁出刀$",
		"^出刀解锁$",
		"^结束出刀$",
		"^出刀结束$",
		"^取消申请出刀$",
		"^踢出队列$",
		"^撤退$",
		"^逃跑$",
	}

}

// 取消出刀
func (l Quxiaochudao) Run(mess plugins.MeassageData, cm string, atqq int64) {
	qqid := mess.FromUserID
	clanGroup, gvgGroupData, err := pcr.GetGvgGroupDataAtGroupIdToCheck(mess.FromGroupID)
	if err != nil {
		bot.Send(mess.FromGroupID, mess.SendToType, err.Error())
		return
	}
	if gvgGroupData.ChallengeStratQqid > 0 {
		if gvgGroupData.ChallengeStratQqid != qqid {
			isAdmin := check.CheckAuthorityGroup(qqid, check.AuthGvgAdmin, clanGroup.GroupId)

			// 必须管理员权限才能取消别人出刀, 3分钟后谁都可以解锁
			if !isAdmin && time.Now().Unix()-gvgGroupData.ChallengeStratTime > 60*3 {
				clanMember, err := clan_member.GetClanMember(gvgGroupData.BossLockQqid, clanGroup.GroupId)
				if err != nil {
					bot.Send(mess.FromGroupID, mess.SendToType, err.Error())
					return
				}
				var gameName string
				if clanMember != nil {
					gameName = clanMember.GameName
				}
				bot.Send(mess.FromGroupID, mess.SendToType, fmt.Sprintf("%s正在挑战boss，但是，您 无 权 取 消", gameName))
				return
			}
		}
		if err = gvg_group.CancelChallenge(gvgGroupData.GvgId); err != nil {
			bot.Send(mess.FromGroupID, mess.SendToType, err.Error())
			return
		}
		bot.Send(mess.FromGroupID, mess.SendToType, fmt.Sprintf("成功取消出刀"))
		if gvgGroupData.BossLockType == 0 {
			msg := fmt.Sprintf("boss挑战已可申请\n现在%d周目，%d号boss\n生命值%s", gvgGroupData.BossCycle, gvgGroupData.BossNum, tools.NumberFormat(gvgGroupData.BossHp))
			bot.Send(mess.FromGroupID, mess.SendToType, msg)
			return
		}
	}
	bot.Send(mess.FromGroupID, mess.SendToType, "当前没有人正在挑战BOSS")
	return
}
