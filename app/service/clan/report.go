package clan

import (
	"errors"
	"fmt"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/util/gconv"
	"github.com/gogf/gf/util/grand"
	"github.com/gogf/gf/util/gvalid"
	"github.com/vivid-vvo/vvbot/app/model/clan_group"
	"github.com/vivid-vvo/vvbot/app/model/clan_member"
	"github.com/vivid-vvo/vvbot/app/service/user"
)

// 创建公会组参数
type ClanGroupCreateInput struct {
	GroupName  string `v:"required#公会名不能为空"` // 公会名
	GameServer string `v:"in:CN,TW,JP,KR#游戏服务器应当在国,台,日,韩选择其一"`
}

// 公会组创建
func ClanGroupCreate(r *ghttp.Request, cgc *ClanGroupCreateInput) error {
	// 输入参数检查
	if e := gvalid.CheckStruct(cgc, nil); e != nil {
		return errors.New(e.FirstString())
	}
	var entity *clan_group.Entity

	qqid := user.GetLoginData2(r).Qqid

	if err := gconv.Struct(cgc, &entity); err != nil {
		return err
	}
	clanGroupData, err := clan_group.GetClanGroupAtName(cgc.GroupName)
	if err != nil {
		return err
	}
	if clanGroupData != nil {
		return errors.New(fmt.Sprintf("公会名已存在"))
	}
	entity.CreatorQqid = qqid
	entity.CreateTime = gtime.Now().Unix()
	entity.GroupId = grand.N(100000000, 1000000000)
	if err := clan_group.ClanGroupCreate(entity); err != nil {
		return err
	}
	return nil
}

// 加入公会 当前用户
func JoinClan(r *ghttp.Request, qqid int64, groupName string) error {
	data, err := clan_group.GetClanGroupAtName(groupName)
	if err != nil {
		return errors.New(fmt.Sprintf("内部错误"))
	}
	if data == nil {
		return errors.New(fmt.Sprintf("指定公会不存在！"))
	}
	if qqid == 0 {
		qqid = user.GetLoginData2(r).Qqid
	}
	userInfo, err := user.GetProfile(r, qqid)
	if err != nil {
		return err
	}
	entity := new(clan_member.Entity)
	entity.GroupId = data.GroupId
	entity.Qqid = qqid
	entity.GameName = userInfo.Nickname
	entity.JoinTime = gtime.Now().Unix()
	err = clan_member.JoinClan(entity)
	if err != nil {
		return errors.New(fmt.Sprintf("内部错误"))
	}
	if err := user.ChangeClanGroupId(r, qqid, data.GroupId); err != nil {
		return errors.New(fmt.Sprintf("内部错误"))
	}
	return nil
}

// 获取公会成员信息
func GetUserClanMemberData(r *ghttp.Request, qqid int64) (*clan_member.Entity, error) {
	if qqid == 0 {
		qqid = user.GetLoginData2(r).Qqid
	}
	//return clan_member.GetClanMember(qqid)
	return nil, nil
}
