package repository

import (
	"errors"
	"time"

	"gorm.io/gorm"

	"github.com/esportsbar/backend/internal/model"
)

// SessionRepository 上机记录仓储。
type SessionRepository struct {
	db *gorm.DB
}

// NewSessionRepository 构造上机记录仓储。
func NewSessionRepository(db *gorm.DB) *SessionRepository {
	return &SessionRepository{db: db}
}

// Create 创建上机记录。
func (r *SessionRepository) Create(s *model.Session) error {
	return r.db.Create(s).Error
}

// FindByID 查询上机记录。
func (r *SessionRepository) FindByID(id uint) (*model.Session, error) {
	var s model.Session
	err := r.db.First(&s, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	return &s, err
}

// FindActiveByStation 查询机位进行中的上机记录。
func (r *SessionRepository) FindActiveByStation(stationID uint) (*model.Session, error) {
	var s model.Session
	err := r.db.Where("station_id = ? AND status = ?", stationID, "active").First(&s).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	return &s, err
}

// Update 更新上机记录。
func (r *SessionRepository) Update(s *model.Session) error {
	return r.db.Save(s).Error
}

// List 分页查询上机记录。
func (r *SessionRepository) List(page, pageSize int, userID uint, status string) ([]model.Session, int64, error) {
	var list []model.Session
	var total int64
	query := r.db.Model(&model.Session{})
	if userID > 0 {
		query = query.Where("user_id = ?", userID)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := query.Order("id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&list).Error
	return list, total, err
}

// Rank 统计指定周期内会员累计上机时长排行榜。
func (r *SessionRepository) Rank(period string, gameType string, limit int) ([]model.Session, int64, error) {
	now := time.Now()
	var since time.Time
	switch period {
	case "day":
		since = now.AddDate(0, 0, -1)
	case "week":
		since = now.AddDate(0, 0, -7)
	default:
		since = now.AddDate(0, -1, 0)
	}
	query := r.db.Model(&model.Session{}).
		Select("user_id, SUM(duration_minutes) AS duration_minutes, COUNT(*) AS session_count").
		Where("status = ? AND start_time >= ?", "completed", since)
	if gameType != "" {
		query = query.Where("game_type = ?", gameType)
	}
	var rows []model.Session
	err := query.Group("user_id").Order("duration_minutes DESC").Limit(limit).Scan(&rows).Error
	return rows, 0, err
}
