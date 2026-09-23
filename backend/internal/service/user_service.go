package service

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/esportsbar/backend/internal/constants"
	"github.com/esportsbar/backend/internal/dto"
	"github.com/esportsbar/backend/internal/model"
	"github.com/esportsbar/backend/internal/repository"
	"github.com/esportsbar/backend/internal/util"
)

// UserService 用户管理服务。
type UserService struct {
	userRepo *repository.UserRepository
	logger   *slog.Logger
}

// NewUserService 构造用户管理服务。
func NewUserService(userRepo *repository.UserRepository, logger *slog.Logger) *UserService {
	return &UserService{userRepo: userRepo, logger: logger}
}

// Create 管理员创建用户。
func (s *UserService) Create(req *dto.CreateUserReq) (*dto.UserItem, error) {
	if _, err := s.userRepo.FindByUsername(req.Username); err == nil {
		return nil, util.NewAppError(constants.CodeUserExists, "用户名已存在，请更换后重试")
	} else if !errors.Is(err, repository.ErrNotFound) {
		return nil, fmt.Errorf("user create check: %w", err)
	}
	hash, err := util.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("user create hash: %w", err)
	}
	user := &model.User{
		Username: req.Username,
		Password: hash,
		Nickname: req.Nickname,
		Phone:    req.Phone,
		Role:     req.Role,
		Status:   "active",
	}
	if err := s.userRepo.Create(user); err != nil {
		return nil, fmt.Errorf("user create: %w", err)
	}
	s.logger.Info(fmt.Sprintf(constants.LogTemplates["user_create_ok"], user.Username, user.Role))
	return toUserItem(user), nil
}

// Update 更新用户。
func (s *UserService) Update(id uint, req *dto.UpdateUserReq) (*dto.UserItem, error) {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, util.NewAppError(constants.CodeUserNotFound, "用户不存在")
		}
		return nil, fmt.Errorf("user update find: %w", err)
	}
	if req.Nickname != "" {
		user.Nickname = req.Nickname
	}
	if req.Phone != "" {
		user.Phone = req.Phone
	}
	if req.Role != "" {
		user.Role = req.Role
	}
	if req.Status != "" {
		user.Status = req.Status
	}
	if err := s.userRepo.Update(user); err != nil {
		return nil, fmt.Errorf("user update: %w", err)
	}
	s.logger.Info(fmt.Sprintf(constants.LogTemplates["user_update_ok"], user.ID, "profile"))
	return toUserItem(user), nil
}

// Delete 删除用户。
func (s *UserService) Delete(id uint) error {
	if err := s.userRepo.Delete(id); err != nil {
		return fmt.Errorf("user delete: %w", err)
	}
	s.logger.Info(fmt.Sprintf(constants.LogTemplates["user_delete_ok"], id))
	return nil
}

// List 分页查询用户。
func (s *UserService) List(page, pageSize int) ([]model.User, int64, error) {
	return s.userRepo.List(page, pageSize)
}

// GetByID 查询用户详情。
func (s *UserService) GetByID(id uint) (*dto.UserItem, error) {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, util.NewAppError(constants.CodeUserNotFound, "用户不存在")
		}
		return nil, fmt.Errorf("user get: %w", err)
	}
	return toUserItem(user), nil
}
