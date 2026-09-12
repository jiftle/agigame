package model

// DeptQueryInput 部门查询
type DeptQueryInput struct {
	Name   string
	Status *int
}

// DeptSaveInput 部门新增/修改
type DeptSaveInput struct {
	Id       int
	ParentId int
	Name     string
	Leader   string
	Phone    string
	Email    string
	Sort     int
	Status   int
	Remark   string
}

// DeptNode 部门树节点
type DeptNode struct {
	Id       int         `json:"id"`
	ParentId int         `json:"parentId"`
	Name     string      `json:"name"`
	Leader   string      `json:"leader"`
	Phone    string      `json:"phone"`
	Email    string      `json:"email"`
	Sort     int         `json:"sort"`
	Status   int         `json:"status"`
	Children []*DeptNode `json:"children,omitempty"`
}
