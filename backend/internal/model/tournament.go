package model

import "time"

// Tournament 电竞赛事。
type Tournament struct {
	ID            uint       `gorm:"primaryKey" json:"id"`
	Name          string     `gorm:"size:128;not null" json:"name"`
	GameType      string     `gorm:"size:16;default:lol" json:"game_type"`
	Description   string     `gorm:"size:512" json:"description"`
	RegisterStart *time.Time `json:"register_start"`
	RegisterEnd   *time.Time `json:"register_end"`
	StartTime     *time.Time `json:"start_time"`
	MaxTeams      int        `gorm:"default:16" json:"max_teams"`
	Status        string     `gorm:"size:16;default:draft" json:"status"` // draft/open/ready/finished
	CreatedBy     uint       `json:"created_by"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

// TableName 指定表名。
func (Tournament) TableName() string { return "tournaments" }
