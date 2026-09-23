package repository

import (
	"errors"

	"gorm.io/gorm"

	"github.com/esportsbar/backend/internal/model"
)

// RegistrationRepository 报名记录仓储。
type RegistrationRepository struct {
	db *gorm.DB
}

// NewRegistrationRepository 构造报名记录仓储。
func NewRegistrationRepository(db *gorm.DB) *RegistrationRepository {
	return &RegistrationRepository{db: db}
}

// Create 创建报名记录。
func (r *RegistrationRepository) Create(reg *model.Registration) error {
	return r.db.Create(reg).Error
}

// FindByID 查询报名记录。
func (r *RegistrationRepository) FindByID(id uint) (*model.Registration, error) {
	var reg model.Registration
	err := r.db.First(&reg, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	return &reg, err
}

// Update 更新报名记录。
func (r *RegistrationRepository) Update(reg *model.Registration) error {
	return r.db.Save(reg).Error
}

// CountByTournamentUser 统计用户在赛事中的报名数（防重复报名）。
func (r *RegistrationRepository) CountByTournamentUser(tournamentID, userID uint) (int64, error) {
	var cnt int64
	err := r.db.Model(&model.Registration{}).
		Where("tournament_id = ? AND user_id = ? AND status IN ?", tournamentID, userID, []string{"pending", "joined"}).
		Count(&cnt).Error
	return cnt, err
}

// ListByTournament 查询赛事报名列表。
func (r *RegistrationRepository) ListByTournament(tournamentID uint) ([]model.Registration, error) {
	var list []model.Registration
	err := r.db.Where("tournament_id = ?", tournamentID).Order("id").Find(&list).Error
	return list, err
}

// AssignGroup 为报名记录分配小组。
func (r *RegistrationRepository) AssignGroup(tx *gorm.DB, id, groupNo int) error {
	return tx.Model(&model.Registration{}).Where("id = ?", id).Update("group_no", groupNo).Error
}
