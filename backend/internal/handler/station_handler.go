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

// StationHandler 机位接口处理器。
type StationHandler struct {
	stationService *service.StationService
	logger         *slog.Logger
}

// NewStationHandler 构造机位接口处理器。
func NewStationHandler(stationService *service.StationService, logger *slog.Logger) *StationHandler {
	return &StationHandler{stationService: stationService, logger: logger}
}

// Create 创建机位。
func (h *StationHandler) Create(c *gin.Context) {
	var req dto.CreateStationReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, constants.CodeValidation, "创建机位参数校验失败："+err.Error())
		return
	}
	station, err := h.stationService.Create(&req)
	if err != nil {
		h.abort(c, err)
		return
	}
	response.OKMessage(c, constants.MsgCreateSuccess, station)
}

// Update 更新机位。
func (h *StationHandler) Update(c *gin.Context) {
	var idReq dto.IDReq
	if err := c.ShouldBindUri(&idReq); err != nil {
		response.Fail(c, 400, constants.CodeValidation, "机位 ID 无效")
		return
	}
	var req dto.UpdateStationReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, constants.CodeValidation, "更新机位参数校验失败："+err.Error())
		return
	}
	station, err := h.stationService.Update(idReq.ID, &req)
	if err != nil {
		h.abort(c, err)
		return
	}
	response.OKMessage(c, constants.MsgUpdateSuccess, station)
}

// UpdateStatus 机位状态流转。
func (h *StationHandler) UpdateStatus(c *gin.Context) {
	var idReq dto.IDReq
	if err := c.ShouldBindUri(&idReq); err != nil {
		response.Fail(c, 400, constants.CodeValidation, "机位 ID 无效")
		return
	}
	var req dto.UpdateStationStatusReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, constants.CodeValidation, "机位状态参数校验失败："+err.Error())
		return
	}
	station, err := h.stationService.UpdateStatus(idReq.ID, &req)
	if err != nil {
		h.abort(c, err)
		return
	}
	response.OKMessage(c, constants.MsgUpdateSuccess, station)
}

// Delete 删除机位。
func (h *StationHandler) Delete(c *gin.Context) {
	var idReq dto.IDReq
	if err := c.ShouldBindUri(&idReq); err != nil {
		response.Fail(c, 400, constants.CodeValidation, "机位 ID 无效")
		return
	}
	if err := h.stationService.Delete(idReq.ID); err != nil {
		h.abort(c, err)
		return
	}
	response.OKMessage(c, constants.MsgDeleteSuccess, nil)
}

// List 分页查询机位。
func (h *StationHandler) List(c *gin.Context) {
	var query dto.StationQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Fail(c, 400, constants.CodeValidation, "机位查询参数校验失败："+err.Error())
		return
	}
	stations, total, err := h.stationService.List(&query)
	if err != nil {
		h.abort(c, err)
		return
	}
	response.OK(c, dto.PageResult{List: stations, Total: total, Page: query.Page, PageSize: query.PageSize})
}

// ListAll 查询全部机位（看板与下拉复用）。
func (h *StationHandler) ListAll(c *gin.Context) {
	stations, err := h.stationService.ListAll()
	if err != nil {
		h.abort(c, err)
		return
	}
	response.OK(c, stations)
}

// Get 机位详情。
func (h *StationHandler) Get(c *gin.Context) {
	var idReq dto.IDReq
	if err := c.ShouldBindUri(&idReq); err != nil {
		response.Fail(c, 400, constants.CodeValidation, "机位 ID 无效")
		return
	}
	station, err := h.stationService.GetByID(idReq.ID)
	if err != nil {
		h.abort(c, err)
		return
	}
	response.OK(c, station)
}

// abort 统一错误处理。
func (h *StationHandler) abort(c *gin.Context, err error) {
	var appErr *util.AppError
	if errors.As(err, &appErr) {
		response.Fail(c, httpStatusFor(appErr.Code), appErr.Code, appErr.Message)
		return
	}
	h.logger.Error(fmt.Sprintf("station handler error: %v", err))
	response.Fail(c, 500, constants.CodeInternal, "服务器内部错误，请稍后重试")
}
