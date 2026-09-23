package handler

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/gin-gonic/gin"

	"github.com/esportsbar/backend/internal/constants"
	"github.com/esportsbar/backend/internal/service"
	"github.com/esportsbar/backend/internal/util"
	"github.com/esportsbar/backend/pkg/response"
)

// DashboardHandler 看板接口处理器。
type DashboardHandler struct {
	dashboardService *service.DashboardService
	logger           *slog.Logger
}

// NewDashboardHandler 构造看板接口处理器。
func NewDashboardHandler(dashboardService *service.DashboardService, logger *slog.Logger) *DashboardHandler {
	return &DashboardHandler{dashboardService: dashboardService, logger: logger}
}

// Summary 看板汇总。
func (h *DashboardHandler) Summary(c *gin.Context) {
	summary, err := h.dashboardService.GetSummary()
	if err != nil {
		var appErr *util.AppError
		if errors.As(err, &appErr) {
			response.Fail(c, httpStatusFor(appErr.Code), appErr.Code, appErr.Message)
			return
		}
		h.logger.Error(fmt.Sprintf("dashboard handler error: %v", err))
		response.Fail(c, 500, constants.CodeInternal, "服务器内部错误，请稍后重试")
		return
	}
	response.OK(c, summary)
}
