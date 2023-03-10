package plugin

import (
	"fmt"
	"github.com/vivid-vvo/vvbot/app/model/gvg_group"
	"github.com/vivid-vvo/vvbot/app/model/gvg_member_extra"
	"github.com/vivid-vvo/vvbot/library/tools"
	"github.com/vivid-vvo/vvbot/qqbot/model/bot"
	"github.com/vivid-vvo/vvbot/qqbot/model/pcr"
	"github.com/vivid-vvo/vvbot/qqbot/plugins"
	"regexp"
	"strings"
)

type Guashu struct {
}

func (l Guashu) CheckOrder(cm string) bool {
	return true
}

func (l Guashu) IsNotCheckOrder() bool {
	return false
}

func (l Guashu) GetOrders() []string {
	return []string{
		"^(ζ(ζ |π²)|δΈ(ζ |π²))( +)?(\\?|οΌ)?$",
		"^(ζ(ζ |π²)|δΈ(ζ |π²))( +)?(\\:|οΌ)",

		"^(ζ(ζ |π²)|δΈ(ζ |π²))(εζΆ|ζ€ι)( +)?(\\?|οΌ)?$",
		"^(εζΆ|ζ€ι)(ζ(ζ |π²)|δΈ(ζ |π²))( +)?(\\?|οΌ)?$",
		"^δΈ(ζ |π²)( +)?$",
	}
}

func (l Guashu) Run(mess plugins.MeassageData, cm string, atqq int64) {
	var agentQqid, qqid int64
	if atqq != 0 {
		agentQqid = mess.FromUserID
		qqid = atqq
	} else {
		qqid = mess.FromUserID
	}
	cm = strings.ToLower(cm)
	isCancel := tools.Compare(cm, "^(εζΆ|ζ€ι|δΈζ )|(εζΆ|ζ€ι)$")
	isQuery := tools.Compare(cm, "[?|οΌ]")

	reg := regexp.MustCompile(`[:|οΌ](.+)`)
	message := reg.FindString(cm)
	if message != "" {
		message = string([]rune(message)[1:])
	}
	_, clanMember, err := pcr.GetClanGroupAndUserGroupToCheck(mess.FromGroupID, qqid)
	if err != nil {
		bot.Send(mess.FromGroupID, mess.SendToType, err.Error())
		return
	}
	_, gvgGroup, err := pcr.GetGvgGroupDataAtGroupIdToCheck(mess.FromGroupID)
	if err != nil {
		bot.Send(mess.FromGroupID, mess.SendToType, err.Error())
		return
	}
	treeData, err := gvg_member_extra.GetIsUpTree(qqid, gvgGroup.GvgId)
	if err != nil {
		bot.Send(mess.FromGroupID, mess.SendToType, err.Error())
		return
	}
	if treeData != nil && !isCancel {
		bot.Send(mess.FromGroupID, mess.SendToType, fmt.Sprintf("%s ζ¨ε·²ε¨ζ δΈ [%s]", clanMember.GameName, treeData.Message))
		return
	}
	if isQuery || isCancel && treeData == nil {
		bot.Send(mess.FromGroupID, mess.SendToType, fmt.Sprintf("%s ζ¨ζͺε¨ζ δΈ", clanMember.GameName))
		return
	}
	if isCancel {
		// ε ζ 
		err = gvg_member_extra.ReportDownTree(qqid, gvgGroup.GvgId, 0)
		if err != nil {
			bot.Send(mess.FromGroupID, mess.SendToType, err.Error())
			return
		}
		bot.Send(mess.FromGroupID, mess.SendToType, fmt.Sprintf("%s ε·²δΈζ ", clanMember.GameName))
	} else {
		// ζ₯εδΈζ 
		err = gvg_member_extra.ReportUpTree(qqid, agentQqid, gvgGroup.GvgId, message)
		if err != nil {
			bot.Send(mess.FromGroupID, mess.SendToType, err.Error())
			return
		}
		bot.Send(mess.FromGroupID, mess.SendToType, fmt.Sprintf("%s ε·²θ?°ε½ζζ ", clanMember.GameName))
		if qqid == gvgGroup.ChallengeStratQqid {
			// εΊειε?θ§£ι
			if err = gvg_group.CancelChallenge(gvgGroup.GvgId); err != nil {
				bot.Send(mess.FromGroupID, mess.SendToType, err.Error())
				return
			}
			msg := fmt.Sprintf("bossζζε·²ε―η³θ―·\nη°ε¨%dε¨η?οΌ%dε·boss\nηε½εΌ%s", gvgGroup.BossCycle, gvgGroup.BossNum, tools.NumberFormat(gvgGroup.BossHp))
			bot.Send(mess.FromGroupID, mess.SendToType, msg)
		}
	}

}
