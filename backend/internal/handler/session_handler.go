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

// SessionHandler 上机记录接口处理器。
type SessionHandler struct {
	sessionService *service.SessionService
	logger         *slog.Logger
}

// NewSessionHandler 构造上机记录接口处理器。
func NewSessionHandler(sessionService *service.SessionService, logger *slog.Logger) *SessionHandler {
	return &SessionHandler{sessionService: sessionService, logger: logger}
}

// Start 上机开机。
func (h *SessionHandler) Start(c *gin.Context) {
	var req dto.StartSessionReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, constants.CodeValidation, "开机参数校验失败："+err.Error())
		return
	}
	userID, _ := c.Get("user_id")
	uid, _ := userID.(uint)
	sess, err := h.sessionService.Start(uid, &req)
	if err != nil {
		h.abort(c, err)
		return
	}
	response.OKMessage(c, constants.MsgCheckInOK, sess)
}

// Renew 续费。
func (h *SessionHandler) Renew(c *gin.Context) {
	var idReq dto.IDReq
	if err := c.ShouldBindUri(&idReq); err != nil {
		response.Fail(c, 400, constants.CodeValidation, "上机记录 ID 无效")
		return
	}
	var req dto.RenewSessionReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, constants.CodeValidation, "续费参数校验失败："+err.Error())
		return
	}
	userID, _ := c.Get("user_id")
	uid, _ := userID.(uint)
	sess, err := h.sessionService.Renew(uid, idReq.ID, &req)
	if err != nil {
		h.abort(c, err)
		return
	}
	response.OKMessage(c, constants.MsgRenewOK, sess)
}

// End 下机结算。
func (h *SessionHandler) End(c *gin.Context) {
	var idReq dto.IDReq
	if err := c.ShouldBindUri(&idReq); err != nil {
		response.Fail(c, 400, constants.CodeValidation, "上机记录 ID 无效")
		return
	}
	var req dto.EndSessionReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, constants.CodeValidation, "下机参数校验失败："+err.Error())
		return
	}
	userID, _ := c.Get("user_id")
	uid, _ := userID.(uint)
	sess, err := h.sessionService.End(uid, idReq.ID, &req)
	if err != nil {
		h.abort(c, err)
		return
	}
	response.OKMessage(c, constants.MsgCheckoutOK, sess)
}

// List 分页查询上机记录。
func (h *SessionHandler) List(c *gin.Context) {
	var page dto.PageQuery
	if err := c.ShouldBindQuery(&page); err != nil {
		response.Fail(c, 400, constants.CodeValidation, "分页参数校验失败："+err.Error())
		return
	}
	status := c.DefaultQuery("status", "")
	userIDStr := c.Query("user_id")
	var uid uint
	if userIDStr != "" {
		fmt.Sscanf(userIDStr, "%d", &uid)
	}
	list, total, err := h.sessionService.List(page.Page, page.PageSize, uid, status)
	if err != nil {
		h.abort(c, err)
		return
	}
	response.OK(c, dto.PageResult{List: list, Total: total, Page: page.Page, PageSize: page.PageSize})
}

// Rank 排行榜。
func (h *SessionHandler) Rank(c *gin.Context) {
	var query dto.RankQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Fail(c, 400, constants.CodeValidation, "排行榜参数校验失败："+err.Error())
		return
	}
	list, err := h.sessionService.Rank(&query)
	if err != nil {
		h.abort(c, err)
		return
	}
	response.OK(c, list)
}

// abort 统一错误处理。
func (h *SessionHandler) abort(c *gin.Context, err error) {
	var appErr *util.AppError
	if errors.As(err, &appErr) {
		response.Fail(c, httpStatusFor(appErr.Code), appErr.Code, appErr.Message)
		return
	}
	h.logger.Error(fmt.Sprintf("session handler error: %v", err))
	response.Fail(c, 500, constants.CodeInternal, "服务器内部错误，请稍后重试")
}
