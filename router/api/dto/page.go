// Package model
// @author： Boice
// @createTime：2022/5/27 16:29
package dto

type Page struct {
	PageNo   int `json:"page_no"`
	PageSize int `json:"page_size"`
}

type PageResult struct {
	Total int64       `json:"total"`
	Data  interface{} `json:"data"`
}

func DefaultPage() *Page {
	return &Page{
		PageNo:   1,
		PageSize: 10,
	}
}

func NewPageResult(total int64, data interface{}) *PageResult {
	return &PageResult{
		total,
		data,
	}

}
