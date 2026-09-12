# GoBoy 控制台

GoBoy-LLM 的控制台，基于 **AdminBase**（GoFrame v2 + Vite/React + Ant Design + SQLite）改造适配，内置 RBAC 权限、部门、字典、参数配置、登录/操作日志等基础模块；将新增模拟器会话/实验管理并与 `agigame/emulator/engine` 对接。

- 后端：http://127.0.0.1:8000
- 前端（dev）：http://127.0.0.1:8001（`/api` 代理到后端）
- 默认账号：`admin` / `123456`

---

## 一、技术选型

### 后端

| 领域 | 选型 | 版本 | 说明 |
| --- | --- | --- | --- |
| 语言 | Go | 1.26 | |
| 应用框架 | GoFrame | v2.10.3 | 一体化框架：`ghttp` 路由、`gdb` ORM、`gcmd` 命令、`glog` 日志、`gcfg` 配置，无需再拼装第三方组件 |
| 数据库 | SQLite | 3.x | 单机部署零依赖，数据文件落在运行目录 `data/adminbase.db` |
| 驱动 | `glebarez/go-sqlite`（经 GoFrame `contrib/drivers/sqlite`） | v1.21.2 | 纯 Go 实现，**无需 CGO/gcc**，跨平台编译简单 |
| 认证 | JWT（`golang-jwt/jwt/v5`） | v5.3.1 | access + refresh 双令牌 |
| 密码 | `golang.org/x/crypto/bcrypt` | v0.57.0 | 单向哈希存储 |
| 代码生成 | `gf gen dao` | — | 由表结构生成 `dao / entity / do` 分层代码 |

### 前端

| 领域 | 选型 | 版本 | 说明 |
| --- | --- | --- | --- |
| 构建/开发服务器 | Vite | 6.x | 显式配置，启动快，`/api` 开发代理 |
| 框架 | React | 19.x | 仅用 React 本身，无上层“全家桶”框架 |
| 路由 | react-router-dom | 6.x | 显式路由表 + 登录/权限守卫 |
| 服务端状态 | TanStack Query | 5.x | 请求缓存、重试、失效、loading 全托管 |
| 客户端状态 | Zustand | 5.x | 登录态（token / refreshToken） |
| 组件库 | antd | 6.x | |
| Pro 组件 | `@ant-design/pro-components` | 3.x | `ProLayout`、`ProTable`、`ProForm` 等（按需使用） |
| 图表 | `@ant-design/plots` | 2.x | |
| HTTP | axios | 1.x | 统一实例 + 拦截器（注入 token / 401 无感刷新） |
| 样式 | Tailwind CSS | 4.x | 与 antd 并存 |
| 语言/包管理 | TypeScript 5 / pnpm | 5.x / — | |

> ⚠️ **版本说明**：前端当前采用 antd 6 + React 19，配套的 `@ant-design/pro-components` 稳定版仅支持 antd 5，故使用了 **`3.1.14-x` 预发布版**。若追求最大稳定性，可整体回退到 `React 18 + antd 5 + pro-components 2.8.x` 组合（见「十、注意事项」）。

### 选型取向

- **显式优先**：后端用 GoFrame 一体化基建；前端刻意不用 Umi 这类“约定式全家桶”，路由、请求、状态、权限都显式装配，控制流可读可控。
- **零外部依赖**：SQLite 纯 Go 驱动，无需单独部署数据库，便于本地开发与小型单机部署。
- **标准协议**：REST + JSON + JWT，前后端职责清晰。

---

## 二、整体架构

```
┌──────────────────────────────────┐         ┌──────────────────────────────────────┐
│  浏览器 (React + Vite + antd)     │  HTTP   │  后端 (GoFrame :8000)                 │
│  pages → hooks/hooks → api/client │ ──────► │  /api/v1/**                           │
│  router / stores / access        │ ◄────── │  统一响应 {code,message,data}          │
└──────────────────────────────────┘         └──────────────────────────────────────┘
          │ dev: Vite proxy /api → :8000
          │ prod: nginx  location /api/ → :8000
```

前后端**完全分离**：前端只通过 HTTP 调 `/api/v1`，不共享代码；开发用 Vite 代理、生产用 nginx 反代解决同源问题。

### 后端分层

```
main.go
  └─ internal/cmd            启动入口：初始化数据库 + 注册路由
       ├─ boot               首次启动建库 + 种子数据 + 索引迁移
       ├─ api/v1/*           接口出入参（g.Meta 声明 path/method/校验规则）
       ├─ controller/*       控制器：参数绑定后的业务编排 + 权限校验（authz.Check）
       ├─ service/*          服务接口定义 + 注册/获取（依赖倒置）
       ├─ logic/*            业务逻辑实现，通过 init() 注册到 service
       ├─ dao/*              数据访问（gf gen dao 生成，internal/ 为底层）
       ├─ model/
       │    ├─ entity/do     表实体 / 数据对象（生成）
       │    └─ *.go          业务输入输出、上下文（ContextUser）、统一响应
       ├─ middleware/        统一响应、异常恢复、JWT 鉴权、权限、操作日志
       └─ utility/           jwtutil / password / errcode / ctxuser / authz / appcfg
```

**依赖方向**：`controller → service(接口) → logic(实现) → dao → DB`。`logic` 包在 `init()` 里注册实现，`cmd` 通过空导入 `internal/logic` 触发注册（见 `internal/logic/logic.go`），`controller` 只依赖 `service` 接口，从而实现解耦与可测试性。

### 请求处理管道

```
ghttp
  └─ CORS → Response(统一包体) → ErrorHandler(panic 恢复) → Auth(JWT) → OperLog
       → controller → service → logic → dao → SQLite
```

- `Response` 把控制器返回值/错误统一包成 `{code, message, data}`；
- `ErrorHandler` 兜底 panic；
- `Auth` 解析 JWT 并组装 `ContextUser`（角色、权限、是否超管）写入请求上下文；
- `OperLog` 记录写操作（POST/PUT/DELETE）。

### 前端架构

```
index.html / vite.config.ts        Vite 入口与配置（别名 @、/api 开发代理、分包）
.env.development / .env.production 环境变量（VITE_API_BASE 等）
src/
  ├─ main.tsx            应用入口：挂载 <App/>
  ├─ App.tsx             Provider 组装：ErrorBoundary → QueryClient → ConfigProvider → BrowserRouter
  ├─ router/
  │    ├─ index.tsx       显式路由表（React.lazy 懒加载，demo 仅开发环境）
  │    ├─ RequireAuth.tsx 登录守卫：无 token / 拉用户失败 → 跳登录
  │    └─ Permission.tsx  页面级权限守卫：无权限 → 403
  ├─ layouts/BasicLayout.tsx  ProLayout + 侧边菜单(按权限过滤) + SettingDrawer
  ├─ access.tsx          AccessProvider / useAccess / <Access>（由 perms 生成开关）
  ├─ api/
  │    ├─ client.ts        axios 实例 + 拦截器（注入 token / 统一错误 / 401 无感刷新）
  │    └─ queryClient.ts   TanStack QueryClient
  ├─ stores/auth.ts      Zustand：token / refreshToken（localStorage 持久化）
  ├─ hooks/              useCurrentUser（用户信息查询）/ useChartTheme
  ├─ services/           API 封装（auth.ts / system.ts / types.ts）
  ├─ config/             defaultSettings.ts / menu.tsx（侧边菜单配置）
  ├─ components/         AvatarDropdown / AccountModal / ErrorBoundary / ChartCard
  ├─ utils/history.ts    非组件模块可用的导航封装
  ├─ theme.ts            布局/主题设置持久化（useSyncExternalStore）
  └─ pages/              Login / Dashboard / system/* / demo/*
```

---

## 三、前后端通讯机制

### 1. 请求寻址

| 环节 | 开发环境 | 生产环境 |
| --- | --- | --- |
| 前端 baseURL | `/api/v1`（`src/api/client.ts`，由 `VITE_API_BASE` 注入） | 同左 |
| 转发 | Vite dev proxy：`/api` → `http://127.0.0.1:8000`（`vite.config.ts`） | nginx：`location /api/` → `http://127.0.0.1:8000` |
| 后端挂载 | `s.Group("/api/v1", ...)`（`internal/cmd/cmd.go`） | 同左 |

浏览器最终访问 `/api/v1/system/user/list`，由代理转发到后端同名路径。

### 2. 路由注册

后端路由**不写在控制器里**，而是写在各接口结构体的 `g.Meta` 上；`Bind` 反射生成路由：

```go
// api/v1/system/user.go
type UserListReq struct {
    g.Meta `path:"/system/user/list" method:"get" ...`
}

// cmd/cmd.go
protected.Bind(system.NewUser())
```

### 3. 统一响应

所有接口（含错误）一律返回 HTTP 200，业务状态由 `code` 表达：

```json
{ "code": 0, "message": "success", "data": {} }
```

| code | 含义 |
| --- | --- |
| `0` | 成功 |
| `400` | 参数/业务错误 |
| `401` | 未登录 |
| `402` | 令牌过期 |
| `403` | 无权限 |
| `404` | 资源不存在 |

分页数据形如 `data: { total: 100, list: [] }`。

### 4. 认证与鉴权流程

1. **登录** `POST /auth/login` → 返回 `access token` + `refresh token`；
2. 前端存 `localStorage`，请求拦截器自动加 `Authorization: Bearer <token>`；
3. 后端 `Auth` 中间件解析 JWT → `GetContextUser` 组装角色/权限 → 写入请求上下文；
4. 控制器用 `authz.Check(ctx, "system:user:add")` 校验权限码；
5. 前端 `access.tsx`（`useAccess` / `<Access>`）根据 `perms` 控制菜单与按钮显隐；
6. **无感刷新**：响应 `code` 为 `401/402` 时，`src/api/client.ts` 调用 `/auth/refresh` 换新令牌并重放原请求，失败才跳登录。

---

## 四、目录结构

```
adminbase/
├── backend/                       # 后端 (GoFrame)
│   ├── api/v1/                    # 接口定义：auth/ system/
│   ├── internal/
│   │   ├── boot/                  # 数据库初始化与索引迁移
│   │   ├── cmd/                   # 启动入口与路由注册
│   │   ├── controller/            # 控制器（权限校验 + 编排）
│   │   ├── logic/                 # 业务逻辑实现
│   │   ├── service/               # 服务接口（注册/获取）
│   │   ├── dao/                   # 数据访问（生成）
│   │   ├── model/                 # entity / do（生成）+ 业务模型
│   │   ├── consts/                # 常量
│   │   └── middleware/            # 响应/异常/JWT/权限/操作日志
│   ├── utility/                   # jwtutil / password / errcode / ctxuser / authz / appcfg
│   ├── manifest/
│   │   ├── config/config.yaml     # 主配置
│   │   └── sql/init.sql           # 建表 + 种子数据
│   ├── hack/config.yaml           # gf gen dao 配置
│   ├── main.go
│   └── Makefile
├── frontend/                      # 前端 (Vite + React + antd)
│   ├── vite.config.ts             # 别名 / 开发代理 / 分包
│   ├── src/                       # router / api / stores / layouts / pages ...
│   └── package.json
└── Makefile                       # 前后端统一命令
```

---

## 五、数据模型

| 表 | 用途 |
| --- | --- |
| `sys_user` | 用户 |
| `sys_role` | 角色（含数据权限范围 `data_scope`） |
| `sys_menu` | 菜单/目录/按钮权限（`type` = `M`/`C`/`F`） |
| `sys_user_role` | 用户-角色关联 |
| `sys_role_menu` | 角色-菜单关联 |
| `sys_dept` | 部门（`ancestors` 存路径，树形） |
| `sys_dict_type` / `sys_dict_data` | 字典类型 / 字典数据 |
| `sys_config` | 系统参数 |
| `sys_login_log` | 登录日志 |
| `sys_oper_log` | 操作日志 |

> 所有业务表带 `created_at / updated_at / deleted_at`，删除为**软删除**；带唯一键的表使用 `WHERE deleted_at IS NULL` 的部分唯一索引，避免“删除后重建同名”冲突。

---

## 六、快速开始

### 1. 一键启动

```bash
make dev          # 后端 :8000（热加载）+ 前端 :8001（HMR），Ctrl+C 同时退出
```

### 2. 分开启动

```bash
# 后端
cd backend
go mod tidy
go run .          # 首次启动自动建库并导入种子数据

# 前端
cd frontend
pnpm install
pnpm dev          # Vite 默认 :8001，/api 代理到后端 :8000（可用 VITE_PORT 覆盖）
```

- Swagger：http://127.0.0.1:8000/swagger
- OpenAPI：http://127.0.0.1:8000/api.json

### 3. 默认账号

| 用户名 | 密码 | 角色 |
| --- | --- | --- |
| `admin` | `123456` | 超级管理员（拥有全部权限） |

### 4. 其它常用命令

```bash
make install      # 安装前后端依赖
make build        # 构建前后端产物
make db-init      # 手动初始化数据库（需 sqlite3）
make db-reset     # 重置数据库
make tidy         # go mod tidy
make clean        # 清理构建产物与运行数据
```

---

## 七、开发约定

### 新增一个后端接口

1. 在 `api/v1/system/xxx.go` 定义 `XxxReq` / `XxxRes`，`g.Meta` 声明 `path/method`；
2. 在 `controller/system/xxx.go` 实现方法：先 `authz.Check(ctx, "权限码")`，再调 `service`；
3. 在 `service/xxx.go` 定义接口方法，`logic/xxx/xxx.go` 实现并 `service.RegisterXxx`；
4. 若新增了服务，在 `logic/logic.go` 空导入；控制器在 `cmd/cmd.go` 中 `Bind`。

### 数据库结构变更

```bash
cd backend
gf gen dao        # 依据 hack/config.yaml 重新生成 dao/entity/do
```

> 生成后确认 `internal/dao/*.go` 的 import 路径为 `.../server/internal/dao/internal`；gf 在文件已存在时不会覆盖外部 dao 文件，必要时删除后重新生成。

### 统一响应与错误

- 控制器返回 `(res, err)`；`err` 用 `utility/errcode` 构造，会自动带业务码；
- 成功由 `Response` 中间件自动包裹，无需手动拼 `code/message`。

### 新增一个前端页面

1. 在 `src/pages/xxx/index.tsx` 写页面组件（权限按钮用 `@/access` 的 `<Access>` / `useAccess`）；
2. 在 `src/router/index.tsx` 注册路由，页面级权限用 `<Permission perm="...">` 包裹；
3. 若需出现在侧边栏，在 `src/config/menu.tsx` 增加菜单项（`access` 字段对应权限码）；
4. 接口调用统一走 `src/services/*`（基于 `src/api/client.ts`），服务端状态用 TanStack Query。

---

## 八、配置说明

`backend/manifest/config/config.yaml`：

| 配置 | 说明 |
| --- | --- |
| `server.address` | 监听地址，默认 `:8000` |
| `database.default.link` | `sqlite::@file(./data/adminbase.db)` |
| `database.default.debug` | 是否打印 SQL 调试日志 |
| `jwt.secret` | JWT 密钥，**生产环境务必修改**，也可用环境变量 `ADMINBASE_JWT_SECRET` 覆盖 |
| `jwt.expire` / `jwt.refreshExpire` | access / refresh 有效期（秒） |
| `system.superAdminId` | 超级管理员用户 ID |
| `system.defaultPassword` | 新建/重置用户默认密码 |

前端环境变量（`frontend/.env.development` / `.env.production`，仅 `VITE_` 前缀会暴露）：

| 变量 | 说明 |
| --- | --- |
| `VITE_API_BASE` | 接口前缀，默认 `/api/v1` |
| `VITE_PORT` | 开发端口，默认 `8001` |
| `VITE_PROXY_TARGET` | 开发代理目标，默认 `http://127.0.0.1:8000` |

---

## 九、部署（nginx 前后端分离）

```bash
cd frontend && pnpm build          # 产物在 frontend/dist
cd ../backend && go build -o bin/adminbase .
```

```nginx
server {
    listen 80;
    server_name admin.example.com;

    root /var/www/adminbase/dist;
    index index.html;
    location / {
        try_files $uri $uri/ /index.html;   # history 路由
    }

    location /api/ {
        proxy_pass http://127.0.0.1:8000;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

后端使用 SQLite，直接运行二进制即可（数据库文件位于运行目录 `data/adminbase.db`）。

---

## 十、注意事项与已知限制

- **前端版本偏前沿**：antd 6 + React 19 + `pro-components` 预发布版，稳定性有风险；如需保守，可整体回退到 React 18 + antd 5 + `pro-components` 2.8.x。
- **令牌撤销为进程内存**：登出/令牌轮换的黑名单不持久化，重启或横向扩容后失效，多实例部署需改用 Redis 等共享存储。
- **HTTP 状态码恒为 200**：错误仅体现在响应体 `code`，不利于网关/监控按状态码告警，接入监控时需注意。
- **数据权限未落地**：角色 `data_scope` 仅存储，列表查询尚未按范围过滤。
- **动态菜单未接入侧边栏**：菜单管理维护的数据当前未驱动前端导航，侧边栏来自静态 `src/config/menu.tsx`。
- **数据库迁移简易**：目前仅靠 `init.sql` + 启动时的幂等索引迁移，无版本化迁移框架。
- **无自动化测试**：尚未接入单元测试与 CI。
- **代理取 IP**：经 nginx 部署时需在后端配置 `server.clientIpHeader`，否则日志记录的是代理 IP。

---

## 常见问题

- **pnpm 提示 `Ignored build scripts`**：执行 `pnpm approve-builds --all` 允许 esbuild 等构建脚本。
- **React 19 `useRef` 需初值**：`useRef<ActionType | undefined>(undefined)`。
- **开发端口被占用**：Vite 会自动顺延端口，可用 `VITE_PORT=8002 pnpm dev` 指定。
