package handler

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/gin-gonic/gin"

	"github.com/esportsbar/backend/internal/constants"
	"github.com/esportsbar/backend/internal/dto"
	"github.com/esportsbar/backend/internal/service"
	"github.com/esportsbar/backend/internal/util"
	"github.com/esportsbar/backend/pkg/response"
)

// AuditHandler 审计日志接口处理器。
type AuditHandler struct {
	auditService *service.AuditService
	logger       *slog.Logger
}

// NewAuditHandler 构造审计日志接口处理器。
func NewAuditHandler(auditService *service.AuditService, logger *slog.Logger) *AuditHandler {
	return &AuditHandler{auditService: auditService, logger: logger}
}

// List 分页查询审计日志。
func (h *AuditHandler) List(c *gin.Context) {
	var query dto.AuditQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Fail(c, 400, constants.CodeValidation, "审计查询参数校验失败："+err.Error())
		return
	}
	list, total, err := h.auditService.List(query.Page, query.PageSize, query.UserID, query.Module, query.Action)
	if err != nil {
		h.abort(c, err)
		return
	}
	response.OK(c, dto.PageResult{List: list, Total: total, Page: query.Page, PageSize: query.PageSize})
}

// abort 统一错误处理。
func (h *AuditHandler) abort(c *gin.Context, err error) {
	var appErr *util.AppError
	if errors.As(err, &appErr) {
		response.Fail(c, httpStatusFor(appErr.Code), appErr.Code, appErr.Message)
		return
	}
	h.logger.Error(fmt.Sprintf("audit handler error: %v", err))
	response.Fail(c, 500, constants.CodeInternal, "服务器内部错误，请稍后重试")
}
