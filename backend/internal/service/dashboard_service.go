package service

import (
	"fmt"
	"log/slog"

	"gorm.io/gorm"

	"github.com/esportsbar/backend/internal/model"
)

// DashboardService 看板统计服务。
type DashboardService struct {
	db     *gorm.DB
	logger *slog.Logger
}

// NewDashboardService 构造看板统计服务。
func NewDashboardService(db *gorm.DB, logger *slog.Logger) *DashboardService {
	return &DashboardService{db: db, logger: logger}
}

// Summary 看板汇总。
type Summary struct {
	StationTotal      int64 `json:"station_total"`
	StationIdle       int64 `json:"station_idle"`
	StationUsing      int64 `json:"station_using"`
	StationFault      int64 `json:"station_fault"`
	StationReserved   int64 `json:"station_reserved"`
	ActiveSession     int64 `json:"active_session"`
	MemberTotal       int64 `json:"member_total"`
	TournamentRunning int64 `json:"tournament_running"`
}

// GetSummary 汇总统计。
func (s *DashboardService) GetSummary() (*Summary, error) {
	summary := &Summary{}
	var station model.Station
	var session model.Session
	var user model.User
	var tournament model.Tournament

	if err := s.db.Model(&station).Count(&summary.StationTotal).Error; err != nil {
		return nil, fmt.Errorf("dashboard count station: %w", err)
	}
	if err := s.db.Model(&station).Where("status = ?", "idle").Count(&summary.StationIdle).Error; err != nil {
		return nil, fmt.Errorf("dashboard count station idle: %w", err)
	}
	if err := s.db.Model(&station).Where("status = ?", "using").Count(&summary.StationUsing).Error; err != nil {
		return nil, fmt.Errorf("dashboard count station using: %w", err)
	}
	if err := s.db.Model(&station).Where("status = ?", "fault").Count(&summary.StationFault).Error; err != nil {
		return nil, fmt.Errorf("dashboard count station fault: %w", err)
	}
	if err := s.db.Model(&station).Where("status = ?", "reserved").Count(&summary.StationReserved).Error; err != nil {
		return nil, fmt.Errorf("dashboard count station reserved: %w", err)
	}
	if err := s.db.Model(&session).Where("status = ?", "active").Count(&summary.ActiveSession).Error; err != nil {
		return nil, fmt.Errorf("dashboard count session active: %w", err)
	}
	if err := s.db.Model(&user).Where("role = ?", "member").Count(&summary.MemberTotal).Error; err != nil {
		return nil, fmt.Errorf("dashboard count member: %w", err)
	}
	if err := s.db.Model(&tournament).Where("status IN ?", []string{"open", "ready"}).Count(&summary.TournamentRunning).Error; err != nil {
		return nil, fmt.Errorf("dashboard count tournament: %w", err)
	}
	s.logger.Info("dashboard summary queried")
	return summary, nil
}
