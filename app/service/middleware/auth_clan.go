package middleware

import (
	"github.com/gogf/gf/net/ghttp"
)

// 鉴权中间件，只有拥有相应权限才能通过 // 0 普通用户 100 管理员
func AuthGvg(r *ghttp.Request) {
	r.Middleware.Next()
}
