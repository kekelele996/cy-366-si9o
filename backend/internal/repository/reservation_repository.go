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

// CreateTx 在指定事务/连接上创建预约。
func (r *ReservationRepository) CreateTx(tx *gorm.DB, res *model.Reservation) error {
	return tx.Create(res).Error
}

// FindByID 查询预约。
func (r *ReservationRepository) FindByID(id uint) (*model.Reservation, error) {
	return r.FindByIDTx(r.db, id)
}

// FindByIDTx 在指定事务/连接上查询预约。
func (r *ReservationRepository) FindByIDTx(tx *gorm.DB, id uint) (*model.Reservation, error) {
	var res model.Reservation
	err := tx.First(&res, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	return &res, err
}

// FindByIDForUpdateTx 行锁查询预约，事务内并发状态流转使用 SELECT ... FOR UPDATE。
func (r *ReservationRepository) FindByIDForUpdateTx(tx *gorm.DB, id uint) (*model.Reservation, error) {
	var res model.Reservation
	err := tx.Clauses(clauseLocking()).First(&res, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	return &res, err
}

// UpdateTx 在指定事务/连接上更新预约。
func (r *ReservationRepository) UpdateTx(tx *gorm.DB, res *model.Reservation) error {
	return tx.Save(res).Error
}

// UpdateStatusTx 条件更新预约状态：仅当当前状态属于 fromStatuses 时生效，
// 返回 rowsAffected，用于取消等幂等操作，防止重复/并发取消重复兑现候补。
func (r *ReservationRepository) UpdateStatusTx(tx *gorm.DB, id uint, fromStatuses []string, toStatus string) (int64, error) {
	res := tx.Model(&model.Reservation{}).
		Where("id = ? AND status IN ?", id, fromStatuses).
		Update("status", toStatus)
	return res.RowsAffected, res.Error
}

// List 分页查询预约。
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
	err := query.Order("start_time DESC, id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&list).Error
	return list, total, err
}

// CountConflict 统计机位在时段内的占用冲突数（含候补，排队不允许插队）。
func (r *ReservationRepository) CountConflict(stationID uint, start, end time.Time, excludeID uint) (int64, error) {
	return r.CountConflictTx(r.db, stationID, start, end, excludeID)
}

// CountConflictTx 在指定事务/连接上统计机位在时段内的占用冲突数（含候补）。
func (r *ReservationRepository) CountConflictTx(tx *gorm.DB, stationID uint, start, end time.Time, excludeID uint) (int64, error) {
	var cnt int64
	query := tx.Model(&model.Reservation{}).
		Where("station_id = ? AND status IN ?", stationID, blockingStatuses()).
		Where("start_time < ? AND end_time > ?", end, start)
	if excludeID > 0 {
		query = query.Where("id <> ?", excludeID)
	}
	err := query.Count(&cnt).Error
	return cnt, err
}

// CountActiveConflictTx 统计机位在时段内“有效预约（不含候补）”的冲突数，
// 用于候补兑现时判断最早候补是否已不再与其他有效预约冲突。
func (r *ReservationRepository) CountActiveConflictTx(tx *gorm.DB, stationID uint, start, end time.Time, excludeID uint) (int64, error) {
	var cnt int64
	query := tx.Model(&model.Reservation{}).
		Where("station_id = ? AND status IN ?", stationID, activeStatuses()).
		Where("start_time < ? AND end_time > ?", end, start)
	if excludeID > 0 {
		query = query.Where("id <> ?", excludeID)
	}
	err := query.Count(&cnt).Error
	return cnt, err
}

// ListWaitlistedByStationTx 行锁读取机位候补队列，按提交顺序（created_at, id）升序。
func (r *ReservationRepository) ListWaitlistedByStationTx(tx *gorm.DB, stationID uint) ([]model.Reservation, error) {
	var list []model.Reservation
	err := tx.Clauses(clauseLocking()).
		Where("station_id = ? AND status = ?", stationID, constants.ReservationWaitlisted).
		Order("created_at ASC, id ASC").
		Find(&list).Error
	return list, err
}

// CountActiveByStationTx 统计机位当前的有效预约数（pending/confirmed/checked_in），
// 用于取消后决定机位是否可以从 reserved 释放回 idle。
func (r *ReservationRepository) CountActiveByStationTx(tx *gorm.DB, stationID uint) (int64, error) {
	var cnt int64
	err := tx.Model(&model.Reservation{}).
		Where("station_id = ? AND status IN ?", stationID, activeStatuses()).
		Count(&cnt).Error
	return cnt, err
}

// ListWaitlistedAll 读取全部候补预约，用于列表批量回填候补顺位。
func (r *ReservationRepository) ListWaitlistedAll(db *gorm.DB) ([]model.Reservation, error) {
	var list []model.Reservation
	err := db.
		Where("status = ?", constants.ReservationWaitlisted).
		Order("station_id, created_at ASC, id ASC").
		Find(&list).Error
	return list, err
}

// blockingStatuses 阻止新建预约的状态：有效预约与候补（候补占队，不允许后来者插队）。
func blockingStatuses() []string {
	return []string{
		constants.ReservationPending,
		constants.ReservationWaitlisted,
		constants.ReservationConfirmed,
		constants.ReservationCheckedIn,
	}
}

// activeStatuses 真正占用机位时段的有效预约状态（候补未占用时段，不在其中）。
func activeStatuses() []string {
	return []string{
		constants.ReservationPending,
		constants.ReservationConfirmed,
		constants.ReservationCheckedIn,
	}
}
