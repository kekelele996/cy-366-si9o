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

// Create 创建预约：机位时段空闲则直接确认并占用机位；与已确认/已开机预约冲突时转为候补排队。
// 并发安全：事务内先锁定机位，再复查机位状态与时段冲突。
func (s *ReservationService) Create(userID uint, req *dto.CreateReservationReq) (*model.Reservation, error) {
	if !req.EndTime.After(req.StartTime) {
		return nil, util.NewAppError(constants.CodeValidation, "预约结束时间必须晚于开始时间")
	}
	if _, err := s.stationService.GetByID(req.StationID); err != nil {
		return nil, err
	}
	res := &model.Reservation{
		UserID:    userID,
		StationID: req.StationID,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		Remark:    req.Remark,
	}
	err := s.db.Transaction(func(tx *gorm.DB) error {
		// 所有写机位/预约的事务均按“先锁机位、再锁预约”的顺序加锁，避免并发死锁。
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
			// 热门机位时段被占用：候补不占机位，进入同机位候补队列。
			res.Status = constants.ReservationWaitlisted
			if err := s.reservationRepo.CreateTx(tx, res); err != nil {
				return fmt.Errorf("reservation waitlist create: %w", err)
			}
			return nil
		}
		res.Status = constants.ReservationConfirmed
		if locked.Status == constants.StationIdle {
			locked.Status = constants.StationReserved
			if err := tx.Save(locked).Error; err != nil {
				return err
			}
		}
		if err := s.reservationRepo.CreateTx(tx, res); err != nil {
			return fmt.Errorf("reservation create: %w", err)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	if res.Status == constants.ReservationWaitlisted {
		s.logger.Info(fmt.Sprintf(constants.LogTemplates["reservation_waitlist_ok"], userID, req.StationID, req.StartTime.Format("2006-01-02 15:04"), res.ID))
	} else {
		s.logger.Info(fmt.Sprintf(constants.LogTemplates["reservation_create_ok"], userID, req.StationID, req.StartTime.Format("2006-01-02 15:04")))
	}
	return res, nil
}

// Confirm 确认预约（staff/admin）。候补预约只能由取消事件自动兑现，不允许人工确认。
func (s *ReservationService) Confirm(id uint) (*model.Reservation, error) {
	err := s.db.Transaction(func(tx *gorm.DB) error {
		res, err := s.getReservationForUpdate(tx, id)
		if err != nil {
			return err
		}
		if res.Status == constants.ReservationWaitlisted {
			return util.NewAppError(constants.CodeWaitlisted, "候补预约需等待取消自动兑现，无法人工确认")
		}
		if res.Status != constants.ReservationPending && res.Status != constants.ReservationConfirmed {
			return util.NewAppError(constants.CodeReservation, "仅待确认或已确认的预约可以确认")
		}
		res.Status = constants.ReservationConfirmed
		return s.reservationRepo.UpdateTx(tx, res)
	})
	if err != nil {
		return nil, err
	}
	s.logger.Info(fmt.Sprintf(constants.LogTemplates["reservation_confirm_ok"], id))
	return s.getReservation(id)
}

// Cancel 取消预约。
// 已确认预约取消时，在同机位、与该时段重叠的候补里取最早提交且此刻不再冲突的一位自动确认，
// 一次取消至多兑现一人；重复/并发取消只会有一次生效，机位状态按剩余预约重新核对。
func (s *ReservationService) Cancel(id uint, userID uint) (*model.Reservation, error) {
	// 事务外无锁预读以确定机位；权威状态以事务内行锁复查为准。
	pre, err := s.getReservation(id)
	if err != nil {
		return nil, err
	}
	stationID := pre.StationID
	var promotedID uint
	err = s.db.Transaction(func(tx *gorm.DB) error {
		// 先锁机位，再锁预约行，与 Create/CheckIn/上机 保持相同加锁顺序。
		if _, err := s.stationService.LockForUpdate(tx, stationID); err != nil {
			return err
		}
		res, err := s.getReservationForUpdate(tx, id)
		if err != nil {
			return err
		}
		from := res.Status
		if !cancelable(from) {
			// 已取消/已完成等状态重复取消直接报错，保证一次取消只兑现一次候补。
			return util.NewAppError(constants.CodeReservation, "当前状态不可取消")
		}
		res.Status = constants.ReservationCancelled
		if err := s.reservationRepo.UpdateTx(tx, res); err != nil {
			return err
		}
		// 仅已确认预约的取消触发候补兑现；候补/已开机的取消不触发。
		if from == constants.ReservationConfirmed {
			promotedID, err = s.promoteFromWaitlist(tx, stationID, res.StartTime, res.EndTime)
			if err != nil {
				return err
			}
		}
		return reconcileStationReserved(tx, s.stationService, s.reservationRepo, stationID)
	})
	if err != nil {
		return nil, err
	}
	if promotedID > 0 {
		s.logger.Info(fmt.Sprintf(constants.LogTemplates["reservation_promote_ok"], id, promotedID, stationID))
	}
	s.logger.Info(fmt.Sprintf(constants.LogTemplates["reservation_cancel_ok"], id))
	return s.getReservation(id)
}

// promoteFromWaitlist 取消后兑现候补：同机位、与取消时段重叠，按提交先后遍历，
// 最早一位与现存占用预约不再冲突的候补转为已确认。一次取消最多兑现一人，其余候补顺位自然前移。
func (s *ReservationService) promoteFromWaitlist(tx *gorm.DB, stationID uint, start, end time.Time) (uint, error) {
	candidates, err := s.reservationRepo.ListWaitlistedOverlapping(tx, stationID, start, end)
	if err != nil {
		return 0, fmt.Errorf("reservation list waitlist: %w", err)
	}
	for i := range candidates {
		cand := &candidates[i]
		if !cand.EndTime.After(time.Now()) {
			// 候补窗口已结束，不再兑现，避免浪费本次取消仅有的一个名额。
			continue
		}
		cnt, err := s.reservationRepo.CountConflictTx(tx, cand.StationID, cand.StartTime, cand.EndTime, cand.ID)
		if err != nil {
			return 0, fmt.Errorf("reservation promote conflict check: %w", err)
		}
		if cnt > 0 {
			// 仍与其他已确认/已开机预约冲突，保留候补，继续看后面的人。
			continue
		}
		cand.Status = constants.ReservationConfirmed
		if err := s.reservationRepo.UpdateTx(tx, cand); err != nil {
			return 0, fmt.Errorf("reservation promote update: %w", err)
		}
		return cand.ID, nil
	}
	return 0, nil
}

// CheckIn 到店扫码开机：预约状态流转为 checked_in，机位置为使用中。候补未转正不能开机。
func (s *ReservationService) CheckIn(id uint) (*model.Reservation, error) {
	pre, err := s.getReservation(id)
	if err != nil {
		return nil, err
	}
	stationID := pre.StationID
	err = s.db.Transaction(func(tx *gorm.DB) error {
		station, err := s.stationService.LockForUpdate(tx, stationID)
		if err != nil {
			return err
		}
		res, err := s.getReservationForUpdate(tx, id)
		if err != nil {
			return err
		}
		if res.Status == constants.ReservationWaitlisted {
			return util.NewAppError(constants.CodeWaitlisted, "候补尚未转正，无法开机")
		}
		if res.Status != constants.ReservationConfirmed {
			return util.NewAppError(constants.CodeReservation, "仅已确认的预约可以开机")
		}
		res.Status = constants.ReservationCheckedIn
		if err := s.reservationRepo.UpdateTx(tx, res); err != nil {
			return err
		}
		station.Status = constants.StationUsing
		return tx.Save(station).Error
	})
	if err != nil {
		return nil, fmt.Errorf("reservation checkin tx: %w", err)
	}
	s.logger.Info(fmt.Sprintf(constants.LogTemplates["reservation_checkin_ok"], id))
	return s.getReservation(id)
}

// List 分页查询预约，候补中记录携带当前顺位。
func (s *ReservationService) List(query *dto.ReservationQuery) ([]model.Reservation, int64, error) {
	page := query.Page
	pageSize := query.PageSize
	if page <= 0 {
		page = constants.DefaultPage
	}
	if pageSize <= 0 {
		pageSize = constants.DefaultPageSize
	}
	return s.reservationRepo.List(page, pageSize, query.Status, query.UserID)
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

// getReservationForUpdate 事务内行锁查询预约并统一处理错误。
func (s *ReservationService) getReservationForUpdate(tx *gorm.DB, id uint) (*model.Reservation, error) {
	res, err := s.reservationRepo.FindByIDForUpdate(tx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, util.NewAppError(constants.CodeNotFound, "预约记录不存在")
		}
		return nil, fmt.Errorf("reservation find for update: %w", err)
	}
	return res, nil
}

// cancelable 允许取消的预约状态。
func cancelable(status string) bool {
	switch status {
	case constants.ReservationPending, constants.ReservationConfirmed,
		constants.ReservationWaitlisted, constants.ReservationCheckedIn:
		return true
	}
	return false
}

// reconcileStationReserved 取消/兑现后核对机位状态：
// 机位处于“已预约”且当前不存在未结束的占用预约时释放为空闲；使用中/故障等状态不动。
func reconcileStationReserved(tx *gorm.DB, svc *StationService, repo *repository.ReservationRepository, stationID uint) error {
	station, err := svc.LockForUpdate(tx, stationID)
	if err != nil {
		return err
	}
	if station.Status != constants.StationReserved {
		return nil
	}
	cnt, err := repo.CountBlockingFuture(tx, stationID, time.Now())
	if err != nil {
		return err
	}
	if cnt == 0 {
		station.Status = constants.StationIdle
		return tx.Save(station).Error
	}
	return nil
}
