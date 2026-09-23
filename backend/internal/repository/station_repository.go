package repository

import (
	"errors"

	"gorm.io/gorm"

	"github.com/esportsbar/backend/internal/model"
)

// StationRepository 机位仓储。
type StationRepository struct {
	db *gorm.DB
}

// NewStationRepository 构造机位仓储。
func NewStationRepository(db *gorm.DB) *StationRepository {
	return &StationRepository{db: db}
}

// Create 创建机位。
func (r *StationRepository) Create(s *model.Station) error {
	return r.db.Create(s).Error
}

// FindByID 查询机位。
func (r *StationRepository) FindByID(id uint) (*model.Station, error) {
	var s model.Station
	err := r.db.First(&s, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	return &s, err
}

// Update 更新机位。
func (r *StationRepository) Update(s *model.Station) error {
	return r.db.Save(s).Error
}

// Delete 删除机位。
func (r *StationRepository) Delete(id uint) error {
	return r.db.Delete(&model.Station{}, id).Error
}

// List 分页查询机位，支持区域与状态筛选。
func (r *StationRepository) List(page, pageSize int, area, status string) ([]model.Station, int64, error) {
	var list []model.Station
	var total int64
	query := r.db.Model(&model.Station{})
	if area != "" {
		query = query.Where("area = ?", area)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := query.Order("area, id").Offset((page - 1) * pageSize).Limit(pageSize).Find(&list).Error
	return list, total, err
}

// ListAll 查询全部机位（看板/下拉选择复用）。
func (r *StationRepository) ListAll() ([]model.Station, error) {
	var list []model.Station
	err := r.db.Order("area, id").Find(&list).Error
	return list, err
}

// LockByID 行锁查询机位（并发状态流转使用）。
func (r *StationRepository) LockByID(tx *gorm.DB, id uint) (*model.Station, error) {
	var s model.Station
	err := tx.Clauses(clauseLocking()).First(&s, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	return &s, err
}
