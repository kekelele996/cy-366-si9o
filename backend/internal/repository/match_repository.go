package repository

import (
	"errors"

	"gorm.io/gorm"

	"github.com/esportsbar/backend/internal/model"
)

// MatchRepository 比赛场次仓储。
type MatchRepository struct {
	db *gorm.DB
}

// NewMatchRepository 构造比赛场次仓储。
func NewMatchRepository(db *gorm.DB) *MatchRepository {
	return &MatchRepository{db: db}
}

// Create 创建比赛场次。
func (r *MatchRepository) Create(m *model.Match) error {
	return r.db.Create(m).Error
}

// FindByID 查询比赛场次。
func (r *MatchRepository) FindByID(id uint) (*model.Match, error) {
	var m model.Match
	err := r.db.First(&m, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	return &m, err
}

// Update 更新比赛场次。
func (r *MatchRepository) Update(m *model.Match) error {
	return r.db.Save(m).Error
}

// ListByTournament 查询赛事比赛场次。
func (r *MatchRepository) ListByTournament(tournamentID uint) ([]model.Match, error) {
	var list []model.Match
	err := r.db.Where("tournament_id = ?", tournamentID).Order("round, group_no, id").Find(&list).Error
	return list, err
}
