package ctxuser

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"

	"agigame/console/backend/internal/consts"
	"agigame/console/backend/internal/model"
)

// Set 将当前登录用户写入请求上下文
func Set(ctx context.Context, user *model.ContextUser) {
	if r := g.RequestFromCtx(ctx); r != nil {
		r.SetCtxVar(consts.CtxUserKey, user)
	}
}

// Get 获取当前登录用户
func Get(ctx context.Context) *model.ContextUser {
	r := g.RequestFromCtx(ctx)
	if r == nil {
		return nil
	}
	v := r.GetCtxVar(consts.CtxUserKey)
	if v.IsNil() {
		return nil
	}
	if u, ok := v.Val().(*model.ContextUser); ok {
		return u
	}
	return nil
}
