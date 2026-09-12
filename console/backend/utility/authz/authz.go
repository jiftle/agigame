package authz

import (
	"context"
	"strings"

	"agigame/console/backend/utility/ctxuser"
	"agigame/console/backend/utility/errcode"
)

// Check 校验当前用户是否拥有指定权限
func Check(ctx context.Context, perm string) error {
	return CheckAny(ctx, perm)
}

// CheckAny 校验当前用户是否拥有其中任意一个权限
func CheckAny(ctx context.Context, perms ...string) error {
	user := ctxuser.Get(ctx)
	if user == nil {
		return errcode.Unauthorized("未登录或登录已过期")
	}
	for _, perm := range perms {
		if user.HasPerm(perm) {
			return nil
		}
	}
	return errcode.Forbidden("没有操作权限：" + strings.Join(perms, " 或 "))
}
