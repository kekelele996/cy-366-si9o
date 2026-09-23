package model

import "time"

// Match 比赛场次。
type Match struct {
	ID           uint       `gorm:"primaryKey" json:"id"`
	TournamentID uint       `gorm:"index;not null" json:"tournament_id"`
	Round        int        `gorm:"default:1" json:"round"`
	GroupNo      int        `gorm:"default:0" json:"group_no"`
	TeamAID      uint       `gorm:"index" json:"team_a_id"`
	TeamBID      uint       `gorm:"index" json:"team_b_id"`
	ScoreA       int        `gorm:"default:0" json:"score_a"`
	ScoreB       int        `gorm:"default:0" json:"score_b"`
	WinnerID     uint       `json:"winner_id"`
	Status       string     `gorm:"size:16;default:pending" json:"status"` // pending/playing/finished
	ScheduledAt  *time.Time `json:"scheduled_at"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

// TableName 指定表名。
func (Match) TableName() string { return "matches" }
