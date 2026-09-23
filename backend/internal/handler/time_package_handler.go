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

// TimePackageHandler 时长包管理接口处理器。
type TimePackageHandler struct {
	packageService *service.TimePackageService
	logger         *slog.Logger
}

// NewTimePackageHandler 构造时长包管理接口处理器。
func NewTimePackageHandler(packageService *service.TimePackageService, logger *slog.Logger) *TimePackageHandler {
	return &TimePackageHandler{packageService: packageService, logger: logger}
}

// Create 创建时长包。
func (h *TimePackageHandler) Create(c *gin.Context) {
	var req dto.CreatePackageReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, constants.CodeValidation, "创建时长包参数校验失败："+err.Error())
		return
	}
	pkg, err := h.packageService.Create(&req)
	if err != nil {
		h.abort(c, err)
		return
	}
	response.OKMessage(c, constants.MsgCreateSuccess, pkg)
}

// Update 更新时长包。
func (h *TimePackageHandler) Update(c *gin.Context) {
	var idReq dto.IDReq
	if err := c.ShouldBindUri(&idReq); err != nil {
		response.Fail(c, 400, constants.CodeValidation, "时长包 ID 无效")
		return
	}
	var req dto.UpdatePackageReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, constants.CodeValidation, "更新时长包参数校验失败："+err.Error())
		return
	}
	pkg, err := h.packageService.Update(idReq.ID, &req)
	if err != nil {
		h.abort(c, err)
		return
	}
	response.OKMessage(c, constants.MsgUpdateSuccess, pkg)
}

// Delete 删除时长包。
func (h *TimePackageHandler) Delete(c *gin.Context) {
	var idReq dto.IDReq
	if err := c.ShouldBindUri(&idReq); err != nil {
		response.Fail(c, 400, constants.CodeValidation, "时长包 ID 无效")
		return
	}
	if err := h.packageService.Delete(idReq.ID); err != nil {
		h.abort(c, err)
		return
	}
	response.OKMessage(c, constants.MsgDeleteSuccess, nil)
}

// List 分页查询时长包。
func (h *TimePackageHandler) List(c *gin.Context) {
	var page dto.PageQuery
	if err := c.ShouldBindQuery(&page); err != nil {
		response.Fail(c, 400, constants.CodeValidation, "分页参数校验失败："+err.Error())
		return
	}
	list, total, err := h.packageService.List(page.Page, page.PageSize)
	if err != nil {
		h.abort(c, err)
		return
	}
	response.OK(c, dto.PageResult{List: list, Total: total, Page: page.Page, PageSize: page.PageSize})
}

// ListActive 在售时长包。
func (h *TimePackageHandler) ListActive(c *gin.Context) {
	list, err := h.packageService.ListActive()
	if err != nil {
		h.abort(c, err)
		return
	}
	response.OK(c, list)
}

// Get 时长包详情。
func (h *TimePackageHandler) Get(c *gin.Context) {
	var idReq dto.IDReq
	if err := c.ShouldBindUri(&idReq); err != nil {
		response.Fail(c, 400, constants.CodeValidation, "时长包 ID 无效")
		return
	}
	pkg, err := h.packageService.GetByID(idReq.ID)
	if err != nil {
		h.abort(c, err)
		return
	}
	response.OK(c, pkg)
}

// abort 统一错误处理。
func (h *TimePackageHandler) abort(c *gin.Context, err error) {
	var appErr *util.AppError
	if errors.As(err, &appErr) {
		response.Fail(c, httpStatusFor(appErr.Code), appErr.Code, appErr.Message)
		return
	}
	h.logger.Error(fmt.Sprintf("time package handler error: %v", err))
	response.Fail(c, 500, constants.CodeInternal, "服务器内部错误，请稍后重试")
}
