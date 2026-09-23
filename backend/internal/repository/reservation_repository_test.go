package repository

import (
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"

	"github.com/esportsbar/backend/internal/model"
)

// TestReservationCountConflictExcludesWaitlisted 冲突统计只认 pending/confirmed/checked_in，候补不算占用。
func TestReservationCountConflictExcludesWaitlisted(t *testing.T) {
	gdb, mock := newMockDB(t)
	repo := NewReservationRepository(gdb)
	start := time.Date(2026, 9, 23, 18, 0, 0, 0, time.UTC)
	end := start.Add(2 * time.Hour)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM `reservations` WHERE (station_id = ? AND status IN (?,?,?)) AND (start_time < ? AND end_time > ?)")).
		WithArgs(uint(7), "pending", "confirmed", "checked_in", end, start).
		WillReturnRows(sqlmockRowsCount(0))
	cnt, err := repo.CountConflict(7, start, end, 0)
	if err != nil {
		t.Fatalf("CountConflict error: %v", err)
	}
	if cnt != 0 {
		t.Fatalf("expected 0 conflict, got %d", cnt)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectations not met: %v", err)
	}
}

// TestReservationQueuePosition 候补顺位 = 同机位时段重叠且更早提交的候补数 + 1。
func TestReservationQueuePosition(t *testing.T) {
	gdb, mock := newMockDB(t)
	repo := NewReservationRepository(gdb)
	createdAt := time.Date(2026, 9, 23, 17, 30, 0, 0, time.UTC)
	start := time.Date(2026, 9, 23, 18, 0, 0, 0, time.UTC)
	end := start.Add(2 * time.Hour)
	mock.ExpectQuery("SELECT count\\(\\*\\) FROM `reservations` WHERE").
		WillReturnRows(sqlmockRowsCount(2))
	pos, err := repo.QueuePosition(&model.Reservation{ID: 9, StationID: 3, CreatedAt: createdAt, StartTime: start, EndTime: end})
	if err != nil {
		t.Fatalf("QueuePosition error: %v", err)
	}
	if pos != 3 {
		t.Fatalf("expected position 3, got %d", pos)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectations not met: %v", err)
	}
}

// TestReservationListFillsQueuePosition 列表为候补记录回填顺位。
func TestReservationListFillsQueuePosition(t *testing.T) {
	gdb, mock := newMockDB(t)
	repo := NewReservationRepository(gdb)
	createdAt := time.Date(2026, 9, 23, 17, 30, 0, 0, time.UTC)
	start := time.Date(2026, 9, 23, 18, 0, 0, 0, time.UTC)
	end := start.Add(2 * time.Hour)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM `reservations`")).
		WillReturnRows(sqlmockRowsCount(2))
	mock.ExpectQuery("SELECT \\* FROM `reservations`").
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "station_id", "start_time", "end_time", "status", "created_at", "updated_at"}).
			AddRow(10, 1, 5, start, end, "waitlisted", createdAt, createdAt).
			AddRow(11, 2, 6, start, end, "confirmed", createdAt, createdAt))
	mock.ExpectQuery("SELECT count\\(\\*\\) FROM `reservations` WHERE").
		WillReturnRows(sqlmockRowsCount(0))
	list, total, err := repo.List(1, 10, "", 0)
	if err != nil {
		t.Fatalf("List error: %v", err)
	}
	if total != 2 || len(list) != 2 {
		t.Fatalf("unexpected total=%d len=%d", total, len(list))
	}
	if list[0].Status != "waitlisted" || list[0].QueuePosition != 1 {
		t.Fatalf("waitlisted queue position not filled: %+v", list[0])
	}
	if list[1].QueuePosition != 0 {
		t.Fatalf("confirmed reservation should have no queue position: %+v", list[1])
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectations not met: %v", err)
	}
}
