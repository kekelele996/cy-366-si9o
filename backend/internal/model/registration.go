package model

import "time"

// Registration 赛事报名记录。
type Registration struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	TournamentID uint      `gorm:"index;not null" json:"tournament_id"`
	UserID       uint      `gorm:"index;not null" json:"user_id"`
	TeamID       uint      `gorm:"index" json:"team_id"`
	Mode         string    `gorm:"size:16;default:solo" json:"mode"` // solo/team
	GroupNo      int       `gorm:"default:0" json:"group_no"`
	Status       string    `gorm:"size:16;default:pending" json:"status"` // pending/joined/rejected
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// TableName 指定表名。
func (Registration) TableName() string { return "registrations" }
