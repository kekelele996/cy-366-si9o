package repository

import (
	"errors"

	"gorm.io/gorm"

	"github.com/esportsbar/backend/internal/model"
)

// TeamRepository 战队仓储。
type TeamRepository struct {
	db *gorm.DB
}

// NewTeamRepository 构造战队仓储。
func NewTeamRepository(db *gorm.DB) *TeamRepository {
	return &TeamRepository{db: db}
}

// Create 创建战队。
func (r *TeamRepository) Create(t *model.Team) error {
	return r.db.Create(t).Error
}

// FindByID 查询战队。
func (r *TeamRepository) FindByID(id uint) (*model.Team, error) {
	var t model.Team
	err := r.db.First(&t, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	return &t, err
}

// Update 更新战队。
func (r *TeamRepository) Update(t *model.Team) error {
	return r.db.Save(t).Error
}

// ListByTournament 查询赛事下战队。
func (r *TeamRepository) ListByTournament(tournamentID uint) ([]model.Team, error) {
	var list []model.Team
	err := r.db.Where("tournament_id = ?", tournamentID).Order("id").Find(&list).Error
	return list, err
}

// ListByUser 查询用户创建的战队。
func (r *TeamRepository) ListByUser(userID uint) ([]model.Team, error) {
	var list []model.Team
	err := r.db.Where("leader_id = ?", userID).Order("id DESC").Find(&list).Error
	return list, err
}
