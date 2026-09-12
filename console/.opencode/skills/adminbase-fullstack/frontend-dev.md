# 前端开发流程（Vite + React 19 + antd 6）

## 任务定位

面向 AdminBase 前端的完整流程：**读接口 → 确认规范 → 需求确认 → 实现 → 输出总结**。
默认遵循同目录 `code-style.md`；项目已有规范（`AGENTS.md`、`README.md`）优先。

---

## 第零步：确认接口来源

接口契约优先以代码为准，按以下顺序获取：

1. `docs/接口文档.md`（若项目维护）
2. `backend/api/v1/**` 的 `g.Meta` 与 Req/Res 结构体（**真实契约**）
3. Swagger：`http://127.0.0.1:8000/swagger`

**每次收到新需求或需求变更，都重新读取接口，因为后端可能已改。**

---

## 第一步：确认代码规范

默认读取同目录 `code-style.md`。若项目存在 `AGENTS.md` 或 `README.md` 中的前端约定，以其为准；用户可指定不使用规范约束。

---

## 第二步：需求分析与确认（确认后再实现）

### 2.1 页面结构

- 需要新建 / 修改的页面与组件
- 布局与交互流程
- 需要调用的接口列表

### 2.2 接口对齐

- 列出涉及接口，核对请求参数与响应字段是否满足页面
- 不满足则标注「需要后端补充」并告知用户

### 2.3 待确认项

交互细节、字段不匹配、其他需用户拍板的内容，逐条列出。

---

## 第三步：代码实现（新增页面的 4 处改动）

1. **`frontend/src/services/<模块>.ts`** — 用 `request`（`@/api/client`）调接口，返回 `res.data`；类型用 `ApiResult<T>` / `PageResult<T>`（`services/types.ts`）。参照 `services/system.ts` 的写法。
2. **`frontend/src/pages/<模块>/index.tsx`** — 列表页惯用 `ProTable` + `request` 回调 + `actionRef.reload()`；弹窗用 `ModalForm`。**不要臆造所有页面都 `useQuery`**（目前仅 `useCurrentUser` 用 `useQuery`）。
3. **`frontend/src/router/index.tsx`** — 页面级权限用 `guard('system:xxx:list', <XxxPage />)`（内部是 `<Permission perm>`）。
4. **`frontend/src/config/menu.tsx`** — 加菜单项，带 `access: 'system:xxx:list'`（无权限自动隐藏）。**不要建 `config/routes.ts`，不要引入 `@umijs/max`**。

按钮/操作项权限：`useAccess()` + `<Access accessible={access['system:xxx:edit']}>`（`src/access.tsx`），不是 `<Permission>`。

### 实现顺序

1. 类型定义（接口请求 / 响应）
2. API 封装（`services/<模块>.ts`）
3. 页面主结构（`pages/<模块>/index.tsx`）
4. 页面级子组件（弹窗、表单、列表）
5. 路由 / 菜单 / 权限注册

### 实现原则

- **遵循项目惯例**：先读同类页面（如 `pages/system/user/index.tsx`、`role/index.tsx`），保持一致。
- **优先复用**：`ProTable` / `ModalForm` / `ChartCard` / `@/utils/antdApp` 等既有能力，不重复造轮子。
- **接口对齐**：参数与响应严格按后端结构，不擅自改字段名或结构。
- **类型安全**：TS 类型完整，避免 `any`。
- **最小改动**：只做需求要求的功能。
- **不引入新依赖**（除非与用户确认）。

### 请求与 UI 约定

- 请求统一走 `src/api/client.ts`（注入 token；`401/402` 自动刷新并重放）；**不要另建 axios 实例**。
- 提示/弹框用 `@/utils/antdApp` 的 `message`/`modal`（`AntdAppBridge` 绑定 App 上下文）；**不要用 antd 静态 `message`/`Modal.confirm`**（暗色主题样式错乱）。
- 主题色用 antd `theme.useToken()` 取，别写死 `rgba(0,0,0,...)` / `#ff4d4f`。
- 写 antd 代码前**先查再写**（`antd` skill 的 `antd info/demo/token/lint`），不要凭记忆。

---

## 第四步：输出总结

- **已实现功能**：功能点 + 对应文件
- **调用的接口**：接口及用途
- **未实现内容**：未完成或后续事项
- **接口问题**：接口文档 / 结构与需求不匹配之处，建议后端调整

---

## 迭代开发

用户提出新需求或修改时：

1. **重新读取接口**（后端可能已更新）
2. 分析增量需求 → 与用户确认方案 → 实现
3. 若接口变更影响已有代码，主动提醒并修复

## 接口变更检测

每次重读接口后对比上次：

- 新增接口 → 提醒是否使用
- 参数/响应变更 → 检查是否影响已实现页面，如有影响则提醒并修复
- 接口删除 → 检查是否破坏已有功能，如有则提醒用户

## 验证

```bash
cd frontend && pnpm tsc && pnpm build
```
