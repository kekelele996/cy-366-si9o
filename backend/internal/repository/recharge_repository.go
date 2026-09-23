package repository

import (
	"gorm.io/gorm"

	"github.com/esportsbar/backend/internal/model"
)

// RechargeRepository 充值记录仓储。
type RechargeRepository struct {
	db *gorm.DB
}

// NewRechargeRepository 构造充值记录仓储。
func NewRechargeRepository(db *gorm.DB) *RechargeRepository {
	return &RechargeRepository{db: db}
}

// Create 创建充值记录。
func (r *RechargeRepository) Create(rc *model.Recharge) error {
	return r.db.Create(rc).Error
}

// ListByUser 查询用户充值记录。
func (r *RechargeRepository) ListByUser(userID uint, page, pageSize int) ([]model.Recharge, int64, error) {
	var list []model.Recharge
	var total int64
	query := r.db.Model(&model.Recharge{}).Where("user_id = ?", userID)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := query.Order("id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&list).Error
	return list, total, err
}
