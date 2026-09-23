package dto

// AuditQuery 审计日志查询参数。
type AuditQuery struct {
	Page     int    `form:"page" binding:"omitempty,min=1"`
	PageSize int    `form:"page_size" binding:"omitempty,min=1,max=100"`
	UserID   uint   `form:"user_id"`
	Module   string `form:"module"`
	Action   string `form:"action"`
}
