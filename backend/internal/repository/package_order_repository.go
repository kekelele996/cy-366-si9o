package repository

import (
	"errors"

	"gorm.io/gorm"

	"github.com/esportsbar/backend/internal/model"
)

// PackageOrderRepository 时长包订单仓储。
type PackageOrderRepository struct {
	db *gorm.DB
}

// NewPackageOrderRepository 构造时长包订单仓储。
func NewPackageOrderRepository(db *gorm.DB) *PackageOrderRepository {
	return &PackageOrderRepository{db: db}
}

// Create 创建订单。
func (r *PackageOrderRepository) Create(o *model.PackageOrder) error {
	return r.db.Create(o).Error
}

// FindByID 查询订单。
func (r *PackageOrderRepository) FindByID(id uint) (*model.PackageOrder, error) {
	var o model.PackageOrder
	err := r.db.First(&o, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	return &o, err
}

// Update 更新订单。
func (r *PackageOrderRepository) Update(o *model.PackageOrder) error {
	return r.db.Save(o).Error
}

// ListByUser 查询用户订单。
func (r *PackageOrderRepository) ListByUser(userID uint, page, pageSize int) ([]model.PackageOrder, int64, error) {
	var list []model.PackageOrder
	var total int64
	query := r.db.Model(&model.PackageOrder{}).Where("user_id = ?", userID)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := query.Order("id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&list).Error
	return list, total, err
}
