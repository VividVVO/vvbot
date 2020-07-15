package middleware

import (
	"github.com/vivid-vvo/vvbot/app/service/user"
	"github.com/vivid-vvo/vvbot/library/response"

	"github.com/gogf/gf/net/ghttp"
)

// 鉴权中间件，只有拥有相应权限才能通过 // 0 普通用户 100 管理员
func AuthGvg(r *ghttp.Request) {
	authorityGroup, err := user.GetUserAuthorityGroup(r, 0)
	if err != nil {
		response.JsonExit(r, 1, "", err.Error())
	}
	if authorityGroup > 0 {
		r.Middleware.Next()
	} else {
		response.JsonExit(r, -200, "", "权限不足")
	}
}
