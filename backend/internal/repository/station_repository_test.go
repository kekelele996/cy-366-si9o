package repository

import (
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/gorm"
)

func TestStationRepositoryList(t *testing.T) {
	gdb, mock := newMockDB(t)
	repo := NewStationRepository(gdb)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM `stations`")).
		WillReturnRows(sqlmockRowsCount(3))
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `stations` ORDER BY area, id LIMIT ?")).
		WithArgs(10).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "status"}).
			AddRow(1, "A区-01", "idle").
			AddRow(2, "A区-02", "using").
			AddRow(3, "B区-01", "fault"))
	list, total, err := repo.List(1, 10, "", "")
	if err != nil {
		t.Fatalf("List error: %v", err)
	}
	if total != 3 || len(list) != 3 {
		t.Fatalf("unexpected result total=%d len=%d", total, len(list))
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectations not met: %v", err)
	}
}

func TestStationRepositoryFindByID(t *testing.T) {
	gdb, mock := newMockDB(t)
	repo := NewStationRepository(gdb)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `stations` WHERE `stations`.`id` = ? ORDER BY `stations`.`id` LIMIT ?")).
		WithArgs(5, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "status"}).AddRow(5, "包厢区-01", "idle"))
	s, err := repo.FindByID(5)
	if err != nil {
		t.Fatalf("FindByID error: %v", err)
	}
	if s.ID != 5 {
		t.Fatalf("unexpected station: %+v", s)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectations not met: %v", err)
	}
}

var _ = gorm.ErrRecordNotFound
