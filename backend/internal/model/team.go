package model

import "time"

// Team 参赛战队。
type Team struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	TournamentID uint      `gorm:"index;not null" json:"tournament_id"`
	Name         string    `gorm:"size:64;not null" json:"name"`
	LeaderID     uint      `gorm:"index" json:"leader_id"`
	MemberCount  int       `gorm:"default:1" json:"member_count"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// TableName 指定表名。
func (Team) TableName() string { return "teams" }
