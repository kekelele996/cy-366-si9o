package repository

import (
	"errors"

	"gorm.io/gorm"

	"github.com/esportsbar/backend/internal/model"
)

// 仓储层哨兵错误。
var (
	ErrNotFound = errors.New("record not found")
	ErrConflict = errors.New("record conflict")
)

// UserRepository 用户仓储。
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository 构造用户仓储。
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create 创建用户。
func (r *UserRepository) Create(u *model.User) error {
	return r.db.Create(u).Error
}

// FindByUsername 按用户名查询用户。
func (r *UserRepository) FindByUsername(username string) (*model.User, error) {
	var u model.User
	err := r.db.Where("username = ?", username).First(&u).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	return &u, err
}

// FindByID 按 ID 查询用户。
func (r *UserRepository) FindByID(id uint) (*model.User, error) {
	var u model.User
	err := r.db.First(&u, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	return &u, err
}

// Update 更新用户。
func (r *UserRepository) Update(u *model.User) error {
	return r.db.Save(u).Error
}

// Delete 删除用户。
func (r *UserRepository) Delete(id uint) error {
	return r.db.Delete(&model.User{}, id).Error
}

// List 分页查询用户。
func (r *UserRepository) List(page, pageSize int) ([]model.User, int64, error) {
	var users []model.User
	var total int64
	query := r.db.Model(&model.User{})
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := query.Order("id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&users).Error
	return users, total, err
}

// UpdateBalance 更新余额（扣款时校验余额充足）。
func (r *UserRepository) UpdateBalance(userID uint, delta float64) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var u model.User
		if err := tx.Clauses(clauseLocking()).First(&u, userID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrNotFound
			}
			return err
		}
		if u.Balance+delta < 0 {
			return ErrConflict
		}
		return tx.Model(&model.User{}).Where("id = ?", userID).Update("balance", u.Balance+delta).Error
	})
}
