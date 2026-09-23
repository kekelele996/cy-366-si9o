package model

import "time"

// TimePackage 时长包（10小时/30小时/月卡）。
type TimePackage struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"size:64;not null" json:"name"`
	Hours     float64   `gorm:"type:decimal(10,2);not null" json:"hours"`
	Price     float64   `gorm:"type:decimal(10,2);not null" json:"price"`
	ValidDays int       `gorm:"default:30" json:"valid_days"`
	Status    string    `gorm:"size:16;default:active" json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName 指定表名。
func (TimePackage) TableName() string { return "time_packages" }
