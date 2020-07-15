package order

import (
	"fmt"
	"github.com/vivid-vvo/vvbot/app/model/clan_group"
	"github.com/vivid-vvo/vvbot/app/model/clan_member"
	"github.com/vivid-vvo/vvbot/app/model/gvg_group"
	"github.com/vivid-vvo/vvbot/app/model/user"
	"github.com/vivid-vvo/vvbot/qqbot/getter"
	"github.com/vivid-vvo/vvbot/qqbot/model/bot"
	"github.com/vivid-vvo/vvbot/qqbot/model/pcr"

	"github.com/vivid-vvo/vvbot/library/Tools"
	"time"
)

type Jiesuo struct {
}

func (l Jiesuo) CheckOrder(cm string) bool {
	return true
}

func (l Jiesuo) IsNotCheckOrder() bool {
	return false
}

func (l Jiesuo) GetOrders() []string {
	return []string{
		"^解锁$",
		"^解除$",
		"^解除锁定$",
		"^锁定解除$"}
}

// 解锁
func (l Jiesuo) Run(mess getter.MeassageData, cm string, atqq int) {
	qqid := int(mess.FromUserID)
	clanGroupData, err := pcr.GetClanGroupAndChack(mess.FromGroupID)
	if err != nil {
		bot.Send(mess.FromGroupID, 2, err.Error())
		return
	}
	clanData, err := clan_group.GetClanData(clanGroupData.GroupId)
	if err != nil {
		bot.Send(mess.FromGroupID, 2, err.Error())
		return
	}
	gvgGroupData, err := gvg_group.GetGvgGroupData(clanData.GvgId)
	if err != nil {
		bot.Send(mess.FromGroupID, 2, err.Error())
		return
	}
	if gvgGroupData.BossLockType == 0 && gvgGroupData.ChallengeStratQqid == 0 {
		bot.Send(mess.FromGroupID, 2, "BOSS未锁定")
		return
	} else {
		if !user.CheckUserAuthorityGroup(qqid, user.AuthAdmin) && gvgGroupData.BossLockType == 1 {
			bot.Send(mess.FromGroupID, 2, "您无权解锁")
			return
		}
	}
	if gvgGroupData.ChallengeStratQqid > 0 {
		if gvgGroupData.ChallengeStratQqid != qqid {
			// 必须管理员权限才能取消别人出刀, 3分钟后谁都可以解锁
			fmt.Println(int(time.Now().Unix()) - gvgGroupData.ChallengeStratTime)
			if !user.CheckUserAuthorityGroup(qqid, user.AuthAdmin) && int(time.Now().Unix())-gvgGroupData.ChallengeStratTime < 60*3 {
				clanMember, err := clan_member.GetClanMember(gvgGroupData.ChallengeStratQqid, clanData.GroupId)
				if err != nil {
					bot.Send(mess.FromGroupID, 2, err.Error())
					return
				}
				var gameName string
				if clanMember != nil {
					gameName = clanMember.GameName
				}

				bot.Send(mess.FromGroupID, 2, fmt.Sprintf("%s正在挑战boss，但是，您无权解锁", gameName))
				return
			}
		}
		if err = gvg_group.CancelChallenge(gvgGroupData.GvgId); err != nil {
			bot.Send(mess.FromGroupID, 2, err.Error())
			return
		}
		// model.Send(mess.FromGroupID, 2, fmt.Sprintf("成功解锁出刀"))
		if gvgGroupData.BossLockType == 0 {
			msg := fmt.Sprintf("boss挑战已可申请\n%s", pcr.GetBossStateStr(mess.FromGroupID))
			bot.Send(mess.FromGroupID, 2, msg)
			return
		}
	}

	if err = gvg_group.UnLockBoss(gvgGroupData.GvgId); err != nil {
		bot.Send(mess.FromGroupID, 2, err.Error())
		return
	}
	msg := fmt.Sprintf("boss挑战已可申请\n现在%d周目，%d号boss\n生命值%s", gvgGroupData.BossCycle,
		gvgGroupData.BossNum, Tools.NumberFormat(gvgGroupData.BossHp))
	bot.Send(mess.FromGroupID, 2, msg)
	return
}
