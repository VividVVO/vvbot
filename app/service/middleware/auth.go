package middleware

import (
	"github.com/vivid-vvo/vvbot/app/service/user"
	"github.com/vivid-vvo/vvbot/library/response"

	"github.com/gogf/gf/net/ghttp"
)

// 鉴权中间件，只有登录成功之后才能通过
func Auth(r *ghttp.Request) {
	if user.IsSignedIn(r) {
		r.Middleware.Next()
	} else {
		response.JsonExit(r, 401, "", "账号未登录")
	}
}
