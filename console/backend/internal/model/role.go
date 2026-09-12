package model

// RoleQueryInput 角色查询
type RoleQueryInput struct {
	PageInput
	Name   string
	Code   string
	Status *int
}

// RoleSaveInput 角色新增/修改
type RoleSaveInput struct {
	Id        int
	Name      string
	Code      string
	Sort      int
	DataScope int
	Status    int
	Remark    string
	MenuIds   []int
}
