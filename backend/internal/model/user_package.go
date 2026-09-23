package model

import "time"

// UserPackage 会员已购时长包。
type UserPackage struct {
	ID             uint       `gorm:"primaryKey" json:"id"`
	UserID         uint       `gorm:"index;not null" json:"user_id"`
	PackageID      uint       `gorm:"index;not null" json:"package_id"`
	PackageName    string     `gorm:"size:64" json:"package_name"`
	TotalHours     float64    `gorm:"type:decimal(10,2);default:0" json:"total_hours"`
	RemainingHours float64    `gorm:"type:decimal(10,2);default:0" json:"remaining_hours"`
	ExpireAt       *time.Time `json:"expire_at"`
	Status         string     `gorm:"size:16;default:active" json:"status"` // active/expired/used
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

// TableName 指定表名。
func (UserPackage) TableName() string { return "user_packages" }
