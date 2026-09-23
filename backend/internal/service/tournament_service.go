package service

import (
	"errors"
	"fmt"
	"log/slog"
	"math/rand"
	"time"

	"gorm.io/gorm"

	"github.com/esportsbar/backend/internal/constants"
	"github.com/esportsbar/backend/internal/dto"
	"github.com/esportsbar/backend/internal/model"
	"github.com/esportsbar/backend/internal/repository"
	"github.com/esportsbar/backend/internal/util"
)

// TournamentService 赛事服务（含报名、战队、抽签、赛果）。
type TournamentService struct {
	tournamentRepo *repository.TournamentRepository
	teamRepo       *repository.TeamRepository
	regRepo        *repository.RegistrationRepository
	matchRepo      *repository.MatchRepository
	db             *gorm.DB
	logger         *slog.Logger
}

// NewTournamentService 构造赛事服务。
func NewTournamentService(
	tournamentRepo *repository.TournamentRepository,
	teamRepo *repository.TeamRepository,
	regRepo *repository.RegistrationRepository,
	matchRepo *repository.MatchRepository,
	db *gorm.DB,
	logger *slog.Logger,
) *TournamentService {
	return &TournamentService{tournamentRepo: tournamentRepo, teamRepo: teamRepo, regRepo: regRepo, matchRepo: matchRepo, db: db, logger: logger}
}

// Create 创建赛事。
func (s *TournamentService) Create(req *dto.CreateTournamentReq, createdBy uint) (*model.Tournament, error) {
	status := constants.TournamentDraft
	if req.RegisterStart != nil {
		status = constants.TournamentOpen
	}
	t := &model.Tournament{
		Name:          req.Name,
		GameType:      req.GameType,
		Description:   req.Description,
		RegisterStart: req.RegisterStart,
		RegisterEnd:   req.RegisterEnd,
		StartTime:     req.StartTime,
		MaxTeams:      req.MaxTeams,
		Status:        status,
		CreatedBy:     createdBy,
	}
	if t.MaxTeams <= 0 {
		t.MaxTeams = 16
	}
	if err := s.tournamentRepo.Create(t); err != nil {
		return nil, fmt.Errorf("tournament create: %w", err)
	}
	s.logger.Info(fmt.Sprintf(constants.LogTemplates["tournament_create_ok"], t.ID, t.Name))
	return t, nil
}

// Update 更新赛事。
func (s *TournamentService) Update(id uint, req *dto.UpdateTournamentReq) (*model.Tournament, error) {
	t, err := s.getTournament(id)
	if err != nil {
		return nil, err
	}
	if req.Name != "" {
		t.Name = req.Name
	}
	if req.Description != "" {
		t.Description = req.Description
	}
	if req.RegisterStart != nil {
		t.RegisterStart = req.RegisterStart
	}
	if req.RegisterEnd != nil {
		t.RegisterEnd = req.RegisterEnd
	}
	if req.StartTime != nil {
		t.StartTime = req.StartTime
	}
	if req.MaxTeams > 0 {
		t.MaxTeams = req.MaxTeams
	}
	if req.Status != "" {
		t.Status = req.Status
	}
	if err := s.tournamentRepo.Update(t); err != nil {
		return nil, fmt.Errorf("tournament update: %w", err)
	}
	if req.Status == constants.TournamentOpen {
		s.logger.Info(fmt.Sprintf(constants.LogTemplates["tournament_publish_ok"], id))
	}
	return t, nil
}

// Delete 删除赛事。
func (s *TournamentService) Delete(id uint) error {
	if err := s.tournamentRepo.Delete(id); err != nil {
		return fmt.Errorf("tournament delete: %w", err)
	}
	return nil
}

// List 分页查询赛事。
func (s *TournamentService) List(page, pageSize int, status string) ([]model.Tournament, int64, error) {
	return s.tournamentRepo.List(page, pageSize, status)
}

// GetByID 查询赛事详情。
func (s *TournamentService) GetByID(id uint) (*model.Tournament, error) {
	return s.getTournament(id)
}

// CreateTeam 创建战队。
func (s *TournamentService) CreateTeam(userID uint, req *dto.CreateTeamReq) (*model.Team, error) {
	team := &model.Team{
		Name:        req.Name,
		LeaderID:    userID,
		MemberCount: req.MemberCount,
	}
	if team.MemberCount <= 0 {
		team.MemberCount = 1
	}
	if err := s.teamRepo.Create(team); err != nil {
		return nil, fmt.Errorf("tournament create team: %w", err)
	}
	return team, nil
}

// ListMyTeams 查询我的战队。
func (s *TournamentService) ListMyTeams(userID uint) ([]model.Team, error) {
	return s.teamRepo.ListByUser(userID)
}

// Register 赛事报名（个人/战队）。
func (s *TournamentService) Register(userID uint, tournamentID uint, req *dto.TournamentRegisterReq) (*model.Registration, error) {
	t, err := s.getTournament(tournamentID)
	if err != nil {
		return nil, err
	}
	if t.Status != constants.TournamentOpen {
		return nil, util.NewAppError(constants.CodeTournament, "赛事未在报名中，无法报名")
	}
	cnt, err := s.regRepo.CountByTournamentUser(tournamentID, userID)
	if err != nil {
		return nil, fmt.Errorf("tournament register count: %w", err)
	}
	if cnt > 0 {
		return nil, util.NewAppError(constants.CodeConflict, "你已报名该赛事，请勿重复报名")
	}
	reg := &model.Registration{
		TournamentID: tournamentID,
		UserID:       userID,
		Mode:         req.Mode,
		Status:       constants.RegistrationJoined,
	}
	if req.Mode == constants.RegistrationTeam {
		if req.TeamID == 0 {
			return nil, util.NewAppError(constants.CodeValidation, "战队报名必须选择战队")
		}
		team, err := s.teamRepo.FindByID(req.TeamID)
		if err != nil {
			if errors.Is(err, repository.ErrNotFound) {
				return nil, util.NewAppError(constants.CodeNotFound, "战队不存在")
			}
			return nil, fmt.Errorf("tournament register find team: %w", err)
		}
		if team.LeaderID != userID {
			return nil, util.NewAppError(constants.CodeForbidden, "只有战队队长可以代表战队报名")
		}
		reg.TeamID = team.ID
	}
	if err := s.regRepo.Create(reg); err != nil {
		return nil, fmt.Errorf("tournament register create: %w", err)
	}
	s.logger.Info(fmt.Sprintf(constants.LogTemplates["registration_create_ok"], tournamentID, req.Mode))
	return reg, nil
}

// ListRegistrations 查询赛事报名列表。
func (s *TournamentService) ListRegistrations(tournamentID uint) ([]model.Registration, error) {
	return s.regRepo.ListByTournament(tournamentID)
}

// DrawGroups 自动抽签分组：将报名记录随机分到小组并生成比赛场次。
func (s *TournamentService) DrawGroups(tournamentID uint) (map[string]any, error) {
	t, err := s.getTournament(tournamentID)
	if err != nil {
		return nil, err
	}
	if t.Status != constants.TournamentOpen && t.Status != constants.TournamentReady {
		return nil, util.NewAppError(constants.CodeTournament, "赛事状态不允许抽签")
	}
	regs, err := s.regRepo.ListByTournament(tournamentID)
	if err != nil {
		return nil, fmt.Errorf("tournament draw list: %w", err)
	}
	if len(regs) < 2 {
		return nil, util.NewAppError(constants.CodeTournament, "报名人数不足，无法抽签")
	}
	rand.New(rand.NewSource(time.Now().UnixNano()))
	rand.Shuffle(len(regs), func(i, j int) { regs[i], regs[j] = regs[j], regs[i] })
	groupCount := (len(regs) + 1) / 2
	err = s.db.Transaction(func(tx *gorm.DB) error {
		for i := range regs {
			groupNo := i%groupCount + 1
			if err := s.regRepo.AssignGroup(tx, int(regs[i].ID), groupNo); err != nil {
				return err
			}
		}
		// 每组内两两生成首轮比赛。
		for g := 1; g <= groupCount; g++ {
			var groupRegs []*model.Registration
			for i := range regs {
				if regs[i].GroupNo == g {
					groupRegs = append(groupRegs, &regs[i])
				}
			}
			for i := 0; i+1 < len(groupRegs); i += 2 {
				m := &model.Match{
					TournamentID: tournamentID,
					Round:        1,
					GroupNo:      g,
					TeamAID:      groupRegs[i].TeamID,
					TeamBID:      groupRegs[i+1].TeamID,
					Status:       constants.MatchPending,
				}
				if err := s.matchRepo.Create(m); err != nil {
					return err
				}
			}
		}
		t.Status = constants.TournamentReady
		return s.tournamentRepo.Update(t)
	})
	if err != nil {
		return nil, fmt.Errorf("tournament draw tx: %w", err)
	}
	s.logger.Info(fmt.Sprintf(constants.LogTemplates["tournament_draw_ok"], tournamentID, groupCount))
	return map[string]any{"groups": groupCount, "registrations": len(regs)}, nil
}

// ListMatches 查询赛事比赛场次。
func (s *TournamentService) ListMatches(tournamentID uint) ([]model.Match, error) {
	return s.matchRepo.ListByTournament(tournamentID)
}

// SubmitMatchResult 提交比赛结果并记录战绩。
func (s *TournamentService) SubmitMatchResult(matchID uint, req *dto.SubmitMatchReq) (*model.Match, error) {
	m, err := s.matchRepo.FindByID(matchID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, util.NewAppError(constants.CodeNotFound, "比赛场次不存在")
		}
		return nil, fmt.Errorf("tournament match find: %w", err)
	}
	if m.Status == constants.MatchFinished {
		return nil, util.NewAppError(constants.CodeConflict, "比赛结果已提交，不可重复提交")
	}
	if req.WinnerID == 0 {
		if req.ScoreA > req.ScoreB {
			req.WinnerID = m.TeamAID
		} else if req.ScoreB > req.ScoreA {
			req.WinnerID = m.TeamBID
		}
	}
	m.ScoreA = req.ScoreA
	m.ScoreB = req.ScoreB
	m.WinnerID = req.WinnerID
	m.Status = constants.MatchFinished
	if err := s.matchRepo.Update(m); err != nil {
		return nil, fmt.Errorf("tournament match update: %w", err)
	}
	s.logger.Info(fmt.Sprintf(constants.LogTemplates["match_result_ok"], matchID, m.WinnerID))
	return m, nil
}

// getTournament 查询赛事并统一错误处理。
func (s *TournamentService) getTournament(id uint) (*model.Tournament, error) {
	t, err := s.tournamentRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, util.NewAppError(constants.CodeNotFound, "赛事不存在")
		}
		return nil, fmt.Errorf("tournament find: %w", err)
	}
	return t, nil
}
