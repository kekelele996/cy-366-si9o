package service

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/esportsbar/backend/internal/constants"
	"github.com/esportsbar/backend/internal/dto"
	"github.com/esportsbar/backend/internal/model"
	"github.com/esportsbar/backend/internal/repository"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// 候补端到端集成测试：使用内存 SQLite 复刻 GORM 模型与真实事务。
// FOR UPDATE 在 SQLite 上退化为普通读，行的正确性靠事务内条件复查保证；加锁顺序与 MySQL 版本一致。
var waitlistDBSN uint64

func newWaitlistTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	dsn := fmt.Sprintf("file:mem%d?mode=memory&cache=shared", atomic.AddUint64(&waitlistDBSN, 1))
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	// 单连接串行化写事务，规避 SQLite 库级锁；不影响对“事务内状态复查、取消幂等”的验证。
	sqlDB, _ := db.DB()
	sqlDB.SetMaxOpenConns(1)
	if err := db.AutoMigrate(&model.Station{}, &model.Reservation{}); err != nil {
		t.Fatalf("automigrate: %v", err)
	}
	return db
}

func newWaitlistService(t *testing.T) (*gorm.DB, *ReservationService) {
	db := newWaitlistTestDB(t)
	station := &model.Station{Name: "A区-01", Area: "A区", StationType: constants.StationIdle, Status: constants.StationIdle, PricePerHour: 8}
	if err := db.Create(station).Error; err != nil {
		t.Fatalf("seed station: %v", err)
	}
	resRepo := repository.NewReservationRepository(db)
	stationRepo := repository.NewStationRepository(db)
	logger := newTestLogger()
	stationSvc := NewStationService(stationRepo, logger)
	return db, NewReservationService(resRepo, stationSvc, db, logger)
}

func mkReqAt(start time.Time, stationID uint) *dto.CreateReservationReq {
	return &dto.CreateReservationReq{StationID: stationID, StartTime: start, EndTime: start.Add(2 * time.Hour)}
}

func mkReqWindow(start time.Time, dur time.Duration, stationID uint) *dto.CreateReservationReq {
	return &dto.CreateReservationReq{StationID: stationID, StartTime: start, EndTime: start.Add(dur)}
}

// TestWaitlistCreateAndQueuePosition 时段冲突时进入候补并带顺位，候补不占机位。
func TestWaitlistCreateAndQueuePosition(t *testing.T) {
	db, svc := newWaitlistService(t)
	base := time.Date(2026, 9, 25, 19, 0, 0, 0, time.UTC)

	r1, err := svc.Create(101, mkReqAt(base, 1))
	if err != nil {
		t.Fatalf("create r1: %v", err)
	}
	if r1.Status != constants.ReservationConfirmed {
		t.Fatalf("r1 should be confirmed, got %s", r1.Status)
	}

	// 与 r1 完全重叠：候补第 1 位
	w1, err := svc.Create(102, mkReqAt(base.Add(10*time.Minute), 1))
	if err != nil {
		t.Fatalf("create w1: %v", err)
	}
	if w1.Status != constants.ReservationWaitlisted {
		t.Fatalf("w1 should be waitlisted, got %s", w1.Status)
	}
	// 部分重叠（19:30-21:30 与 19:00-21:00 重叠）：候补第 2 位
	w2, err := svc.Create(103, mkReqAt(base.Add(30*time.Minute), 1))
	if err != nil || w2.Status != constants.ReservationWaitlisted {
		t.Fatalf("w2 should be waitlisted, err=%v", err)
	}
	// 完全不重叠的晚场（22:00-24:00）：直接确认
	r2, err := svc.Create(104, mkReqAt(base.Add(3*time.Hour), 1))
	if err != nil || r2.Status != constants.ReservationConfirmed {
		t.Fatalf("r2 should be confirmed, err=%v", err)
	}

	// 候补不占机位：机位仍为 reserved（由 r1/r2 占用）。
	var st model.Station
	db.First(&st, 1)
	if st.Status != constants.StationReserved {
		t.Fatalf("station should stay reserved, got %s", st.Status)
	}

	list, total, err := svc.List(&dto.ReservationQuery{Page: 1, PageSize: 50})
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if total != 4 {
		t.Fatalf("total = %d, want 4", total)
	}
	posByID := map[uint]int{}
	for _, r := range list {
		posByID[r.ID] = r.QueuePosition
	}
	if posByID[w1.ID] != 1 {
		t.Fatalf("w1 position = %d, want 1", posByID[w1.ID])
	}
	if posByID[w2.ID] != 2 {
		t.Fatalf("w2 position = %d, want 2", posByID[w2.ID])
	}
}

// TestCancelPromotesEarliestNonConflictingWaitlist 取消已确认预约：
// 同机位重叠候补里最早提交且不再冲突的一位转正，其余顺位前移，一次只兑现一人。
func TestCancelPromotesEarliestNonConflictingWaitlist(t *testing.T) {
	db, svc := newWaitlistService(t)
	base := time.Date(2026, 9, 25, 19, 0, 0, 0, time.UTC)

	r1, _ := svc.Create(101, mkReqAt(base, 1)) // confirmed 19:00-21:00
	w1, _ := svc.Create(102, mkReqAt(base.Add(10*time.Minute), 1))
	w2, _ := svc.Create(103, mkReqAt(base.Add(30*time.Minute), 1))
	w3, _ := svc.Create(105, mkReqAt(base.Add(60*time.Minute), 1))

	cancelled, err := svc.Cancel(r1.ID, 101)
	if err != nil {
		t.Fatalf("cancel r1: %v", err)
	}
	if cancelled.Status != constants.ReservationCancelled {
		t.Fatalf("r1 status = %s, want cancelled", cancelled.Status)
	}

	assertStatus := func(id uint, want, name string) {
		var r model.Reservation
		db.First(&r, id)
		if r.Status != want {
			t.Fatalf("%s status = %s, want %s", name, r.Status, want)
		}
	}
	// w1 最早提交且 r1 取消后不再冲突 -> 转正；w2/w3 仍候补（与 w1 冲突）
	assertStatus(w1.ID, constants.ReservationConfirmed, "w1")
	assertStatus(w2.ID, constants.ReservationWaitlisted, "w2")
	assertStatus(w3.ID, constants.ReservationWaitlisted, "w3")

	// 顺位前移：w2 现在是第 1 位（w1 已离开队列）
	list, _, err := svc.List(&dto.ReservationQuery{Page: 1, PageSize: 50})
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	for _, r := range list {
		if r.ID == w2.ID && r.QueuePosition != 1 {
			t.Fatalf("w2 position after promotion = %d, want 1", r.QueuePosition)
		}
	}

	// 再次取消已转正的 w1：一次取消仍只兑现一人（w2 转正，w3 留队）
	if _, err := svc.Cancel(w1.ID, 102); err != nil {
		t.Fatalf("cancel w1: %v", err)
	}
	assertStatus(w2.ID, constants.ReservationConfirmed, "w2 after second cancel")
	assertStatus(w3.ID, constants.ReservationWaitlisted, "w3 stays waitlisted")

	// 机位仍被 w2 占用：reserved，不能被误置 idle
	var st model.Station
	db.First(&st, 1)
	if st.Status != constants.StationReserved {
		t.Fatalf("station status = %s, want reserved", st.Status)
	}
}

// TestSkipConflictingCandidatePromotesLaterOne 取消后最早候补仍与另一已确认预约冲突时，
// 顺延兑现后面不再冲突的一位。
func TestSkipConflictingCandidatePromotesLaterOne(t *testing.T) {
	db, svc := newWaitlistService(t)
	base := time.Date(2026, 9, 25, 19, 0, 0, 0, time.UTC)

	// r1 19:00-20:00（将被取消）；r2 20:30-22:00（另一已确认，不重叠）
	r1, _ := svc.Create(101, mkReqWindow(base, time.Hour, 1))
	r2, _ := svc.Create(102, mkReqWindow(base.Add(90*time.Minute), 90*time.Minute, 1))
	if r2.Status != constants.ReservationConfirmed {
		t.Fatalf("seed r2 should be confirmed (non-overlapping), got %s", r2.Status)
	}
	// w1 先提交：19:30-21:00，与 r1、r2 都重叠。r1 取消后仍与 r2 冲突 -> 跳过
	w1, _ := svc.Create(103, mkReqWindow(base.Add(30*time.Minute), 90*time.Minute, 1))
	if w1.Status != constants.ReservationWaitlisted {
		t.Fatalf("w1 should be waitlisted, got %s", w1.Status)
	}
	// w2 后提交：19:10-20:00，与 r1 重叠、与 r2 不重叠。r1 取消后可兑现
	w2, _ := svc.Create(104, mkReqWindow(base.Add(10*time.Minute), 50*time.Minute, 1))
	if w2.Status != constants.ReservationWaitlisted {
		t.Fatalf("w2 should be waitlisted, got %s", w2.Status)
	}

	if _, err := svc.Cancel(r1.ID, 101); err != nil {
		t.Fatalf("cancel r1: %v", err)
	}
	assertStatus := func(id uint, want, name string) {
		var r model.Reservation
		db.First(&r, id)
		if r.Status != want {
			t.Fatalf("%s status = %s, want %s", name, r.Status, want)
		}
	}
	assertStatus(w1.ID, constants.ReservationWaitlisted, "w1 still conflicts with r2")
	assertStatus(w2.ID, constants.ReservationConfirmed, "w2 promoted despite later submit")
}

// TestStationReleasedWhenNoBlockingReservation 取消后无任何占用预约时，机位释放为空闲。
func TestStationReleasedWhenNoBlockingReservation(t *testing.T) {
	db, svc := newWaitlistService(t)
	base := time.Date(2026, 9, 28, 19, 0, 0, 0, time.UTC)
	r1, _ := svc.Create(101, mkReqAt(base, 1))
	// 候补窗口完全在 r1 内，r1 取消后转正；再把它也取消，此时无占用 -> idle
	w1, _ := svc.Create(102, mkReqAt(base, 1))

	if _, err := svc.Cancel(r1.ID, 101); err != nil {
		t.Fatalf("cancel r1: %v", err)
	}
	if _, err := svc.Cancel(w1.ID, 102); err != nil {
		t.Fatalf("cancel w1: %v", err)
	}
	var st model.Station
	db.First(&st, 1)
	if st.Status != constants.StationIdle {
		t.Fatalf("station should be released to idle, got %s", st.Status)
	}
}

// TestDoubleCancelRejected 重复取消同一预约第二次必须失败，不会二次兑现。
func TestDoubleCancelRejected(t *testing.T) {
	db, svc := newWaitlistService(t)
	base := time.Date(2026, 9, 26, 19, 0, 0, 0, time.UTC)
	r1, _ := svc.Create(101, mkReqAt(base, 1))
	w1, _ := svc.Create(102, mkReqAt(base, 1))

	if _, err := svc.Cancel(r1.ID, 101); err != nil {
		t.Fatalf("first cancel: %v", err)
	}
	if _, err := svc.Cancel(r1.ID, 101); err == nil {
		t.Fatal("second cancel of same reservation must be rejected")
	}
	var w model.Reservation
	db.First(&w, w1.ID)
	if w.Status != constants.ReservationConfirmed {
		t.Fatalf("w1 should be promoted exactly once, status=%s", w.Status)
	}
}

// TestWaitlistedCannotConfirmOrCheckIn 候补不能被人工确认或开机。
func TestWaitlistedCannotConfirmOrCheckIn(t *testing.T) {
	_, svc := newWaitlistService(t)
	base := time.Date(2026, 9, 26, 19, 0, 0, 0, time.UTC)
	r1, _ := svc.Create(101, mkReqAt(base, 1))
	w1, _ := svc.Create(102, mkReqAt(base, 1))
	_ = r1

	if _, err := svc.Confirm(w1.ID); err == nil {
		t.Fatal("confirm on waitlisted reservation must be rejected")
	}
	if _, err := svc.CheckIn(w1.ID); err == nil {
		t.Fatal("checkin on waitlisted reservation must be rejected")
	}
}

// TestConcurrentCancelPromotesOnce 并发取消同一预约：恰好一个成功，候补只被兑现一次。
func TestConcurrentCancelPromotesOnce(t *testing.T) {
	db, svc := newWaitlistService(t)
	base := time.Date(2026, 9, 27, 19, 0, 0, 0, time.UTC)
	r1, _ := svc.Create(101, mkReqAt(base, 1))
	ws := make([]*model.Reservation, 3)
	for i := range ws {
		ws[i], _ = svc.Create(uint(200+i), mkReqAt(base.Add(time.Duration(5+i*5)*time.Minute), 1))
	}

	var wg sync.WaitGroup
	results := make(chan error, 8)
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := svc.Cancel(r1.ID, 101)
			results <- err
		}()
	}
	wg.Wait()
	close(results)
	ok, fail := 0, 0
	for err := range results {
		if err == nil {
			ok++
		} else {
			fail++
		}
	}
	if ok != 1 || fail != 7 {
		t.Fatalf("concurrent cancel: success=%d fail=%d, want 1/7", ok, fail)
	}

	var confirmed, waitlisted int64
	db.Model(&model.Reservation{}).Where("status = ?", constants.ReservationConfirmed).Count(&confirmed)
	db.Model(&model.Reservation{}).Where("status = ?", constants.ReservationWaitlisted).Count(&waitlisted)
	// r1 已取消；3 个候补中恰好 1 人转正
	if confirmed != 1 {
		t.Fatalf("confirmed count = %d, want 1 (w1 promoted once)", confirmed)
	}
	if waitlisted != 2 {
		t.Fatalf("waitlisted count = %d, want 2", waitlisted)
	}
}
