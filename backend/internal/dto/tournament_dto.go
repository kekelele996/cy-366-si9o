package dto

import "time"

// CreateTournamentReq 创建赛事请求。
type CreateTournamentReq struct {
	Name          string     `json:"name" binding:"required,max=128"`
	GameType      string     `json:"game_type" binding:"required,oneof=lol csgo kog other"`
	Description   string     `json:"description" binding:"omitempty,max=512"`
	RegisterStart *time.Time `json:"register_start"`
	RegisterEnd   *time.Time `json:"register_end"`
	StartTime     *time.Time `json:"start_time"`
	MaxTeams      int        `json:"max_teams" binding:"omitempty,min=2,max=64"`
}

// UpdateTournamentReq 更新赛事请求。
type UpdateTournamentReq struct {
	Name          string     `json:"name" binding:"omitempty,max=128"`
	Description   string     `json:"description" binding:"omitempty,max=512"`
	RegisterStart *time.Time `json:"register_start"`
	RegisterEnd   *time.Time `json:"register_end"`
	StartTime     *time.Time `json:"start_time"`
	MaxTeams      int        `json:"max_teams" binding:"omitempty,min=2,max=64"`
	Status        string     `json:"status" binding:"omitempty,oneof=draft open ready finished"`
}

// CreateTeamReq 创建战队请求。
type CreateTeamReq struct {
	Name        string `json:"name" binding:"required,max=64"`
	MemberCount int    `json:"member_count" binding:"omitempty,min=1,max=10"`
}

// RegisterReq 赛事报名请求。
type TournamentRegisterReq struct {
	Mode   string `json:"mode" binding:"required,oneof=solo team"`
	TeamID uint   `json:"team_id"`
}

// SubmitMatchReq 提交比赛结果请求。
type SubmitMatchReq struct {
	ScoreA   int  `json:"score_a" binding:"omitempty,min=0"`
	ScoreB   int  `json:"score_b" binding:"omitempty,min=0"`
	WinnerID uint `json:"winner_id"`
}
