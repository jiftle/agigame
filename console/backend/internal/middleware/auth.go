package middleware

import (
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"

	"agigame/console/backend/internal/service"
	"agigame/console/backend/utility/ctxuser"
	"agigame/console/backend/utility/errcode"
	jwtutil "agigame/console/backend/utility/jwtutil"
)

// Auth 登录鉴权中间件
func Auth(r *ghttp.Request) {
	ctx := r.Context()
	token := ExtractToken(r)
	if token == "" {
		r.SetError(errcode.Unauthorized("未登录或登录已过期"))
		return
	}

	claims, err := jwtutil.Parse(ctx, token)
	if err != nil {
		r.SetError(errcode.TokenExpired("登录已过期，请重新登录"))
		return
	}
	if claims.Type != jwtutil.TypeAccess {
		r.SetError(errcode.Unauthorized("无效的令牌类型"))
		return
	}

	user, err := service.Auth().GetContextUser(ctx, claims.UserId)
	if err != nil {
		r.SetError(err)
		return
	}
	ctxuser.Set(ctx, user)
	r.Middleware.Next()
}

// ExtractToken 从请求头或查询参数中提取 token
func ExtractToken(r *ghttp.Request) string {
	header := g.Cfg().MustGet(r.Context(), "jwt.header", "Authorization").String()
	prefix := g.Cfg().MustGet(r.Context(), "jwt.prefix", "Bearer ").String()

	token := r.Header.Get(header)
	if token == "" {
		token = r.Get("token").String()
	}
	token = strings.TrimSpace(token)
	if prefix != "" {
		token = strings.TrimPrefix(token, strings.TrimSpace(prefix))
	}
	return strings.TrimSpace(token)
}
