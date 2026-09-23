package model

import "time"

// PackageOrder 时长包购买订单。
type PackageOrder struct {
	ID            uint       `gorm:"primaryKey" json:"id"`
	OrderNo       string     `gorm:"size:64;uniqueIndex;not null" json:"order_no"`
	UserID        uint       `gorm:"index;not null" json:"user_id"`
	PackageID     uint       `gorm:"index;not null" json:"package_id"`
	PackageName   string     `gorm:"size:64" json:"package_name"`
	Amount        float64    `gorm:"type:decimal(12,2);not null" json:"amount"`
	Hours         float64    `gorm:"type:decimal(10,2)" json:"hours"`
	PaymentMethod string     `gorm:"size:16;default:balance" json:"payment_method"`
	Status        string     `gorm:"size:16;default:pending" json:"status"` // pending/paid/cancelled
	PaidAt        *time.Time `json:"paid_at"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

// TableName 指定表名。
func (PackageOrder) TableName() string { return "package_orders" }
