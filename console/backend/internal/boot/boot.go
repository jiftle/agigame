package boot

import (
	"context"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gfile"
)

// InitDatabase 首次启动时自动初始化 SQLite 数据库与种子数据。
func InitDatabase(ctx context.Context) error {
	_ = gfile.Mkdir("data")

	db := g.DB()
	value, err := db.GetValue(ctx, "SELECT count(*) FROM sqlite_master WHERE type='table' AND name='sys_user'")
	if err != nil {
		return err
	}
	if value.Int() == 0 {
		if err = initSchema(ctx); err != nil {
			return err
		}
	}

	return upgradeSoftDeleteIndexes(ctx)
}

func initSchema(ctx context.Context) error {
	g.Log().Info(ctx, "初始化数据库: manifest/sql/init.sql")
	content := gfile.GetContents("manifest/sql/init.sql")
	if content == "" {
		g.Log().Warning(ctx, "未找到 manifest/sql/init.sql, 跳过初始化")
		return nil
	}
	db := g.DB()
	for _, stmt := range splitStatements(content) {
		if _, err := db.Exec(ctx, stmt); err != nil {
			g.Log().Errorf(ctx, "执行初始化SQL失败: %v\nSQL: %s", err, stmt)
			return err
		}
	}
	g.Log().Info(ctx, "数据库初始化完成")
	return nil
}

// upgradeSoftDeleteIndexes 将唯一索引升级为排除软删除记录的部分索引。
// 旧库的普通唯一索引会导致"删除后重建同名记录"命中外键冲突，这里做幂等迁移。
func upgradeSoftDeleteIndexes(ctx context.Context) error {
	indexes := []struct {
		name string
		sql  string
	}{
		{"uk_sys_user_username", "CREATE UNIQUE INDEX IF NOT EXISTS uk_sys_user_username ON sys_user (username) WHERE deleted_at IS NULL"},
		{"uk_sys_role_code", "CREATE UNIQUE INDEX IF NOT EXISTS uk_sys_role_code ON sys_role (code) WHERE deleted_at IS NULL"},
		{"uk_sys_dict_type_type", "CREATE UNIQUE INDEX IF NOT EXISTS uk_sys_dict_type_type ON sys_dict_type (type) WHERE deleted_at IS NULL"},
		{"uk_sys_config_key", "CREATE UNIQUE INDEX IF NOT EXISTS uk_sys_config_key ON sys_config (config_key) WHERE deleted_at IS NULL"},
	}
	for _, idx := range indexes {
		if err := ensurePartialUniqueIndex(ctx, idx.name, idx.sql); err != nil {
			return err
		}
	}
	return nil
}

func ensurePartialUniqueIndex(ctx context.Context, name, createSQL string) error {
	db := g.DB()
	existing, err := db.GetValue(ctx, "SELECT sql FROM sqlite_master WHERE type='index' AND name=?", name)
	if err != nil {
		return err
	}
	existingSQL := existing.String()
	if existingSQL != "" && strings.Contains(strings.ToUpper(existingSQL), "WHERE") {
		return nil
	}
	if existingSQL != "" {
		if _, err = db.Exec(ctx, "DROP INDEX IF EXISTS "+name); err != nil {
			return err
		}
	}
	_, err = db.Exec(ctx, createSQL)
	return err
}

func splitStatements(content string) []string {
	var sb strings.Builder
	for _, line := range strings.Split(content, "\n") {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "--") {
			continue
		}
		sb.WriteString(line)
		sb.WriteString("\n")
	}
	var stmts []string
	for _, part := range strings.Split(sb.String(), ";") {
		if s := strings.TrimSpace(part); s != "" {
			stmts = append(stmts, s)
		}
	}
	return stmts
}
