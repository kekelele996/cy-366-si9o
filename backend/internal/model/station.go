package model

import "time"

// Station 机位/包厢。
type Station struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Name         string    `gorm:"size:64;not null" json:"name"`
	Area         string    `gorm:"size:64;not null" json:"area"`
	StationType  string    `gorm:"size:16;default:seat" json:"station_type"` // seat 机位 / box 包厢
	PricePerHour float64   `gorm:"type:decimal(10,2);default:0" json:"price_per_hour"`
	Status       string    `gorm:"size:16;default:idle" json:"status"` // idle/using/fault/reserved
	Description  string    `gorm:"size:255" json:"description"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// TableName 指定表名。
func (Station) TableName() string { return "stations" }
