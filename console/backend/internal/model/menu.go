package model

// MenuQueryInput 菜单查询
type MenuQueryInput struct {
	Title  string
	Status *int
}

// MenuSaveInput 菜单新增/修改
type MenuSaveInput struct {
	Id        int
	ParentId  int
	Title     string
	Name      string
	Path      string
	Component string
	Icon      string
	Type      string
	Perms     string
	Sort      int
	Visible   int
	Status    int
	Redirect  string
	IsFrame   int
	IsCache   int
}
