package clan

import (
	"errors"
	"fmt"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/util/gconv"
	"github.com/vivid-vvo/vvbot/app/model/clan_group"
	"github.com/vivid-vvo/vvbot/app/model/clan_member"
	user2 "github.com/vivid-vvo/vvbot/app/model/user"
	"github.com/vivid-vvo/vvbot/app/service/check"
	"github.com/vivid-vvo/vvbot/app/service/user"
)

type UserClanListInput struct {
	GroupId      int    `orm:"group_id"      json:"group_id"`      //
	GroupName    string `orm:"group_name"    json:"group_name"`    //
	GameServer   string `orm:"game_server"   json:"game_server"`   //
	Notification string `orm:"notification"  json:"notification"`  //
	GvgId        int    `orm:"gvg_id"        json:"gvg_id"`        //
	MemberNum    int    `orm:"member_num"       json:"member_num"` //
	BindQqGroup  int64  `orm:"bind_qq_group" json:"bind_qq_group"` //
}

// GetClan 获取用户公会列表
func GetUserClanList(r *ghttp.Request, qqid int64) ([]*UserClanListInput, error) {
	if qqid == 0 {
		qqid = user.GetLoginData2(r).Qqid
	}

	clanList, err := clan_member.GetUserClanList(qqid)
	if err != nil {
		return nil, err
	}
	var userClanList []*UserClanListInput
	var groupIDList []int
	if clanList == nil {
		return nil, nil
	}
	for _, clan := range clanList {
		groupIDList = append(groupIDList, clan.GroupId)
	}
	clanGroupList, err := clan_group.GetClanGroupList(groupIDList)
	if err != nil {
		return nil, err
	}
	for _, clan := range clanGroupList {
		var userClan *UserClanListInput
		if err := gconv.Struct(clan.Entity, &userClan); err != nil {
			return nil, err
		}
		userClan.MemberNum = clan.MemberNum
		userClanList = append(userClanList, userClan)
	}
	return userClanList, nil
}

// GetClan 获取所有公会
func GetAllClan(r *ghttp.Request) ([]*clan_group.ClanGroupEntity, error) {

	clanGroupList, err := clan_group.GetAllClanGroup()
	if err != nil {
		return nil, err
	}
	return clanGroupList, nil
}

// 创建公会组参数
type ChangeClanInfoInput struct {
	GroupId     int    `v:"min:1#公会ID不能为空"`   //
	GroupName   string `v:"required#公会名不能为空"` // 公会名
	GameServer  string `v:"in:CN,TW,JP,KR,#参数错误"`
	BindQqGroup int64
	Apikey      string
}

// ChangeClanInfo 修改公会信息
func ChangeClanInfo(r *ghttp.Request, data ChangeClanInfoInput) error {
	qqid := user.GetLoginData2(r).Qqid
	if !check.CheckAuthorityGroup(qqid, check.AuthClanAdmin, data.GroupId) {
		return errors.New("权限不足")
	}
	return clan_group.ChangeClanInfo(data.GroupId, data.GroupName, data.BindQqGroup, data.GameServer, data.Apikey)
}

// DelClanGroup 删除公会
func DelClanGroup(r *ghttp.Request, groupId int) error {
	qqid := user.GetLoginData2(r).Qqid
	if !check.CheckAuthorityGroup(qqid, check.AuthAdmin, 0) {
		return errors.New("权限不足")
	}

	tx, err := g.DB().Begin()
	if err != nil {
		return errors.New("内部错误")
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	_, err = tx.Table("clan_group").Delete("group_id", groupId)
	if err != nil {
		return errors.New("内部错误")
	}
	_, err = tx.Table("clan_member").Delete("group_id", groupId)
	if err != nil {
		return errors.New("内部错误")
	}
	_, err = tx.Table("gvg_challenge").Delete("clan_group_id", groupId)
	if err != nil {
		return errors.New("内部错误")
	}
	err = tx.Commit()
	if err != nil {
		return errors.New("内部错误")
	}
	return nil
}

func GetClanGroupMembers(r *ghttp.Request, ClanGroupID int) ([]*clan_member.Entity, error) {
	members, err := clan_member.GetAllClanMember(ClanGroupID)
	if err != nil {
		return nil, err
	}
	return members, nil
}

func GetClanDataAtQqid(r *ghttp.Request, qqid int64) (*clan_group.Entity, error) {
	if qqid == 0 {
		qqid = user.GetLoginData2(r).Qqid
	}
	clanGroupData, err := GetUserClanMemberData(r, qqid)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("内部错误"))
	}
	if clanGroupData == nil {
		return nil, errors.New(fmt.Sprintf("当前成员未加入公会"))
	}
	return clan_group.GetClanData(clanGroupData.GroupId)
}

// ChangeMembersData 修改成员信息 公会管理员可以赋予会战管理员权限，管理员及以上才能赋予公会管理员权限及修改自身权限
func ChangeMembersData(r *ghttp.Request, qqid int64, clanGroupID int, gameName string, role int) error {
	thisqqid := user.GetLoginData2(r).Qqid
	clanMember, err := clan_member.GetClanMember(qqid, clanGroupID)
	if err != nil {
		return err
	}
	// 是否修改了用户权限
	if clanMember.Role != role {
		if thisqqid == qqid {
			// 必须拥有管理员权限才能修改自己的公会权限
			if !check.CheckAuthorityGroup(thisqqid, check.AuthAdmin, clanGroupID) {
				return errors.New("权限不足")
			}
		}
		if role == check.AuthUser || role == check.AuthGvgAdmin {
			// 拥有公会管理员权限才能赋予公会成员权限
			if !check.CheckAuthorityGroup(thisqqid, check.AuthClanAdmin, clanGroupID) {
				return errors.New("权限不足")
			}
		} else if role == check.AuthClanAdmin {
			// 拥有管理员权限才能赋予公会管理员权限
			if !check.CheckAuthorityGroup(thisqqid, check.AuthAdmin, 0) {
				return errors.New("权限不足")
			}
		} else {
			return errors.New("参数错误")
		}
	} else if thisqqid != qqid {
		// 必须拥有公会管理员权限才能修改成员信息
		if !check.CheckAuthorityGroup(thisqqid, check.AuthClanAdmin, 0) {
			return errors.New("权限不足")
		}
	}
	return clan_member.ChangeMembersData(qqid, clanGroupID, gameName, role)
}

// 成员退出公会
func MemberExitGroup(r *ghttp.Request, qqid int64, clanGroupID int) error {
	thisqqid := user.GetLoginData2(r).Qqid
	if qqid == 0 {
		qqid = thisqqid
	}
	clanGroup, _, err := check.GetClanGroupAndUserGroupToCheck(clanGroupID, qqid)
	if err != nil {
		return err
	}
	if qqid != thisqqid {
		// 必须拥有公会管理员权限才能将他人踢出公会
		if !check.CheckAuthorityGroup(thisqqid, check.AuthClanAdmin, clanGroup.GroupId) {
			return errors.New("权限不足")
		}
	}
	err = clan_member.MemberExitGroupAtQqid(qqid, clanGroupID)
	if err != nil {
		return err
	}
	if err := user2.ChangeClanGroupId(qqid, 0); err != nil {
		return err
	}
	return nil
}
