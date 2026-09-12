package model

// Response 统一响应结构
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// PageInput 通用分页入参
type PageInput struct {
	PageNum  int `json:"pageNum"  d:"1"`
	PageSize int `json:"pageSize" d:"10"`
}

// Offset 计算偏移量
func (p PageInput) Offset() int {
	if p.PageNum < 1 {
		p.PageNum = 1
	}
	return (p.PageNum - 1) * p.PageSize
}

// Limit 返回分页大小
func (p PageInput) Limit() int {
	if p.PageSize < 1 {
		p.PageSize = 10
	}
	if p.PageSize > 1000 {
		p.PageSize = 1000
	}
	return p.PageSize
}
