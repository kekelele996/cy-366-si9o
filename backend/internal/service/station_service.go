package service

import (
	"errors"
	"fmt"
	"log/slog"

	"gorm.io/gorm"

	"github.com/esportsbar/backend/internal/constants"
	"github.com/esportsbar/backend/internal/dto"
	"github.com/esportsbar/backend/internal/model"
	"github.com/esportsbar/backend/internal/repository"
	"github.com/esportsbar/backend/internal/util"
)

// StationService 机位服务。
type StationService struct {
	stationRepo *repository.StationRepository
	logger      *slog.Logger
}

// NewStationService 构造机位服务。
func NewStationService(stationRepo *repository.StationRepository, logger *slog.Logger) *StationService {
	return &StationService{stationRepo: stationRepo, logger: logger}
}

// Create 创建机位。
func (s *StationService) Create(req *dto.CreateStationReq) (*model.Station, error) {
	station := &model.Station{
		Name:         req.Name,
		Area:         req.Area,
		StationType:  req.StationType,
		PricePerHour: req.PricePerHour,
		Description:  req.Description,
		Status:       constants.StationIdle,
	}
	if err := s.stationRepo.Create(station); err != nil {
		return nil, fmt.Errorf("station create: %w", err)
	}
	s.logger.Info(fmt.Sprintf(constants.LogTemplates["station_create_ok"], station.Name, station.Area))
	return station, nil
}

// Update 更新机位。
func (s *StationService) Update(id uint, req *dto.UpdateStationReq) (*model.Station, error) {
	station, err := s.stationRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, util.NewAppError(constants.CodeNotFound, "机位不存在")
		}
		return nil, fmt.Errorf("station update find: %w", err)
	}
	if req.Name != "" {
		station.Name = req.Name
	}
	if req.Area != "" {
		station.Area = req.Area
	}
	if req.StationType != "" {
		station.StationType = req.StationType
	}
	if req.PricePerHour > 0 {
		station.PricePerHour = req.PricePerHour
	}
	if req.Description != "" {
		station.Description = req.Description
	}
	if err := s.stationRepo.Update(station); err != nil {
		return nil, fmt.Errorf("station update: %w", err)
	}
	s.logger.Info(fmt.Sprintf(constants.LogTemplates["station_update_ok"], station.ID, "profile"))
	return station, nil
}

// UpdateStatus 机位状态流转：idle<->using/reserved/fault。
func (s *StationService) UpdateStatus(id uint, req *dto.UpdateStationStatusReq) (*model.Station, error) {
	station, err := s.stationRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, util.NewAppError(constants.CodeNotFound, "机位不存在")
		}
		return nil, fmt.Errorf("station status find: %w", err)
	}
	from := station.Status
	if !allowedStationTransition(from, req.Status) {
		return nil, util.NewAppError(constants.CodeConflict,
			fmt.Sprintf("机位状态不允许从 %s 变更为 %s，请先处理当前状态", util.StatusText(from), util.StatusText(req.Status)))
	}
	station.Status = req.Status
	if err := s.stationRepo.Update(station); err != nil {
		return nil, fmt.Errorf("station status update: %w", err)
	}
	s.logger.Info(fmt.Sprintf(constants.LogTemplates["station_status_change"], station.ID, from, req.Status, "operator"))
	return station, nil
}

// Delete 删除机位。
func (s *StationService) Delete(id uint) error {
	station, err := s.stationRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil
		}
		return fmt.Errorf("station delete find: %w", err)
	}
	if station.Status == constants.StationUsing {
		return util.NewAppError(constants.CodeStationBusy, "使用中的机位不能删除，请先下机")
	}
	if err := s.stationRepo.Delete(id); err != nil {
		return fmt.Errorf("station delete: %w", err)
	}
	s.logger.Info(fmt.Sprintf(constants.LogTemplates["station_delete_ok"], id))
	return nil
}

// List 分页查询机位。
func (s *StationService) List(query *dto.StationQuery) ([]model.Station, int64, error) {
	page := query.Page
	pageSize := query.PageSize
	if page <= 0 {
		page = constants.DefaultPage
	}
	if pageSize <= 0 {
		pageSize = constants.DefaultPageSize
	}
	return s.stationRepo.List(page, pageSize, query.Area, query.Status)
}

// ListAll 查询全部机位（看板/下拉复用）。
func (s *StationService) ListAll() ([]model.Station, error) {
	return s.stationRepo.ListAll()
}

// GetByID 查询机位详情。
func (s *StationService) GetByID(id uint) (*model.Station, error) {
	station, err := s.stationRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, util.NewAppError(constants.CodeNotFound, "机位不存在")
		}
		return nil, fmt.Errorf("station get: %w", err)
	}
	return station, nil
}

// LockForUpdate 事务内行锁机位，供预约/上机服务复用。
func (s *StationService) LockForUpdate(tx *gorm.DB, id uint) (*model.Station, error) {
	return s.stationRepo.LockByID(tx, id)
}

// allowedStationTransition 机位状态机：using 只能由 idle 进入，且必须通过上机/下机服务流转。
func allowedStationTransition(from, to string) bool {
	if from == to {
		return true
	}
	switch from {
	case constants.StationIdle:
		return to == constants.StationFault || to == constants.StationReserved || to == constants.StationUsing
	case constants.StationFault:
		return to == constants.StationIdle
	case constants.StationReserved:
		return to == constants.StationIdle || to == constants.StationUsing
	case constants.StationUsing:
		return to == constants.StationIdle
	}
	return false
}
