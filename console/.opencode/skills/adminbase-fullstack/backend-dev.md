# 后端开发流程（GoFrame v2）

## 任务定位

面向 AdminBase 后端的完整流程：**读需求 → 需求确认 → 实现 → 接口规范 → 输出总结**。
框架细节以项目既有代码为准，**先阅读已有实现，再照猫画虎**，不自造风格。

---

## 第一步：读取需求与现状

按顺序读取（存在则读）：

1. `docs/产品文档.md`（默认路径，缺失则询问用户）
2. `docs/接口文档.md`（已有的接口约定）
3. 相关模块的既有代码：`backend/api/v1/**`、`controller/**`、`service/**`、`logic/**`、`internal/dao`、`manifest/sql/init.sql`
4. 项目约定：根 `AGENTS.md`、`README.md`

---

## 第二步：需求分析与确认（必须等待用户确认）

输出以下内容，等用户确认后再进入实现：

### 2.1 需求理解

- 列出本次要实现的功能点，标注优先级（核心 / 次要 / 可选）
- 有歧义的需求单列「待确认项」

### 2.2 技术方案

- 数据模型：表结构、字段、关联关系、索引（软删除 `deleted_at`；带唯一键用部分唯一索引）
- API 设计：路径、方法、请求参数、响应字段、权限码
- 核心业务逻辑与副作用（级联删除等）
- 涉及的文件清单（按下面「5 处改动」列出）

### 2.3 待确认项

逐条列出等用户确认。**未确认不动手。**

---

## 第三步：代码实现（新增接口要动 5 处）

按顺序改，缺一处会 panic（service 未注册时 `Xxx()` 会 panic `forgot register?`）。

### 1. `backend/api/v1/system/xxx.go` — Req/Res + `g.Meta`

```go
type XxxListReq struct {
    g.Meta   `path:"/system/xxx/list" method:"get" tags:"Xxx管理" summary:"Xxx列表"`
    PageNum  int    `json:"pageNum" d:"1"`
    PageSize int    `json:"pageSize" d:"10"`
    Name     string `json:"name"`
}
type XxxListRes struct {
    Total int           `json:"total"`
    List  []*entity.Xxx `json:"list"`
}
type XxxCreateReq struct {
    g.Meta `path:"/system/xxx" method:"post" tags:"Xxx管理" summary:"新增Xxx"`
    Name   string `json:"name" v:"required#请输入名称"`
}
type XxxCreateRes struct{}
```

- 路径参数：`Id int \`json:"id" in:"path" v:"required#ID不能为空"\``；默认值 `d:"1"`。
- 可选过滤用指针区分「未传」与零值：`Status *int`。
- 特殊动作直接拼路径，如 `/system/user/resetPwd`；批量删除收 `Ids []int`。

### 2. `backend/internal/controller/system/xxx.go` — 先 `authz.Check` 再调 service（不拼 JSON）

```go
type cXxx struct{}
func NewXxx() *cXxx { return &cXxx{} }

func (c *cXxx) List(ctx context.Context, req *api.XxxListReq) (*api.XxxListRes, error) {
    if err := authz.Check(ctx, "system:xxx:list"); err != nil {
        return nil, err
    }
    total, list, err := service.Xxx().List(ctx, &model.XxxQueryInput{
        PageInput: model.PageInput{PageNum: req.PageNum, PageSize: req.PageSize},
    })
    if err != nil {
        return nil, err
    }
    return &api.XxxListRes{Total: total, List: list}, nil
}
```

- 权限码命名：`system:<模块>:list|add|edit|remove`；多权限用 `authz.CheckAny(ctx, "a", "b")`。

### 3. `backend/internal/service/xxx.go` — 接口 + 全局 getter + Register

```go
type IXxx interface { List(ctx context.Context, in *model.XxxQueryInput) (int, []*entity.Xxx, error) }
var localXxx IXxx
func Xxx() IXxx {
    if localXxx == nil { panic("implement not found for interface IXxx, forgot register?") }
    return localXxx
}
func RegisterXxx(i IXxx) { localXxx = i }
```

### 4. `backend/internal/logic/xxx/xxx.go` — 实现 + `init()` 注册

```go
type sXxx struct{}
func init() { service.RegisterXxx(New()) }
func New() *sXxx { return &sXxx{} }

// DAO：dao.Xxx.Ctx(ctx).Where/WhereLike/Page/Count/Scan；写库用 do.Xxx{...}
// 列名用 dao.Xxx.Columns().Xxx；软删除由 GF 自动处理（Delete 即软删）
// 业务错误返回 errcode.NotFound/BadRequest/Forbidden(...)
```

### 5. `backend/internal/cmd/cmd.go` — 在受保护分组的 `protected.Bind(...)` 里加 `system.NewXxx()`。

### 6.（仅当新增 logic 包）在 `backend/internal/logic/logic.go` 加空导入 `_ ".../server/internal/logic/xxx"`。

### 实现原则

- **遵循项目惯例**，优先复用既有 `logic`/`utility`，最小改动，不做无关重构。
- 职责分离：controller 只做权限校验与参数转换，业务逻辑放 logic。
- **不手改生成物**（`internal/dao`、`dao/internal`、`model/entity`、`model/do`）。

---

## 第四步：响应、错误码与权限

- 响应**一律 HTTP 200**，业务码在 body `{code,message,data}`；`middleware.Response` 自动包裹，成功固定 `code:0`，**别手拼 JSON**。
- 错误用 `utility/errcode`：`Business`(1) / `BadRequest`(400) / `Unauthorized`(401) / `TokenExpired`(402) / `Forbidden`(403) / `NotFound`(404)。
- `middleware.Auth` 只挂在受保护分组；公开接口（登录/刷新）单独 `group.Bind`。

---

## 第五步：数据库与代码生成

- 改表：编辑 `manifest/sql/init.sql` → `make db-reset`（清数据重建）→ `cd backend && gf gen dao`（配置 `backend/hack/config.yaml`，`removePrefix: sys_`，`withTime: true`）。
- 生成 `internal/dao`、`dao/internal`、`model/entity`、`model/do`；dao import 前缀为 `.../server/internal/dao/internal`。
- 业务表软删除（`deleted_at`）；**带唯一键的表用部分唯一索引** `... (col) WHERE deleted_at IS NULL`，不要普通唯一索引（否则删后重建同名冲突）。
- 启动只在 `sys_user` 表不存在时执行 init.sql；**已有库不会自动增量迁移**，加/改表必须 `make db-reset`。

---

## 第六步：接口规范（建议同步，不强制维护文件）

若项目维护 `docs/接口文档.md`，每个接口按以下格式同步；未维护文件时，至少保证 `g.Meta` 的 `tags`/`summary` 准确、Swagger 可用。

1. **接口路径** — `METHOD /api/v1/xxx/yyy`
2. **请求参数** — 表格（参数名、类型、是否必填、说明）
3. **请求示例** — JSON
4. **响应示例** — JSON（成功 + 失败）
5. **响应字段说明** — 表格（字段名、类型、说明）
6. **副作用说明** — 如有级联操作必须写明

同步规则：新增追加、修改更新、删除移除；按模块分组；示例与代码逻辑一致。

---

## 第七步：输出总结

- **已实现功能**：功能点 + 对应文件
- **接口变更**：新增 / 修改 / 删除的接口
- **数据库变更**：新增 / 修改的表和字段
- **待处理项**：未完成或后续事项

---

## 迭代开发

用户提出新需求或修改时重复本流程：重读需求与 `docs/接口文档.md`（可能已更新）→ 分析增量 → 确认方案 → 实现 → 校正接口描述。

## 验证

```bash
cd backend && go vet ./... && go build ./...
```
