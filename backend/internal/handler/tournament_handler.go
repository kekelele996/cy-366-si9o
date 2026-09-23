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

// TournamentHandler 赛事接口处理器。
type TournamentHandler struct {
	tournamentService *service.TournamentService
	logger            *slog.Logger
}

// NewTournamentHandler 构造赛事接口处理器。
func NewTournamentHandler(tournamentService *service.TournamentService, logger *slog.Logger) *TournamentHandler {
	return &TournamentHandler{tournamentService: tournamentService, logger: logger}
}

// Create 创建赛事。
func (h *TournamentHandler) Create(c *gin.Context) {
	var req dto.CreateTournamentReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, constants.CodeValidation, "创建赛事参数校验失败："+err.Error())
		return
	}
	userID, _ := c.Get("user_id")
	uid, _ := userID.(uint)
	t, err := h.tournamentService.Create(&req, uid)
	if err != nil {
		h.abort(c, err)
		return
	}
	response.OKMessage(c, constants.MsgCreateSuccess, t)
}

// Update 更新赛事。
func (h *TournamentHandler) Update(c *gin.Context) {
	var idReq dto.IDReq
	if err := c.ShouldBindUri(&idReq); err != nil {
		response.Fail(c, 400, constants.CodeValidation, "赛事 ID 无效")
		return
	}
	var req dto.UpdateTournamentReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, constants.CodeValidation, "更新赛事参数校验失败："+err.Error())
		return
	}
	t, err := h.tournamentService.Update(idReq.ID, &req)
	if err != nil {
		h.abort(c, err)
		return
	}
	response.OKMessage(c, constants.MsgUpdateSuccess, t)
}

// Delete 删除赛事。
func (h *TournamentHandler) Delete(c *gin.Context) {
	var idReq dto.IDReq
	if err := c.ShouldBindUri(&idReq); err != nil {
		response.Fail(c, 400, constants.CodeValidation, "赛事 ID 无效")
		return
	}
	if err := h.tournamentService.Delete(idReq.ID); err != nil {
		h.abort(c, err)
		return
	}
	response.OKMessage(c, constants.MsgDeleteSuccess, nil)
}

// List 分页查询赛事。
func (h *TournamentHandler) List(c *gin.Context) {
	var page dto.PageQuery
	if err := c.ShouldBindQuery(&page); err != nil {
		response.Fail(c, 400, constants.CodeValidation, "分页参数校验失败："+err.Error())
		return
	}
	status := c.DefaultQuery("status", "")
	list, total, err := h.tournamentService.List(page.Page, page.PageSize, status)
	if err != nil {
		h.abort(c, err)
		return
	}
	response.OK(c, dto.PageResult{List: list, Total: total, Page: page.Page, PageSize: page.PageSize})
}

// Get 赛事详情。
func (h *TournamentHandler) Get(c *gin.Context) {
	var idReq dto.IDReq
	if err := c.ShouldBindUri(&idReq); err != nil {
		response.Fail(c, 400, constants.CodeValidation, "赛事 ID 无效")
		return
	}
	t, err := h.tournamentService.GetByID(idReq.ID)
	if err != nil {
		h.abort(c, err)
		return
	}
	response.OK(c, t)
}

// CreateTeam 创建战队。
func (h *TournamentHandler) CreateTeam(c *gin.Context) {
	var req dto.CreateTeamReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, constants.CodeValidation, "创建战队参数校验失败："+err.Error())
		return
	}
	userID, _ := c.Get("user_id")
	uid, _ := userID.(uint)
	team, err := h.tournamentService.CreateTeam(uid, &req)
	if err != nil {
		h.abort(c, err)
		return
	}
	response.OKMessage(c, constants.MsgCreateSuccess, team)
}

// ListMyTeams 我的战队。
func (h *TournamentHandler) ListMyTeams(c *gin.Context) {
	userID, _ := c.Get("user_id")
	uid, _ := userID.(uint)
	list, err := h.tournamentService.ListMyTeams(uid)
	if err != nil {
		h.abort(c, err)
		return
	}
	response.OK(c, list)
}

// Register 赛事报名。
func (h *TournamentHandler) Register(c *gin.Context) {
	var idReq dto.IDReq
	if err := c.ShouldBindUri(&idReq); err != nil {
		response.Fail(c, 400, constants.CodeValidation, "赛事 ID 无效")
		return
	}
	var req dto.TournamentRegisterReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, constants.CodeValidation, "报名参数校验失败："+err.Error())
		return
	}
	userID, _ := c.Get("user_id")
	uid, _ := userID.(uint)
	reg, err := h.tournamentService.Register(uid, idReq.ID, &req)
	if err != nil {
		h.abort(c, err)
		return
	}
	response.OKMessage(c, constants.MsgCreateSuccess, reg)
}

// ListRegistrations 赛事报名列表。
func (h *TournamentHandler) ListRegistrations(c *gin.Context) {
	var idReq dto.IDReq
	if err := c.ShouldBindUri(&idReq); err != nil {
		response.Fail(c, 400, constants.CodeValidation, "赛事 ID 无效")
		return
	}
	list, err := h.tournamentService.ListRegistrations(idReq.ID)
	if err != nil {
		h.abort(c, err)
		return
	}
	response.OK(c, list)
}

// DrawGroups 抽签分组。
func (h *TournamentHandler) DrawGroups(c *gin.Context) {
	var idReq dto.IDReq
	if err := c.ShouldBindUri(&idReq); err != nil {
		response.Fail(c, 400, constants.CodeValidation, "赛事 ID 无效")
		return
	}
	result, err := h.tournamentService.DrawGroups(idReq.ID)
	if err != nil {
		h.abort(c, err)
		return
	}
	response.OKMessage(c, constants.MsgDrawOK, result)
}

// ListMatches 比赛场次。
func (h *TournamentHandler) ListMatches(c *gin.Context) {
	var idReq dto.IDReq
	if err := c.ShouldBindUri(&idReq); err != nil {
		response.Fail(c, 400, constants.CodeValidation, "赛事 ID 无效")
		return
	}
	list, err := h.tournamentService.ListMatches(idReq.ID)
	if err != nil {
		h.abort(c, err)
		return
	}
	response.OK(c, list)
}

// SubmitMatchResult 提交比赛结果。
func (h *TournamentHandler) SubmitMatchResult(c *gin.Context) {
	var idReq dto.IDReq
	if err := c.ShouldBindUri(&idReq); err != nil {
		response.Fail(c, 400, constants.CodeValidation, "比赛 ID 无效")
		return
	}
	var req dto.SubmitMatchReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, constants.CodeValidation, "比赛结果参数校验失败："+err.Error())
		return
	}
	m, err := h.tournamentService.SubmitMatchResult(idReq.ID, &req)
	if err != nil {
		h.abort(c, err)
		return
	}
	response.OKMessage(c, constants.MsgUpdateSuccess, m)
}

// abort 统一错误处理。
func (h *TournamentHandler) abort(c *gin.Context, err error) {
	var appErr *util.AppError
	if errors.As(err, &appErr) {
		response.Fail(c, httpStatusFor(appErr.Code), appErr.Code, appErr.Message)
		return
	}
	h.logger.Error(fmt.Sprintf("tournament handler error: %v", err))
	response.Fail(c, 500, constants.CodeInternal, "服务器内部错误，请稍后重试")
}
