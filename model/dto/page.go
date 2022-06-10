package dto

type Pagination struct {
	Total    int64 `json:"total" form:"total"`
	Perpage  int64 `json:"per_page" form:"per_page"`
	Page     int64 `json:"page" form:"page"`
	LastPage int64 `json:"last_page" form:"last_page"`
	From     int64 `json:"from" form:"from"`
	To       int64 `json:"to" form:"to"`
}

type PageResult struct {
	Pagination *Pagination `json:"pagination"`
	PageData   interface{} `json:"page_data"`
}

func DefaultPage() *Pagination {
	return &Pagination{
		Page:    1,
		Perpage: 10,
	}
}

func NewPageResult(pagination *Pagination, data interface{}) *PageResult {
	return &PageResult{
		pagination,
		data,
	}

}
