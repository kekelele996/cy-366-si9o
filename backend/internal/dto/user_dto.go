package dto

// RegisterReq 注册请求。
type RegisterReq struct {
	Username string `json:"username" binding:"required,min=3,max=32"`
	Password string `json:"password" binding:"required,min=6,max=64"`
	Nickname string `json:"nickname" binding:"required,max=64"`
	Phone    string `json:"phone" binding:"omitempty,max=20"`
}

// LoginReq 登录请求。
type LoginReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResp 登录响应。
type LoginResp struct {
	Token string    `json:"token"`
	User  *UserItem `json:"user"`
}

// UserItem 用户简要信息。
type UserItem struct {
	ID       uint    `json:"id"`
	Username string  `json:"username"`
	Nickname string  `json:"nickname"`
	Phone    string  `json:"phone"`
	Role     string  `json:"role"`
	Balance  float64 `json:"balance"`
	Status   string  `json:"status"`
}

// CreateUserReq 管理员创建用户请求。
type CreateUserReq struct {
	Username string `json:"username" binding:"required,min=3,max=32"`
	Password string `json:"password" binding:"required,min=6,max=64"`
	Nickname string `json:"nickname" binding:"required,max=64"`
	Phone    string `json:"phone" binding:"omitempty,max=20"`
	Role     string `json:"role" binding:"required,oneof=admin staff member"`
}

// UpdateUserReq 更新用户请求。
type UpdateUserReq struct {
	Nickname string `json:"nickname" binding:"omitempty,max=64"`
	Phone    string `json:"phone" binding:"omitempty,max=20"`
	Role     string `json:"role" binding:"omitempty,oneof=admin staff member"`
	Status   string `json:"status" binding:"omitempty,oneof=active disabled"`
}

// PageQuery 分页查询。
type PageQuery struct {
	Page     int `form:"page" binding:"omitempty,min=1"`
	PageSize int `form:"page_size" binding:"omitempty,min=1,max=100"`
}
