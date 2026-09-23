package repository

import (
	"errors"

	"gorm.io/gorm"

	"github.com/esportsbar/backend/internal/model"
)

// TournamentRepository 赛事仓储。
type TournamentRepository struct {
	db *gorm.DB
}

// NewTournamentRepository 构造赛事仓储。
func NewTournamentRepository(db *gorm.DB) *TournamentRepository {
	return &TournamentRepository{db: db}
}

// Create 创建赛事。
func (r *TournamentRepository) Create(t *model.Tournament) error {
	return r.db.Create(t).Error
}

// FindByID 查询赛事。
func (r *TournamentRepository) FindByID(id uint) (*model.Tournament, error) {
	var t model.Tournament
	err := r.db.First(&t, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	return &t, err
}

// Update 更新赛事。
func (r *TournamentRepository) Update(t *model.Tournament) error {
	return r.db.Save(t).Error
}

// Delete 删除赛事。
func (r *TournamentRepository) Delete(id uint) error {
	return r.db.Delete(&model.Tournament{}, id).Error
}

// List 分页查询赛事。
func (r *TournamentRepository) List(page, pageSize int, status string) ([]model.Tournament, int64, error) {
	var list []model.Tournament
	var total int64
	query := r.db.Model(&model.Tournament{})
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := query.Order("id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&list).Error
	return list, total, err
}
