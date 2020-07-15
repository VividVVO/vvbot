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
	/*o := new(Order)
	s.BindObject("/{.struct}-{.method}", o)*/

	// 分组路由注册方式
	s.Group("/", func(group *ghttp.RouterGroup) {
		ctlUser := new(user.Controller)
		ctlclan := new(clan.Controller)
		ctlgvg := new(gvg.Controller)
		group.Middleware(middleware.CORS)
		group.ALL("/user", ctlUser, "SignIn")
		group.ALL("/user/signup", ctlUser, "SignUp")
		group.ALL("/user/signin", ctlUser, "SignIn")
		group.ALL("/user/checkqqid", ctlUser, "CheckQQid")
		group.Group("/", func(group *ghttp.RouterGroup) {
			group.Middleware(middleware.Auth)
			group.ALL("/user/issignedin", ctlUser, "IsSignedIn")
			group.ALL("/user/profile", ctlUser, "Profile")
			group.ALL("/user/signout", ctlUser, "SignOut")
			group.ALL("/user/changepassword", ctlUser, "ChangePassword")
			group.ALL("/user/joinclan", ctlclan, "JoinClan")
			group.ALL("/user/reportcauseharm", ctlgvg, "ReportCauseHarm")

			group.Group("/", func(group *ghttp.RouterGroup) {
				group.Middleware(middleware.AuthGvg)
				group.ALL("/user/clangroupcreate", ctlclan, "ClanGroupCreate")
				group.ALL("/user/joinclanto", ctlclan, "JoinClanTo")
				group.ALL("/user/gvggroupcreate", ctlgvg, "GvgGroupCreate")
			})
		})

	})

}
