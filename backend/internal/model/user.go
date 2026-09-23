package model

import "time"

// User 会员/管理员/店员账号。
type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Username  string    `gorm:"size:64;uniqueIndex;not null" json:"username"`
	Password  string    `gorm:"size:128;not null" json:"-"`
	Nickname  string    `gorm:"size:64" json:"nickname"`
	Phone     string    `gorm:"size:20" json:"phone"`
	Role      string    `gorm:"size:16;default:member" json:"role"`
	Balance   float64   `gorm:"type:decimal(12,2);default:0" json:"balance"`
	Status    string    `gorm:"size:16;default:active" json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName 指定表名。
func (User) TableName() string { return "users" }
