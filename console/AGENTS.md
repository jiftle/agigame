# AGENTS.md

AdminBase：GoFrame v2 后端 + Vite/React 前端的后台管理基座，前后端通过 `/api/v1` JSON 通讯。文档与提交信息用中文。

## 仓库边界与入口

- `backend/`：Go 模块，**module 名是 `agigame/console/backend`**。新增 Go 文件的 import 必须写 `agigame/console/backend/...`。
- `frontend/`：Vite + React 19 + antd 6 + pro-components + react-router-dom + TanStack Query + Zustand；**已彻底移除 Umi**。
- 后端入口 `backend/main.go` → `internal/cmd`（先建库、再注册路由）。
- 表结构定义在 `backend/manifest/sql/init.sql`；运行数据 `backend/data/adminbase.db`（gitignored）。

## 常用命令

- 一键起前后端：根目录 `make dev`（后端 :8000、前端 :8001；`gf` 已装则热重载，否则 `go run .`）。
- 后端验证：`cd backend && go vet ./... && go build ./...`（无测试框架、无 golangci）。
- 前端验证：`cd frontend && pnpm tsc && pnpm build`（`build` 已含 `tsc --noEmit`）。
- 依赖安装：`make install`；重置数据库：`make db-reset`（删 `data/adminbase.db` 后按 init.sql 重建）。
- 默认账号 `admin` / `123456`。

## 后端约定（易漏）

- 路由**不写在 controller**：写在 API 结构体的 `g.Meta` `path:"..." method:"..."`，再由 `internal/cmd/cmd.go` 的 `group.Bind(...)` 生效。
- 新增接口要动 5 处：`api/v1/**`（Req/Res + g.Meta）→ `controller/**`（先 `authz.Check(ctx, "权限码")` 再调 service）→ `service/**`（接口）→ `logic/**`（实现并在 `init()` 里 `RegisterXxx`）→ `cmd/cmd.go` `Bind`；新增 service 还要在 `internal/logic/logic.go` 空导入。
- 响应**一律 HTTP 200**，业务码在 body：`{code,message,data}`，`0` 成功，`401/402/403` 失败。错误用 `utility/errcode` 构造，`Response` 中间件自动包裹，别手拼。
- 鉴权分层：`middleware.Auth` 挂在受保护分组；细粒度权限用 `authz.Check/CheckAny`。`middleware.Permission` 是未被使用的死代码。
- 超管判定统一用 `appcfg.SuperAdminId(ctx)` 与 `consts.SuperRoleCode`，不要硬编码 id=1。

## 前端约定（易漏）

- **不要引入 `@umijs/max`，不要建 `config/routes.ts`**（Umi 已移除）。路由在 `src/router/index.tsx`，菜单在 `src/config/menu.tsx`，权限在 `src/access.tsx`（`useAccess` / `<Access>`）。
- 请求统一走 `src/api/client.ts`（axios 实例：注入 token、`401/402` 自动刷新并重放原请求）；页面调 `src/services/*`，服务端数据用 TanStack Query。
- 新增页面：`src/pages/...` → `src/router/index.tsx` 注册（页面级权限用 `<Permission perm="...">`）→ 需要侧边栏则加 `src/config/menu.tsx`。
- 提示/弹框用 `@/utils/antdApp` 的 `message`/`modal`（由 `AntdAppBridge` 绑定 antd `App` 上下文）；不要用 antd 静态 `message`/`Modal.confirm`，否则暗色主题下样式错乱。
- 主题相关颜色用 antd `theme.useToken()` 取，别写死 `rgba(0,0,0,...)`（暗色下不可读）。
- 端口/代理看 `vite.config.ts` 与 `.env.development`：`VITE_API_BASE`、`VITE_PORT`（默认 8001）、`VITE_PROXY_TARGET`。

## 数据库与代码生成

- 改表后 `cd backend && gf gen dao`（配置 `hack/config.yaml`），生成 `internal/dao`、`model/entity`、`model/do`，**不要手改生成物**。
- 生成后确认 `internal/dao/*.go` 的 import 为 `agigame/console/backend/internal/dao/internal`；gf 不覆盖已存在的外部 dao 文件，必要时删除后重生成。
- 启动只在 `sys_user` 表不存在时执行 init.sql；**已有库不会自动增量迁移**。加表/改表需 `make db-reset`（会清空数据）。
- 业务表是软删除（`deleted_at`）；带唯一键的表用 `WHERE deleted_at IS NULL` 的部分唯一索引，不要加普通唯一索引（否则“删后重建同名”会冲突）。

## 其它

- JWT 密钥可用环境变量 `ADMINBASE_JWT_SECRET` 覆盖（优先于 config.yaml）。
- 生成物/运行数据已 gitignore，勿提交：`frontend/dist`、`frontend/src/.umi*`、`backend/data`、`backend/logs`、`backend/main`、`backend/bin`。
- 前端无 ESLint/Prettier，类型正确性由 `tsc` 保证。
