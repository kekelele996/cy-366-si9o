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

// RechargeHandler 充值/时长包接口处理器。
type RechargeHandler struct {
	rechargeService *service.RechargeService
	logger          *slog.Logger
}

// NewRechargeHandler 构造充值接口处理器。
func NewRechargeHandler(rechargeService *service.RechargeService, logger *slog.Logger) *RechargeHandler {
	return &RechargeHandler{rechargeService: rechargeService, logger: logger}
}

// Recharge 会员充值。
func (h *RechargeHandler) Recharge(c *gin.Context) {
	var req dto.RechargeReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, constants.CodeValidation, "充值参数校验失败："+err.Error())
		return
	}
	operatorID, _ := c.Get("user_id")
	opID, _ := operatorID.(uint)
	if err := h.rechargeService.Recharge(&req, opID); err != nil {
		h.abort(c, err)
		return
	}
	response.OKMessage(c, constants.MsgRechargeOK, nil)
}

// BuyPackage 购买时长包。
func (h *RechargeHandler) BuyPackage(c *gin.Context) {
	var req dto.BuyPackageReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, constants.CodeValidation, "购买时长包参数校验失败："+err.Error())
		return
	}
	userID, _ := c.Get("user_id")
	uid, _ := userID.(uint)
	order, err := h.rechargeService.BuyPackage(uid, &req)
	if err != nil {
		h.abort(c, err)
		return
	}
	response.OKMessage(c, constants.MsgBuyPackageOK, order)
}

// ListRecharges 充值记录。
func (h *RechargeHandler) ListRecharges(c *gin.Context) {
	var page dto.PageQuery
	if err := c.ShouldBindQuery(&page); err != nil {
		response.Fail(c, 400, constants.CodeValidation, "分页参数校验失败："+err.Error())
		return
	}
	userID, _ := c.Get("user_id")
	uid, _ := userID.(uint)
	list, total, err := h.rechargeService.ListRecharges(uid, page.Page, page.PageSize)
	if err != nil {
		h.abort(c, err)
		return
	}
	response.OK(c, dto.PageResult{List: list, Total: total, Page: page.Page, PageSize: page.PageSize})
}

// ListOrders 时长包订单。
func (h *RechargeHandler) ListOrders(c *gin.Context) {
	var page dto.PageQuery
	if err := c.ShouldBindQuery(&page); err != nil {
		response.Fail(c, 400, constants.CodeValidation, "分页参数校验失败："+err.Error())
		return
	}
	userID, _ := c.Get("user_id")
	uid, _ := userID.(uint)
	list, total, err := h.rechargeService.ListOrders(uid, page.Page, page.PageSize)
	if err != nil {
		h.abort(c, err)
		return
	}
	response.OK(c, dto.PageResult{List: list, Total: total, Page: page.Page, PageSize: page.PageSize})
}

// abort 统一错误处理。
func (h *RechargeHandler) abort(c *gin.Context, err error) {
	var appErr *util.AppError
	if errors.As(err, &appErr) {
		response.Fail(c, httpStatusFor(appErr.Code), appErr.Code, appErr.Message)
		return
	}
	h.logger.Error(fmt.Sprintf("recharge handler error: %v", err))
	response.Fail(c, 500, constants.CodeInternal, "服务器内部错误，请稍后重试")
}
