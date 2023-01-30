package check

import (
	"github.com/gogf/gf/frame/g"
	"github.com/vivid-vvo/vvbot/app/model/clan_member"
	"github.com/vivid-vvo/vvbot/app/model/user"
)

const (
	// 普通用户
	AuthUser = 0

	// 会战管理员
	AuthGvgAdmin = 9

	// 公会管理员
	AuthClanAdmin = 10

	// 管理员
	AuthAdmin = 100

	// 超级管理员
	AuthSuperAdmin = 200
)

// checkUserAuthorityGroup 检测用户组权限 普通用户 < 管理员 < 200超级管理员
func checkUserAuthorityGroup(qqid int64, auth int) bool {
	auth1, err := user.GetUserAuthorityGroup(qqid)
	if err != nil {
		return false
	}
	return auth1 >= auth
}

// CheckAuthorityGroup 检测用户组与公会权限 普通用户 < 会战管理员 < 公会管理员 < 管理员 < 200超级管理员
func CheckAuthorityGroup(qqid int64, auth int, groupID int) bool {
	if CheckIsMaster(qqid) {
		return true
	}
	if groupID > 0 {
		if clan_member.CheckClanMemberAuthorityGroup(qqid, auth, groupID) {
			return true
		}
	}
	if checkUserAuthorityGroup(qqid, auth) {
		return true
	}
	return false
}

var Masters = g.Cfg().GetInterfaces("auth.masters")

// CheckIsMaster 检测是否是机器人的主人
func CheckIsMaster(qqid int64) bool {
	for _, qqid2 := range Masters {
		if qqid == int64(qqid2.(float64)) {
			return true
		}
	}
	return false
}
