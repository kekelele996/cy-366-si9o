package service

import (
	"errors"
	"fmt"
	"log/slog"
	"time"

	"gorm.io/gorm"

	"github.com/esportsbar/backend/internal/constants"
	"github.com/esportsbar/backend/internal/dto"
	"github.com/esportsbar/backend/internal/model"
	"github.com/esportsbar/backend/internal/repository"
	"github.com/esportsbar/backend/internal/util"
)

// ReservationService 机位预约服务。
type ReservationService struct {
	reservationRepo *repository.ReservationRepository
	stationService  *StationService
	db              *gorm.DB
	logger          *slog.Logger
}

// NewReservationService 构造预约服务。
func NewReservationService(
	reservationRepo *repository.ReservationRepository,
	stationService *StationService,
	db *gorm.DB,
	logger *slog.Logger,
) *ReservationService {
	return &ReservationService{reservationRepo: reservationRepo, stationService: stationService, db: db, logger: logger}
}

// Create 创建预约：时段空闲则直接确认；时段冲突则进入候补队列并返回当前顺位。
func (s *ReservationService) Create(userID uint, req *dto.CreateReservationReq) (*model.Reservation, error) {
	if !req.EndTime.After(req.StartTime) {
		return nil, util.NewAppError(constants.CodeValidation, "预约结束时间必须晚于开始时间")
	}
	station, err := s.stationService.GetByID(req.StationID)
	if err != nil {
		return nil, err
	}
	if station.Status != constants.StationIdle && station.Status != constants.StationReserved {
		return nil, util.NewAppError(constants.CodeStationBusy, "机位当前不可预约，请选择其他机位")
	}
	res := &model.Reservation{
		UserID:    userID,
		StationID: req.StationID,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		Status:    constants.ReservationConfirmed,
		Remark:    req.Remark,
	}
	err = s.db.Transaction(func(tx *gorm.DB) error {
		// 先锁机位：同机位的创建/取消经行锁串行化，杜绝并发下的重复确认与状态错乱。
		locked, err := s.stationService.LockForUpdate(tx, req.StationID)
		if err != nil {
			return err
		}
		if locked.Status != constants.StationIdle && locked.Status != constants.StationReserved {
			return util.NewAppError(constants.CodeStationBusy, "机位当前不可预约，请选择其他机位")
		}
		cnt, err := s.reservationRepo.CountConflictTx(tx, req.StationID, req.StartTime, req.EndTime, 0)
		if err != nil {
			return fmt.Errorf("reservation count conflict: %w", err)
		}
		if cnt > 0 {
			// 热门机位已被订走：不拒绝会员，进入候补队列。
			res.Status = constants.ReservationWaitlisted
			if err := s.reservationRepo.CreateTx(tx, res); err != nil {
				return err
			}
			queue, err := s.reservationRepo.ListWaitlistedByStationTx(tx, req.StationID)
			if err != nil {
				return err
			}
			res.WaitlistPosition = waitlistPositionInQueue(res, queue)
			return nil
		}
		if locked.Status == constants.StationIdle {
			locked.Status = constants.StationReserved
			if err := tx.Save(locked).Error; err != nil {
				return err
			}
		}
		return s.reservationRepo.CreateTx(tx, res)
	})
	if err != nil {
		return nil, fmt.Errorf("reservation create tx: %w", err)
	}
	if res.Status == constants.ReservationWaitlisted {
		s.logger.Info(fmt.Sprintf(constants.LogTemplates["reservation_wait_ok"], res.ID, userID, req.StationID, res.WaitlistPosition))
	} else {
		s.logger.Info(fmt.Sprintf(constants.LogTemplates["reservation_create_ok"], userID, req.StationID, req.StartTime.Format("2006-01-02 15:04")))
	}
	return res, nil
}

// Confirm 确认预约（staff/admin）。
func (s *ReservationService) Confirm(id uint) (*model.Reservation, error) {
	res, err := s.getReservation(id)
	if err != nil {
		return nil, err
	}
	if res.Status != constants.ReservationPending && res.Status != constants.ReservationConfirmed {
		return nil, util.NewAppError(constants.CodeReservation, "仅待确认或已确认的预约可以确认")
	}
	res.Status = constants.ReservationConfirmed
	err = s.db.Transaction(func(tx *gorm.DB) error {
		return s.reservationRepo.UpdateTx(tx, res)
	})
	if err != nil {
		return nil, fmt.Errorf("reservation confirm tx: %w", err)
	}
	s.logger.Info(fmt.Sprintf(constants.LogTemplates["reservation_confirm_ok"], id))
	return res, nil
}

// Cancel 取消预约：已确认预约取消时，在同事务内兑现至多一个候补并回填机位状态。
func (s *ReservationService) Cancel(id uint, userID uint) (*model.Reservation, error) {
	res, err := s.getReservation(id)
	if err != nil {
		return nil, err
	}
	if !cancellable(res.Status) {
		return nil, util.NewAppError(constants.CodeReservation, "当前状态不可取消")
	}
	fromStatus := res.Status
	var promoted *model.Reservation
	err = s.db.Transaction(func(tx *gorm.DB) error {
		// 统一按“机位 → 预约”顺序加锁，避免并发取消交叉等锁造成死锁。
		lockedStation, err := s.stationService.LockForUpdate(tx, res.StationID)
		if err != nil {
			return err
		}
		locked, err := s.reservationRepo.FindByIDForUpdateTx(tx, id)
		if err != nil {
			return err
		}
		if !cancellable(locked.Status) {
			// 并发/重复取消：状态已被前一个请求改走，直接冲突失败，不会再次兑现候补。
			return util.NewAppError(constants.CodeReservation, "当前状态不可取消")
		}
		rows, err := s.reservationRepo.UpdateStatusTx(tx, id,
			[]string{constants.ReservationPending, constants.ReservationWaitlisted, constants.ReservationConfirmed, constants.ReservationCheckedIn},
			constants.ReservationCancelled)
		if err != nil {
			return err
		}
		if rows == 0 {
			return util.NewAppError(constants.CodeReservation, "当前状态不可取消")
		}
		fromStatus = locked.Status
		locked.Status = constants.ReservationCancelled
		res = locked

		if fromStatus == constants.ReservationConfirmed {
			// 仅“已确认”取消触发候补兑现；一次取消只兑现一人。
			promoted, err = s.promoteWaitlistTx(tx, locked.StationID, id)
			if err != nil {
				return err
			}
		}
		return reconcileStationReserved(tx, s.reservationRepo, s.stationService, lockedStation, promoted)
	})
	if err != nil {
		return nil, fmt.Errorf("reservation cancel tx: %w", err)
	}
	if promoted != nil {
		s.logger.Info(fmt.Sprintf(constants.LogTemplates["reservation_promote_ok"], promoted.ID, promoted.UserID, promoted.StationID, id))
	}
	s.logger.Info(fmt.Sprintf(constants.LogTemplates["reservation_cancel_ok"], id))
	return res, nil
}

// CheckIn 到店扫码开机：预约状态流转为 checked_in，机位置为使用中。
func (s *ReservationService) CheckIn(id uint) (*model.Reservation, error) {
	res, err := s.getReservation(id)
	if err != nil {
		return nil, err
	}
	if res.Status != constants.ReservationConfirmed {
		return nil, util.NewAppError(constants.CodeReservation, "仅已确认的预约可以开机")
	}
	err = s.db.Transaction(func(tx *gorm.DB) error {
		lockedStation, err := s.stationService.LockForUpdate(tx, res.StationID)
		if err != nil {
			return err
		}
		locked, err := s.reservationRepo.FindByIDForUpdateTx(tx, id)
		if err != nil {
			return err
		}
		if locked.Status != constants.ReservationConfirmed {
			return util.NewAppError(constants.CodeReservation, "仅已确认的预约可以开机")
		}
		locked.Status = constants.ReservationCheckedIn
		if err := s.reservationRepo.UpdateTx(tx, locked); err != nil {
			return err
		}
		lockedStation.Status = constants.StationUsing
		return tx.Save(lockedStation).Error
	})
	if err != nil {
		return nil, fmt.Errorf("reservation checkin tx: %w", err)
	}
	s.logger.Info(fmt.Sprintf(constants.LogTemplates["reservation_checkin_ok"], id))
	return res, nil
}

// List 分页查询预约，并为候补预约回填当前候补顺位。
func (s *ReservationService) List(query *dto.ReservationQuery) ([]model.Reservation, int64, error) {
	page := query.Page
	pageSize := query.PageSize
	if page <= 0 {
		page = constants.DefaultPage
	}
	if pageSize <= 0 {
		pageSize = constants.DefaultPageSize
	}
	list, total, err := s.reservationRepo.List(page, pageSize, query.Status, query.UserID)
	if err != nil {
		return nil, 0, err
	}
	if err := s.fillWaitlistPositions(list); err != nil {
		return nil, 0, fmt.Errorf("reservation list positions: %w", err)
	}
	return list, total, nil
}

// promoteWaitlistTx 兑现候补：在同机位、与被取消时段重叠的候补里，
// 按提交顺序挑选第一个不再与任何有效预约冲突者转为 confirmed；一次只兑现一人。
func (s *ReservationService) promoteWaitlistTx(tx *gorm.DB, stationID, cancelledID uint) (*model.Reservation, error) {
	cancelled, err := s.reservationRepo.FindByIDTx(tx, cancelledID)
	if err != nil {
		return nil, err
	}
	queue, err := s.reservationRepo.ListWaitlistedByStationTx(tx, stationID)
	if err != nil {
		return nil, err
	}
	candidate := pickPromotable(queue, cancelled, func(c *model.Reservation) (bool, error) {
		cnt, err := s.reservationRepo.CountActiveConflictTx(tx, stationID, c.StartTime, c.EndTime, cancelledID)
		if err != nil {
			return false, err
		}
		return cnt > 0, nil
	})
	if candidate == nil {
		return nil, nil
	}
	// 条件更新兜底：并发下只有一个事务能把该候补从 waitlisted 改成 confirmed。
	rows, err := s.reservationRepo.UpdateStatusTx(tx, candidate.ID,
		[]string{constants.ReservationWaitlisted}, constants.ReservationConfirmed)
	if err != nil {
		return nil, err
	}
	if rows == 0 {
		// 已被并发事务兑现或取消，本次取消不再补选第二人（一次取消至多兑现一个候补）。
		return nil, nil
	}
	candidate.Status = constants.ReservationConfirmed
	candidate.WaitlistPosition = 0
	return candidate, nil
}

// pickPromotable 按提交顺序选出第一个与 cancelled 时段重叠且 hasConflict=false 的候补；
// 仍冲突的候补跳过但不淘汰，前面的人兑现或退出后顺位自然前移。
func pickPromotable(queue []model.Reservation, cancelled *model.Reservation, hasConflict func(*model.Reservation) (bool, error)) *model.Reservation {
	for i := range queue {
		candidate := &queue[i]
		if !overlaps(candidate.StartTime, candidate.EndTime, cancelled.StartTime, cancelled.EndTime) {
			continue
		}
		conflict, err := hasConflict(candidate)
		if err != nil || conflict {
			continue
		}
		return candidate
	}
	return nil
}

// fillWaitlistPositions 批量回填列表中候补预约的当前顺位，避免 N+1 查询。
func (s *ReservationService) fillWaitlistPositions(list []model.Reservation) error {
	hasWaitlisted := false
	stationIDs := map[uint]struct{}{}
	for i := range list {
		if list[i].Status == constants.ReservationWaitlisted {
			hasWaitlisted = true
			stationIDs[list[i].StationID] = struct{}{}
		}
	}
	if !hasWaitlisted {
		return nil
	}
	waitlist, err := s.reservationRepo.ListWaitlistedAll(s.db)
	if err != nil {
		return err
	}
	queues := map[uint][]model.Reservation{}
	for i := range waitlist {
		if _, ok := stationIDs[waitlist[i].StationID]; ok {
			queues[waitlist[i].StationID] = append(queues[waitlist[i].StationID], waitlist[i])
		}
	}
	for i := range list {
		if list[i].Status == constants.ReservationWaitlisted {
			list[i].WaitlistPosition = waitlistPositionInQueue(&list[i], queues[list[i].StationID])
		}
	}
	return nil
}

// getReservation 查询预约并统一处理错误。
func (s *ReservationService) getReservation(id uint) (*model.Reservation, error) {
	res, err := s.reservationRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, util.NewAppError(constants.CodeNotFound, "预约记录不存在")
		}
		return nil, fmt.Errorf("reservation find: %w", err)
	}
	return res, nil
}

// cancellable 允许取消的预约状态。
func cancellable(status string) bool {
	switch status {
	case constants.ReservationPending, constants.ReservationWaitlisted, constants.ReservationConfirmed, constants.ReservationCheckedIn:
		return true
	}
	return false
}

// overlaps 判断两个半开时段是否重叠（start 相同视为相邻不冲突，与 SQL 中 < / > 判定一致）。
func overlaps(startA, endA, startB, endB time.Time) bool {
	return startA.Before(endB) && endA.After(startB)
}

// waitlistPositionInQueue 计算候补在同机位队列中的顺位：
// 提交时间更早（created_at, id 决胜）且时段重叠的候补人数 + 1。
func waitlistPositionInQueue(target *model.Reservation, queue []model.Reservation) int {
	position := 1
	for i := range queue {
		c := &queue[i]
		if c.ID == target.ID {
			continue
		}
		if earlier(c, target) && overlaps(c.StartTime, c.EndTime, target.StartTime, target.EndTime) {
			position++
		}
	}
	return position
}

// earlier 判断 a 是否比 b 更早提交（created_at 相同则 id 更小者在前）。
func earlier(a, b *model.Reservation) bool {
	if a.CreatedAt.Equal(b.CreatedAt) {
		return a.ID < b.ID
	}
	return a.CreatedAt.Before(b.CreatedAt)
}

// reconcileStationReserved 取消后回填机位状态：有有效预约或刚兑现候补则保持 reserved，
// 仅在确无占用时把 reserved 释放回 idle，避免并发取消把机位状态改乱。
func reconcileStationReserved(tx *gorm.DB, repo *repository.ReservationRepository, svc *StationService, station *model.Station, promoted *model.Reservation) error {
	if promoted != nil {
		return nil
	}
	cnt, err := repo.CountActiveByStationTx(tx, station.ID)
	if err != nil {
		return err
	}
	if cnt == 0 && station.Status == constants.StationReserved {
		station.Status = constants.StationIdle
		return tx.Save(station).Error
	}
	return nil
}
