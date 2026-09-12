package appcfg

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"

	"agigame/console/backend/internal/consts"
)

// SuperAdminId 超级管理员用户ID（可在 config.yaml 的 system.superAdminId 配置）
func SuperAdminId(ctx context.Context) int {
	return g.Cfg().MustGet(ctx, "system.superAdminId", consts.SuperAdminId).Int()
}

// DefaultPassword 新建/重置用户默认密码（可在 config.yaml 的 system.defaultPassword 配置）
func DefaultPassword(ctx context.Context) string {
	return g.Cfg().MustGet(ctx, "system.defaultPassword", "123456").String()
}
