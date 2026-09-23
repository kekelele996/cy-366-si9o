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

// AuthHandler 认证接口处理器。
type AuthHandler struct {
	authService *service.AuthService
	userService *service.UserService
	logger      *slog.Logger
}

// NewAuthHandler 构造认证接口处理器。
func NewAuthHandler(authService *service.AuthService, userService *service.UserService, logger *slog.Logger) *AuthHandler {
	return &AuthHandler{authService: authService, userService: userService, logger: logger}
}

// Register 注册。
func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, constants.CodeValidation, "注册参数校验失败："+err.Error())
		return
	}
	resp, err := h.authService.Register(&req)
	if err != nil {
		h.abort(c, err)
		return
	}
	response.OKMessage(c, constants.MsgRegisterOK, resp)
}

// Login 登录。
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, constants.CodeValidation, "登录参数校验失败："+err.Error())
		return
	}
	resp, err := h.authService.Login(&req)
	if err != nil {
		h.abort(c, err)
		return
	}
	response.OKMessage(c, constants.MsgLoginSuccess, resp)
}

// Profile 当前用户信息。
func (h *AuthHandler) Profile(c *gin.Context) {
	userID, _ := c.Get("user_id")
	uid, _ := userID.(uint)
	user, err := h.userService.GetByID(uid)
	if err != nil {
		h.abort(c, err)
		return
	}
	response.OK(c, user)
}

// abort 将服务错误转为统一响应（再次包装）。
func (h *AuthHandler) abort(c *gin.Context, err error) {
	var appErr *util.AppError
	if errors.As(err, &appErr) {
		response.Fail(c, httpStatusFor(appErr.Code), appErr.Code, appErr.Message)
		return
	}
	h.logger.Error(fmt.Sprintf("auth handler error: %v", err))
	response.Fail(c, 500, constants.CodeInternal, "服务器内部错误，请稍后重试")
}
