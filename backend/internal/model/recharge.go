package model

import "time"

// Recharge 会员充值记录。
type Recharge struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	UserID        uint      `gorm:"index;not null" json:"user_id"`
	Amount        float64   `gorm:"type:decimal(12,2);not null" json:"amount"`
	PaymentMethod string    `gorm:"size:16;default:balance" json:"payment_method"`
	OperatorID    uint      `gorm:"index" json:"operator_id"`
	Remark        string    `gorm:"size:255" json:"remark"`
	CreatedAt     time.Time `json:"created_at"`
}

// TableName 指定表名。
func (Recharge) TableName() string { return "recharges" }
