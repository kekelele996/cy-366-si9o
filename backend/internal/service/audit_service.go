package service

import (
	"fmt"
	"log/slog"

	"github.com/esportsbar/backend/internal/constants"
	"github.com/esportsbar/backend/internal/model"
	"github.com/esportsbar/backend/internal/repository"
)

// AuditService 审计日志服务。
type AuditService struct {
	auditRepo *repository.AuditRepository
	logger    *slog.Logger
}

// NewAuditService 构造审计日志服务。
func NewAuditService(auditRepo *repository.AuditRepository, logger *slog.Logger) *AuditService {
	return &AuditService{auditRepo: auditRepo, logger: logger}
}

// Record 写入审计日志。
func (s *AuditService) Record(entry *model.AuditLog) error {
	if err := s.auditRepo.Create(entry); err != nil {
		s.logger.Error("audit record failed", "err", err)
		return fmt.Errorf("audit record: %w", err)
	}
	s.logger.Info(fmt.Sprintf(constants.LogTemplates["audit_write_ok"], entry.Action, entry.Module))
	return nil
}

// List 分页查询审计日志。
func (s *AuditService) List(page, pageSize int, userID uint, module, action string) ([]model.AuditLog, int64, error) {
	return s.auditRepo.List(page, pageSize, userID, module, action)
}
