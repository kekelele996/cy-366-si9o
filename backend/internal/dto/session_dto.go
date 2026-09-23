package dto

// StartSessionReq 上机开机请求。
type StartSessionReq struct {
	StationID     uint   `json:"station_id" binding:"required"`
	ReservationID uint   `json:"reservation_id"`
	GameType      string `json:"game_type" binding:"omitempty,oneof=lol csgo kog other"`
}

// RenewSessionReq 续费请求。
type RenewSessionReq struct {
	AddMinutes int `json:"add_minutes" binding:"required,min=1"`
}

// EndSessionReq 下机请求。
type EndSessionReq struct {
	GameType string `json:"game_type" binding:"omitempty,oneof=lol csgo kog other"`
}

// RankQuery 排行榜查询参数。
type RankQuery struct {
	Period   string `form:"period" binding:"omitempty,oneof=day week month"`
	GameType string `form:"game_type" binding:"omitempty,oneof=lol csgo kog other"`
	Limit    int    `form:"limit" binding:"omitempty,min=1,max=100"`
}
