package dto

// CreateStationReq 创建机位请求。
type CreateStationReq struct {
	Name         string  `json:"name" binding:"required,max=64"`
	Area         string  `json:"area" binding:"required,max=64"`
	StationType  string  `json:"station_type" binding:"required,oneof=seat box"`
	PricePerHour float64 `json:"price_per_hour" binding:"omitempty,min=0"`
	Description  string  `json:"description" binding:"omitempty,max=255"`
}

// UpdateStationReq 更新机位请求。
type UpdateStationReq struct {
	Name         string  `json:"name" binding:"omitempty,max=64"`
	Area         string  `json:"area" binding:"omitempty,max=64"`
	StationType  string  `json:"station_type" binding:"omitempty,oneof=seat box"`
	PricePerHour float64 `json:"price_per_hour" binding:"omitempty,min=0"`
	Description  string  `json:"description" binding:"omitempty,max=255"`
}

// UpdateStationStatusReq 机位状态变更请求。
type UpdateStationStatusReq struct {
	Status string `json:"status" binding:"required,oneof=idle using fault reserved"`
}

// StationQuery 机位查询参数。
type StationQuery struct {
	Page     int    `form:"page" binding:"omitempty,min=1"`
	PageSize int    `form:"page_size" binding:"omitempty,min=1,max=100"`
	Area     string `form:"area"`
	Status   string `form:"status" binding:"omitempty,oneof=idle using fault reserved"`
}
