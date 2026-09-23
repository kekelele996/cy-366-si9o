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

// UserHandler 用户管理接口处理器。
type UserHandler struct {
	userService *service.UserService
	logger      *slog.Logger
}

// NewUserHandler 构造用户管理接口处理器。
func NewUserHandler(userService *service.UserService, logger *slog.Logger) *UserHandler {
	return &UserHandler{userService: userService, logger: logger}
}

// Create 创建用户。
func (h *UserHandler) Create(c *gin.Context) {
	var req dto.CreateUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, constants.CodeValidation, "创建用户参数校验失败："+err.Error())
		return
	}
	user, err := h.userService.Create(&req)
	if err != nil {
		h.abort(c, err)
		return
	}
	response.OKMessage(c, constants.MsgCreateSuccess, user)
}

// Update 更新用户。
func (h *UserHandler) Update(c *gin.Context) {
	var idReq dto.IDReq
	if err := c.ShouldBindUri(&idReq); err != nil {
		response.Fail(c, 400, constants.CodeValidation, "用户 ID 无效")
		return
	}
	var req dto.UpdateUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, constants.CodeValidation, "更新用户参数校验失败："+err.Error())
		return
	}
	user, err := h.userService.Update(idReq.ID, &req)
	if err != nil {
		h.abort(c, err)
		return
	}
	response.OKMessage(c, constants.MsgUpdateSuccess, user)
}

// Delete 删除用户。
func (h *UserHandler) Delete(c *gin.Context) {
	var idReq dto.IDReq
	if err := c.ShouldBindUri(&idReq); err != nil {
		response.Fail(c, 400, constants.CodeValidation, "用户 ID 无效")
		return
	}
	if err := h.userService.Delete(idReq.ID); err != nil {
		h.abort(c, err)
		return
	}
	response.OKMessage(c, constants.MsgDeleteSuccess, nil)
}

// List 分页查询用户。
func (h *UserHandler) List(c *gin.Context) {
	var page dto.PageQuery
	if err := c.ShouldBindQuery(&page); err != nil {
		response.Fail(c, 400, constants.CodeValidation, "分页参数校验失败："+err.Error())
		return
	}
	users, total, err := h.userService.List(page.Page, page.PageSize)
	if err != nil {
		h.abort(c, err)
		return
	}
	response.OK(c, dto.PageResult{List: users, Total: total, Page: page.Page, PageSize: page.PageSize})
}

// Get 用户详情。
func (h *UserHandler) Get(c *gin.Context) {
	var idReq dto.IDReq
	if err := c.ShouldBindUri(&idReq); err != nil {
		response.Fail(c, 400, constants.CodeValidation, "用户 ID 无效")
		return
	}
	user, err := h.userService.GetByID(idReq.ID)
	if err != nil {
		h.abort(c, err)
		return
	}
	response.OK(c, user)
}

// abort 统一错误处理。
func (h *UserHandler) abort(c *gin.Context, err error) {
	var appErr *util.AppError
	if errors.As(err, &appErr) {
		response.Fail(c, httpStatusFor(appErr.Code), appErr.Code, appErr.Message)
		return
	}
	h.logger.Error(fmt.Sprintf("user handler error: %v", err))
	response.Fail(c, 500, constants.CodeInternal, "服务器内部错误，请稍后重试")
}
