package dto

// CreatePackageReq 创建时长包请求。
type CreatePackageReq struct {
	Name      string  `json:"name" binding:"required,max=64"`
	Hours     float64 `json:"hours" binding:"required,gt=0"`
	Price     float64 `json:"price" binding:"required,gt=0"`
	ValidDays int     `json:"valid_days" binding:"omitempty,min=1"`
}

// UpdatePackageReq 更新时长包请求。
type UpdatePackageReq struct {
	Name      string  `json:"name" binding:"omitempty,max=64"`
	Hours     float64 `json:"hours" binding:"omitempty,gt=0"`
	Price     float64 `json:"price" binding:"omitempty,gt=0"`
	ValidDays int     `json:"valid_days" binding:"omitempty,min=1"`
	Status    string  `json:"status" binding:"omitempty,oneof=active disabled"`
}
