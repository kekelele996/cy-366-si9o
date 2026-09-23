package repository

import (
	"errors"

	"gorm.io/gorm"

	"github.com/esportsbar/backend/internal/model"
)

// TimePackageRepository 时长包仓储。
type TimePackageRepository struct {
	db *gorm.DB
}

// NewTimePackageRepository 构造时长包仓储。
func NewTimePackageRepository(db *gorm.DB) *TimePackageRepository {
	return &TimePackageRepository{db: db}
}

// Create 创建时长包。
func (r *TimePackageRepository) Create(p *model.TimePackage) error {
	return r.db.Create(p).Error
}

// FindByID 查询时长包。
func (r *TimePackageRepository) FindByID(id uint) (*model.TimePackage, error) {
	var p model.TimePackage
	err := r.db.First(&p, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	return &p, err
}

// Update 更新时长包。
func (r *TimePackageRepository) Update(p *model.TimePackage) error {
	return r.db.Save(p).Error
}

// Delete 删除时长包。
func (r *TimePackageRepository) Delete(id uint) error {
	return r.db.Delete(&model.TimePackage{}, id).Error
}

// List 分页查询时长包。
func (r *TimePackageRepository) List(page, pageSize int) ([]model.TimePackage, int64, error) {
	var list []model.TimePackage
	var total int64
	query := r.db.Model(&model.TimePackage{})
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := query.Order("id").Offset((page - 1) * pageSize).Limit(pageSize).Find(&list).Error
	return list, total, err
}

// ListActive 查询在售时长包（前端购买页复用）。
func (r *TimePackageRepository) ListActive() ([]model.TimePackage, error) {
	var list []model.TimePackage
	err := r.db.Where("status = ?", "active").Order("price").Find(&list).Error
	return list, err
}
