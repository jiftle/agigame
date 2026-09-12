package middleware

import (
	"github.com/gogf/gf/v2/net/ghttp"

	"agigame/console/backend/utility/ctxuser"
	"agigame/console/backend/utility/errcode"
)

// Permission 权限校验中间件工厂
func Permission(perm string) ghttp.HandlerFunc {
	return func(r *ghttp.Request) {
		user := ctxuser.Get(r.Context())
		if user == nil {
			r.SetError(errcode.Unauthorized("未登录或登录已过期"))
			return
		}
		if user.HasPerm(perm) {
			r.Middleware.Next()
			return
		}
		r.SetError(errcode.Forbidden("没有操作权限：" + perm))
	}
}
