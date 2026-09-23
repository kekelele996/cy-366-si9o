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

// TimePackageService 时长包管理服务。
type TimePackageService struct {
	packageRepo *repository.TimePackageRepository
	logger      *slog.Logger
}

// NewTimePackageService 构造时长包管理服务。
func NewTimePackageService(packageRepo *repository.TimePackageRepository, logger *slog.Logger) *TimePackageService {
	return &TimePackageService{packageRepo: packageRepo, logger: logger}
}

// Create 创建时长包。
func (s *TimePackageService) Create(req *dto.CreatePackageReq) (*model.TimePackage, error) {
	pkg := &model.TimePackage{
		Name:      req.Name,
		Hours:     req.Hours,
		Price:     req.Price,
		ValidDays: req.ValidDays,
		Status:    "active",
	}
	if err := s.packageRepo.Create(pkg); err != nil {
		return nil, fmt.Errorf("time package create: %w", err)
	}
	return pkg, nil
}

// Update 更新时长包。
func (s *TimePackageService) Update(id uint, req *dto.UpdatePackageReq) (*model.TimePackage, error) {
	pkg, err := s.packageRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, util.NewAppError(constants.CodeNotFound, "时长包不存在")
		}
		return nil, fmt.Errorf("time package update find: %w", err)
	}
	if req.Name != "" {
		pkg.Name = req.Name
	}
	if req.Hours > 0 {
		pkg.Hours = req.Hours
	}
	if req.Price > 0 {
		pkg.Price = req.Price
	}
	if req.ValidDays > 0 {
		pkg.ValidDays = req.ValidDays
	}
	if req.Status != "" {
		pkg.Status = req.Status
	}
	if err := s.packageRepo.Update(pkg); err != nil {
		return nil, fmt.Errorf("time package update: %w", err)
	}
	return pkg, nil
}

// Delete 删除时长包。
func (s *TimePackageService) Delete(id uint) error {
	if err := s.packageRepo.Delete(id); err != nil {
		return fmt.Errorf("time package delete: %w", err)
	}
	return nil
}

// List 分页查询时长包。
func (s *TimePackageService) List(page, pageSize int) ([]model.TimePackage, int64, error) {
	return s.packageRepo.List(page, pageSize)
}

// ListActive 查询在售时长包。
func (s *TimePackageService) ListActive() ([]model.TimePackage, error) {
	return s.packageRepo.ListActive()
}

// GetByID 查询时长包详情。
func (s *TimePackageService) GetByID(id uint) (*model.TimePackage, error) {
	pkg, err := s.packageRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, util.NewAppError(constants.CodeNotFound, "时长包不存在")
		}
		return nil, fmt.Errorf("time package get: %w", err)
	}
	return pkg, nil
}
