package gvg

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/gogf/gf/net/ghttp"
	"github.com/vivid-vvo/vvbot/app/model/clan_member"
	"github.com/vivid-vvo/vvbot/app/model/gvg_challenge"
	"github.com/vivid-vvo/vvbot/app/service/check"
	excel2 "github.com/vivid-vvo/vvbot/library/excel"
	"github.com/vivid-vvo/vvbot/library/tools"
	"strconv"
)

type GetAllChallengeToExcelInput struct {
	ExcelBuffer *bytes.Buffer
	TimeStr     string `json:"timeStr"`
}

type challengeExcelType struct {
	qqid     int64
	gameName string
	c1       string
	c1s      string
	c2       string
	c2s      string
	c3       string
	c3s      string
	total    string
	finished int
}

func DownAllChallengeToExcel(r *ghttp.Request, clanGroupId int, TimeType string) (*GetAllChallengeToExcelInput, error) {
	_, err := check.GetClanGroupAndChack(clanGroupId)
	if err != nil {
		return nil, err
	}
	challenges, err := GetAllChallenge(r, clanGroupId, TimeType, "", "")
	if err != nil {
		return nil, err
	}
	members, err := clan_member.GetAllClanMember(clanGroupId)
	if err != nil {
		return nil, err
	}
	challengeTable := bossChallengeTable(challenges.List, members)
	var finished int
	var challengeExcel []*challengeExcelType
	for _, m := range challengeTable {
		challengeExcel = append(challengeExcel, &challengeExcelType{
			qqid:     m.member.Qqid,
			gameName: m.member.GameName,
			c1:       challengeExcelStr(m.detail[0]), c1s: challengeExcelStr(m.detail[1]),
			c2: challengeExcelStr(m.detail[2]), c2s: challengeExcelStr(m.detail[3]),
			c3: challengeExcelStr(m.detail[4]), c3s: challengeExcelStr(m.detail[5]),
			total:    tools.NumberFormat(challengeTotalDamage(m.detail)),
			finished: m.finished,
		})
		finished += m.finished
	}
	file := excelize.NewFile()
	sheet := excel2.NewSheet(file, "战斗记录", 11, 22)

	sheet.SetAllColsWidth(7, 13, 25, 7, 16, 16, 16, 16, 16, 16, 12)

	sheet.WriteRow(challenges.TimeStr + " 战斗记录")
	sheet.WriteRow("序号", "QQ号", "游戏名", "出刀数", "第一刀", "第二刀", "第三刀", "余刀一", "余刀二", "余刀三", "总伤害")
	for i, e := range challengeExcel {
		sheet.WriteRow(strconv.Itoa(i+1), strconv.FormatInt(e.qqid, 10), e.gameName, strconv.Itoa(e.finished), e.c1, e.c2, e.c3, e.c1s, e.c2s, e.c3s, e.total)
	}
	sheet.WriteRow(" ", " ", " ", strconv.Itoa(finished), " ", " ", " ", " ", " ", " ", tools.NumberFormat(challengeAllTotalDamage(challengeTable)))
	buf, err := file.WriteToBuffer()
	if err != nil {
		return nil, errors.New("内部错误")
	}
	getAllChallengeToExcelInput := &GetAllChallengeToExcelInput{
		ExcelBuffer: buf,
		TimeStr:     challenges.TimeStr,
	}
	return getAllChallengeToExcelInput, nil
}

func challengeAllTotalDamage(bossDamageList []*bossChallengeTableInput) int {
	var totalDamage = 0
	for _, m := range bossDamageList {
		totalDamage += challengeTotalDamage(m.detail)
	}
	return totalDamage

}

func challengeTotalDamage(detailList []*gvg_challenge.Entity) int {
	if detailList == nil {
		return 0
	}
	var totalDamage = 0
	for _, data := range detailList {
		if data != nil {
			totalDamage += data.ChallengeDamage
		}
	}
	return totalDamage

}

func challengeExcelStr(detail *gvg_challenge.Entity) string {
	if detail == nil {
		return " "
	}
	return fmt.Sprintf("(%d-%d)%s", detail.BossCycle, detail.BossNum, tools.NumberFormat(detail.ChallengeDamage))
}

type bossChallengeTableInput struct {
	// detail 0为第一刀，1为第一刀的尾刀； 2为第二刀，3为第二刀的尾刀、、、
	detail   []*gvg_challenge.Entity
	member   *clan_member.Entity
	finished int
}

type userIndexType struct {
	index int
	num   int
}

func bossChallengeTable(gvgChallenges []*gvg_challenge.Entity, members []*clan_member.Entity) []*bossChallengeTableInput {
	var challengeTable []*bossChallengeTableInput
	var userNum int
	userIndex := make(map[int64]*userIndexType)
	for _, m := range members {
		userIndex[m.Qqid] = &userIndexType{index: userNum, num: 0}
		challengeTable = append(challengeTable, &bossChallengeTableInput{})
		challengeTable[userNum].detail = make([]*gvg_challenge.Entity, 6)
		challengeTable[userNum].member = m
		userNum++
	}
	for _, damage := range gvgChallenges {
		if _, ok := userIndex[damage.Qqid]; !ok {
			continue
		}
		var index = userIndex[damage.Qqid].index
		if damage.IsSurplus == 0 && userIndex[damage.Qqid].num > 0 {
			if userIndex[damage.Qqid].num > 0 && challengeTable[index].detail[userIndex[damage.Qqid].num-1].IsSurplus == 0 {
				userIndex[damage.Qqid].num++
			}
		}
		if damage.IsContinue == 0 || damage.IsSurplus > 0 {
			challengeTable[index].finished++
		}
		challengeTable[index].detail[userIndex[damage.Qqid].num] = damage
		userIndex[damage.Qqid].num++
	}
	return challengeTable

}
