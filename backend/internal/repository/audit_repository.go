package repository

import (
	"gorm.io/gorm"

	"github.com/esportsbar/backend/internal/model"
)

// AuditRepository 审计日志仓储。
type AuditRepository struct {
	db *gorm.DB
}

// NewAuditRepository 构造审计日志仓储。
func NewAuditRepository(db *gorm.DB) *AuditRepository {
	return &AuditRepository{db: db}
}

// Create 写入审计日志。
func (r *AuditRepository) Create(log *model.AuditLog) error {
	return r.db.Create(log).Error
}

// List 分页查询审计日志。
func (r *AuditRepository) List(page, pageSize int, userID uint, module, action string) ([]model.AuditLog, int64, error) {
	var list []model.AuditLog
	var total int64
	query := r.db.Model(&model.AuditLog{})
	if userID > 0 {
		query = query.Where("user_id = ?", userID)
	}
	if module != "" {
		query = query.Where("module = ?", module)
	}
	if action != "" {
		query = query.Where("action = ?", action)
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := query.Order("id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&list).Error
	return list, total, err
}
