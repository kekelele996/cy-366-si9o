package model

import "time"

// Reservation 机位预约。
type Reservation struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"index;not null" json:"user_id"`
	StationID uint      `gorm:"index;not null" json:"station_id"`
	StartTime time.Time `gorm:"not null" json:"start_time"`
	EndTime   time.Time `gorm:"not null" json:"end_time"`
	Status    string    `gorm:"size:16;default:pending" json:"status"` // pending/waitlisted/confirmed/checked_in/completed/cancelled
	Remark    string    `gorm:"size:255" json:"remark"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// WaitlistPosition 候补顺位（不落库），仅 waitlisted 状态由 service 按提交顺序动态计算填充。
	WaitlistPosition int `gorm:"-" json:"waitlist_position,omitempty"`
}

// TableName 指定表名。
func (Reservation) TableName() string { return "reservations" }
