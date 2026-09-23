package repository

import (
	"errors"
	"time"

	"gorm.io/gorm"

	"github.com/esportsbar/backend/internal/constants"
	"github.com/esportsbar/backend/internal/model"
)

// ReservationRepository 预约仓储。
type ReservationRepository struct {
	db *gorm.DB
}

// NewReservationRepository 构造预约仓储。
func NewReservationRepository(db *gorm.DB) *ReservationRepository {
	return &ReservationRepository{db: db}
}

// Create 创建预约。
func (r *ReservationRepository) Create(res *model.Reservation) error {
	return r.db.Create(res).Error
}

// CreateTx 事务内创建预约。
func (r *ReservationRepository) CreateTx(tx *gorm.DB, res *model.Reservation) error {
	return tx.Create(res).Error
}

// FindByID 查询预约。
func (r *ReservationRepository) FindByID(id uint) (*model.Reservation, error) {
	return r.findByID(r.db, id)
}

// FindByIDForUpdate 事务内行锁查询预约（取消/开机等并发状态流转使用）。
func (r *ReservationRepository) FindByIDForUpdate(tx *gorm.DB, id uint) (*model.Reservation, error) {
	var res model.Reservation
	err := tx.Clauses(clauseLocking()).First(&res, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	return &res, err
}

func (r *ReservationRepository) findByID(q *gorm.DB, id uint) (*model.Reservation, error) {
	var res model.Reservation
	err := q.First(&res, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	return &res, err
}

// Update 更新预约。
func (r *ReservationRepository) Update(res *model.Reservation) error {
	return r.db.Save(res).Error
}

// UpdateTx 事务内更新预约。
func (r *ReservationRepository) UpdateTx(tx *gorm.DB, res *model.Reservation) error {
	return tx.Save(res).Error
}

// List 分页查询预约，并为候补中记录回填当前候补顺位。
func (r *ReservationRepository) List(page, pageSize int, status string, userID uint) ([]model.Reservation, int64, error) {
	var list []model.Reservation
	var total int64
	query := r.db.Model(&model.Reservation{})
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if userID > 0 {
		query = query.Where("user_id = ?", userID)
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := query.Order("start_time DESC, id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	for i := range list {
		if list[i].Status == constants.ReservationWaitlisted {
			pos, err := r.QueuePosition(&list[i])
			if err != nil {
				return nil, 0, err
			}
			list[i].QueuePosition = pos
		}
	}
	return list, total, nil
}

// CountConflict 统计机位在时段内的冲突预约数（占用机位的状态：待确认/已确认/已开机）。
func (r *ReservationRepository) CountConflict(stationID uint, start, end time.Time, excludeID uint) (int64, error) {
	return r.CountConflictTx(r.db, stationID, start, end, excludeID)
}

// CountConflictTx 事务内统计机位在时段内的冲突预约数。
func (r *ReservationRepository) CountConflictTx(tx *gorm.DB, stationID uint, start, end time.Time, excludeID uint) (int64, error) {
	var cnt int64
	query := tx.Model(&model.Reservation{}).
		Where("station_id = ? AND status IN ?", stationID, blockingStatuses).
		Where("start_time < ? AND end_time > ?", end, start)
	if excludeID > 0 {
		query = query.Where("id <> ?", excludeID)
	}
	err := query.Count(&cnt).Error
	return cnt, err
}

// ListWaitlistedOverlapping 查出同机位、与给定时段重叠且仍在候补的预约，
// 按提交先后（created_at, id）升序，供取消后顺位兑现使用。
func (r *ReservationRepository) ListWaitlistedOverlapping(tx *gorm.DB, stationID uint, start, end time.Time) ([]model.Reservation, error) {
	var list []model.Reservation
	err := tx.Clauses(clauseLocking()).
		Where("station_id = ? AND status = ?", stationID, constants.ReservationWaitlisted).
		Where("start_time < ? AND end_time > ?", end, start).
		Order("created_at ASC, id ASC").
		Find(&list).Error
	return list, err
}

// QueuePosition 返回候补预约的当前顺位：同机位、时段重叠、更早提交的候补数 + 1。
func (r *ReservationRepository) QueuePosition(res *model.Reservation) (int, error) {
	return r.QueuePositionTx(r.db, res)
}

// QueuePositionTx 事务内计算候补顺位。
func (r *ReservationRepository) QueuePositionTx(tx *gorm.DB, res *model.Reservation) (int, error) {
	var cnt int64
	err := tx.Model(&model.Reservation{}).
		Where("station_id = ? AND status = ?", res.StationID, constants.ReservationWaitlisted).
		Where("start_time < ? AND end_time > ?", res.EndTime, res.StartTime).
		Where("created_at < ? OR (created_at = ? AND id < ?)", res.CreatedAt, res.CreatedAt, res.ID).
		Count(&cnt).Error
	return int(cnt) + 1, err
}

// blockingStatuses 占用机位、会造成时段冲突的预约状态。候补不占机位，不在其中。
var blockingStatuses = []string{
	constants.ReservationPending,
	constants.ReservationConfirmed,
	constants.ReservationCheckedIn,
}

// CountBlockingFuture 统计机位在 from 时刻之后仍占用机位（已确认/待确认/已开机且未结束）的预约数，
// 取消/兑现后用于核对机位是否还能保持“已预约”。
func (r *ReservationRepository) CountBlockingFuture(tx *gorm.DB, stationID uint, from time.Time) (int64, error) {
	var cnt int64
	err := tx.Model(&model.Reservation{}).
		Where("station_id = ? AND status IN ?", stationID, blockingStatuses).
		Where("end_time > ?", from).
		Count(&cnt).Error
	return cnt, err
}
