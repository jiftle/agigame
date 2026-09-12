package consts

// 上下文键
const (
	CtxUserKey = "AdminBaseCtxUser"
)

// 响应码
const (
	CodeSuccess      = 0
	CodeError        = 1
	CodeBadRequest   = 400
	CodeUnauthorized = 401
	CodeTokenExpired = 402
	CodeForbidden    = 403
	CodeNotFound     = 404
)

// 通用状态
const (
	StatusEnabled  = 1
	StatusDisabled = 0
)

// 菜单类型
const (
	MenuTypeDir    = "M" // 目录
	MenuTypeMenu   = "C" // 菜单
	MenuTypeButton = "F" // 按钮
)

// 数据权限范围
const (
	DataScopeAll       = 1 // 全部数据
	DataScopeCustom    = 2 // 自定义
	DataScopeDept      = 3 // 本部门
	DataScopeDeptChild = 4 // 本部门及以下
	DataScopeSelf      = 5 // 仅本人
)

// 内置
const (
	SuperAdminId   = 1
	SuperRoleCode  = "admin"
	SuperAdminName = "超级管理员"
)

// 性别
const (
	SexUnknown = 0
	SexMale    = 1
	SexFemale  = 2
)
