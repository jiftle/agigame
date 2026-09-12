package model

// DictTypeQueryInput 字典类型查询
type DictTypeQueryInput struct {
	PageInput
	Name   string
	Type   string
	Status *int
}

// DictTypeSaveInput 字典类型新增/修改
type DictTypeSaveInput struct {
	Id     int
	Name   string
	Type   string
	Status int
	Remark string
}

// DictDataQueryInput 字典数据查询
type DictDataQueryInput struct {
	PageInput
	DictType  string
	DictLabel string
	Status    *int
}

// DictDataSaveInput 字典数据新增/修改
type DictDataSaveInput struct {
	Id        int
	DictSort  int
	DictLabel string
	DictValue string
	DictType  string
	IsDefault int
	Status    int
	Remark    string
}
