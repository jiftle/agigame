---
name: adminbase-fullstack
description: >
  AdminBase 全栈开发流程入口（GoFrame v2 后端 + Vite/React 19/antd 6 前端）。
  Use when 在本项目开发功能并需要按「需求确认 → 实现 → 接口对齐」推进时：新增/修改后端接口、
  controller/service/logic/dao、数据库表，新增前端页面/路由/菜单/权限，或任务涉及
  接口文档 / 产品文档 / 全栈开发 / 需求分析 / 参数管理 / 系统管理。
  触发词：新增接口、新增页面、全栈开发、接口文档、需求分析、后端开发、前端开发。
---

# AdminBase 全栈开发流程

本项目技术栈已固定（GoFrame v2 + Vite/React 19/antd 6 + SQLite），**无需运行时检测项目类型**。
本 skill 负责把开发任务路由到后端 / 前端流程，并统一「先确认再动手、前后端以接口对齐」的开发节奏。

内部资源（与本文件同目录）：

| 文件 | 说明 |
| --- | --- |
| `backend-dev.md` | 后端流程：读需求 → 确认方案 → 实现 → 接口规范 → 总结 |
| `frontend-dev.md` | 前端流程：读接口 → 确认规范 → 确认方案 → 实现 → 总结 |
| `code-style.md` | 前端代码规范（React + TS + antd），被 `frontend-dev.md` 引用 |

配套技能：写 antd 代码前先查再写，用 `antd` skill（`antd info/demo/token/lint`）。硬性技术约束以项目根 `AGENTS.md` 与 `README.md` 为准。

---

## 第一步：确认本次开发方向

询问用户本次工作的范围（除非用户已明确指定）：

> 本次开发聚焦哪个方向？
> 1. 后端开发（进入 `backend-dev.md`）
> 2. 前端开发（进入 `frontend-dev.md`）
> 3. 全栈开发（先后端、再前端，串联执行）

- 用户已说「改接口 / 加表 / service」等 → 直接走后端。
- 用户已说「加页面 / 菜单 / 表单」等 → 直接走前端。
- 用户描述一个完整功能（既有接口又有页面）→ 询问，或默认全栈。

---

## 第二步：确认文档位置

工作流的文档约定（**仅作为参考规范，不强制维护文件**）：

- 产品文档：`docs/产品文档.md`
- 接口文档：`docs/接口文档.md`

规则：

1. 文件存在就读；不存在就向用户确认路径，或直接基于代码与需求推进。
2. 接口契约优先以代码为准：`backend/api/v1/**` 的 `g.Meta`（`path/method/tags/summary`）+ Swagger（`http://127.0.0.1:8000/swagger`）即真实契约。
3. 若项目正在维护 `docs/接口文档.md`，后端实现后**建议**同步；未维护时不阻塞流程，但保持 `g.Meta` 描述准确。

---

## 第三步：路由到对应流程

读取对应资源文件，从「第一步」开始完整执行：

| 方向 | 资源 |
| --- | --- |
| 后端 | `backend-dev.md` |
| 前端 | `frontend-dev.md` |
| 全栈 | 先执行 `backend-dev.md`（产出/校正接口），再执行 `frontend-dev.md` |

前端流程引用代码规范时，默认读取同目录 `code-style.md`。

---

## 项目坐标速查

- 后端 module：`agigame/console/backend`。
- 后端入口：`backend/main.go` → `internal/cmd`（先建库、再注册路由）。
- 路由**不写在 controller**，写在 API 结构体 `g.Meta` 的 `path`/`method`，靠 `cmd.go` 的 `group.Bind(...)` 生效。
- 表结构：`backend/manifest/sql/init.sql`；运行库 `backend/data/adminbase.db`（gitignored）。
- 超管判定统一 `appcfg.SuperAdminId(ctx)` 与 `consts.SuperRoleCode`；`middleware.Permission` 是死代码。
- 前端路由 `frontend/src/router/index.tsx`，菜单 `frontend/src/config/menu.tsx`，权限 `frontend/src/access.tsx`。
- 前端请求统一走 `frontend/src/api/client.ts`；**不要引入 `@umijs/max`、不要建 `config/routes.ts`**（Umi 已移除）。

## 验证命令（提交前必跑）

- 后端：`cd backend && go vet ./... && go build ./...`
- 前端：`cd frontend && pnpm tsc && pnpm build`（`build` 已含 `tsc --noEmit`）
- 一键起：根目录 `make dev`（后端 :8000、前端 :8001）；默认账号 `admin` / `123456`。
- 前端无 ESLint/Prettier，类型正确性由 `tsc` 保证。

## 界面验证协议

开发环境侧边栏「示例页面」下有两个组件验证页（`import.meta.env.DEV` 限定，生产构建剔除）：

- `/demo/overview` 组件总览：按通用 / 数据录入 / 数据展示 / 反馈 / 导航 / Pro 组件分区，逐块核对渲染与配色。
- `/demo/status` 状态与反馈：异步三态、权限（`useAccess` / `<Access>`）、`ErrorBoundary`、403 预览。

触发时机：升级 antd / pro-components、修改主题 token 或明暗逻辑、调整全局样式后。

步骤：

1. `cd frontend && pnpm tsc && pnpm build`
2. `antd lint <改动路径>`，确保无 deprecated / 错误用法
3. `make dev`，登录后打开上述两页，切换亮 / 暗与主色，逐块检查无错位、无写死配色、无控制台报错
4. 引入新的组件用法时，同步补进组件总览页

