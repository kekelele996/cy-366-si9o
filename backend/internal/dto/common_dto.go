package dto

// IDReq ID 路径参数。
type IDReq struct {
	ID uint `uri:"id" binding:"required,min=1"`
}

// PageResult 通用分页结果。
type PageResult struct {
	List     any   `json:"list"`
	Total    int64 `json:"total"`
	Page     int   `json:"page"`
	PageSize int   `json:"page_size"`
}
