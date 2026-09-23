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

// AuthService 认证服务。
type AuthService struct {
	userRepo  *repository.UserRepository
	logger    *slog.Logger
	jwtSecret string
	expireSec int
}

// NewAuthService 构造认证服务。
func NewAuthService(userRepo *repository.UserRepository, logger *slog.Logger, jwtSecret string, expireSec int) *AuthService {
	return &AuthService{userRepo: userRepo, logger: logger, jwtSecret: jwtSecret, expireSec: expireSec}
}

// Register 会员注册。
func (s *AuthService) Register(req *dto.RegisterReq) (*dto.LoginResp, error) {
	if _, err := s.userRepo.FindByUsername(req.Username); err == nil {
		return nil, util.NewAppError(constants.CodeUserExists, "用户名已存在，请更换后重试")
	} else if !errors.Is(err, repository.ErrNotFound) {
		return nil, fmt.Errorf("auth register check user: %w", err)
	}
	hash, err := util.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("auth register hash password: %w", err)
	}
	user := &model.User{
		Username: req.Username,
		Password: hash,
		Nickname: req.Nickname,
		Phone:    req.Phone,
		Role:     constants.RoleMember,
		Status:   "active",
	}
	if err := s.userRepo.Create(user); err != nil {
		return nil, fmt.Errorf("auth register create user: %w", err)
	}
	s.logger.Info(fmt.Sprintf(constants.LogTemplates["user_register_ok"], user.Username, user.Role))
	token, err := util.GenerateToken(s.jwtSecret, s.expireSec, user.ID, user.Username, user.Role)
	if err != nil {
		return nil, fmt.Errorf("auth register generate token: %w", err)
	}
	return &dto.LoginResp{Token: token, User: toUserItem(user)}, nil
}

// Login 登录。
func (s *AuthService) Login(req *dto.LoginReq) (*dto.LoginResp, error) {
	user, err := s.userRepo.FindByUsername(req.Username)
	if err != nil {
		s.logger.Warn(fmt.Sprintf(constants.LogTemplates["user_login_fail"], req.Username, "user not found"))
		return nil, util.NewAppError(constants.CodeWrongPass, "用户名或密码错误")
	}
	if !util.CheckPassword(user.Password, req.Password) {
		s.logger.Warn(fmt.Sprintf(constants.LogTemplates["user_login_fail"], req.Username, "wrong password"))
		return nil, util.NewAppError(constants.CodeWrongPass, "用户名或密码错误")
	}
	if user.Status == "disabled" {
		return nil, util.NewAppError(constants.CodeForbidden, "账号已被禁用，请联系管理员")
	}
	s.logger.Info(fmt.Sprintf(constants.LogTemplates["user_login_ok"], user.Username, user.Role))
	token, err := util.GenerateToken(s.jwtSecret, s.expireSec, user.ID, user.Username, user.Role)
	if err != nil {
		return nil, fmt.Errorf("auth login generate token: %w", err)
	}
	return &dto.LoginResp{Token: token, User: toUserItem(user)}, nil
}

// toUserItem 模型转响应对象。
func toUserItem(u *model.User) *dto.UserItem {
	return &dto.UserItem{
		ID:       u.ID,
		Username: u.Username,
		Nickname: u.Nickname,
		Phone:    u.Phone,
		Role:     u.Role,
		Balance:  u.Balance,
		Status:   u.Status,
	}
}
