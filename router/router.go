package router

import (
	"github.com/vivid-vvo/vvbot/app/api/clan"
	"github.com/vivid-vvo/vvbot/app/api/gvg"
	"github.com/vivid-vvo/vvbot/app/api/user"
	"github.com/vivid-vvo/vvbot/app/service/middleware"

	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
)

func init() {
	s := g.Server()
	// 某些浏览器直接请求favicon.ico文件，特别是产生404时
	s.SetRewrite("/favicon.ico", "/resource/image/favicon.ico")
	s.SetNameToUriType(ghttp.URI_TYPE_ALLLOWER)
	// 分组路由注册方式
	s.Group("/api", func(group *ghttp.RouterGroup) {
		group.Middleware(middleware.CORS)
		ctlUser := new(user.Controller)
		group.ALL("/user/login", ctlUser, "Login")
		group.Group("/user", func(group *ghttp.RouterGroup) {
			group.Middleware(middleware.Auth)
			group.ALL("/deluser", ctlUser, "DelUser")
			group.ALL("/changepassword", ctlUser, "ChangePassword")
			group.ALL("/profile", ctlUser, "Profile")
			group.ALL("/signup", ctlUser, "SignUp")
			group.ALL("/checkqqid", ctlUser, "CheckQQid")
			group.ALL("/getuserlist", ctlUser, "GetUserList")
			group.ALL("/changeuserdata", ctlUser, "ChangeUserData")
		})
	})

	// 分组路由注册方式
	s.Group("/api", func(group *ghttp.RouterGroup) {
		group.Middleware(middleware.CORS, middleware.Auth)
		group.Group("/pcr", func(group *ghttp.RouterGroup) {
			ctlclan := new(clan.Controller)
			ctlgvg := new(gvg.Controller)
			group.ALL("/getuserclanlist", ctlclan, "GetUserClanList")
			group.ALL("/getclangroupmembers", ctlclan, "GetClanGroupMembers")
			group.ALL("/changemembersdata", ctlclan, "ChangeMembersData")
			group.ALL("/memberexitgroup", ctlclan, "MemberExitGroup")
			group.ALL("/getallclan", ctlclan, "GetAllClan")
			group.ALL("/changeclaninfo", ctlclan, "ChangeClanInfo")
			group.ALL("/delclangroup", ctlclan, "DelClanGroup")

			group.ALL("/getclangvg", ctlgvg, "GetClanGvg")
			group.ALL("/getchallengeatqq", ctlgvg, "GetChallengeAtQQ")
			group.ALL("/getallchallenge", ctlgvg, "GetAllChallenge")
			group.ALL("/remindchallenge", ctlgvg, "RemindChallenge")
			group.ALL("/changeuserchallenge", ctlgvg, "ChangeUserChallenge")
			group.ALL("/deluserchallenge", ctlgvg, "DelUserChallenge")
			group.ALL("/getallslstate", ctlgvg, "GetAllSlState")
			group.ALL("/downallchallengetoexcel", ctlgvg, "DownAllChallengeToExcel")

		})
	})

}
