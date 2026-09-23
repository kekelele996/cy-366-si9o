package model

import "time"

// Reservation 机位预约。
type Reservation struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"index;not null" json:"user_id"`
	StationID uint      `gorm:"index;not null;index:idx_reservations_waitlist,priority:1" json:"station_id"`
	StartTime time.Time `gorm:"not null;index:idx_reservations_waitlist,priority:3" json:"start_time"`
	EndTime   time.Time `gorm:"not null;index:idx_reservations_waitlist,priority:4" json:"end_time"`
	Status    string    `gorm:"size:16;default:pending;index:idx_reservations_waitlist,priority:2" json:"status"` // pending/confirmed/waitlisted/checked_in/completed/cancelled
	Remark    string    `gorm:"size:255" json:"remark"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// QueuePosition 候补顺位（非持久化，仅列表查询时按“同机位、时段重叠、更早提交”的候补数回填，1 表示队首）。
	QueuePosition int `gorm:"-" json:"queue_position"`
}

// TableName 指定表名。
func (Reservation) TableName() string { return "reservations" }
