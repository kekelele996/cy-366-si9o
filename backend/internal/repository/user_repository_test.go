package repository

import (
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func newMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	t.Helper()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock new error: %v", err)
	}
	gdb, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	if err != nil {
		t.Fatalf("gorm open error: %v", err)
	}
	return gdb, mock
}

func TestUserRepositoryFindByUsername(t *testing.T) {
	gdb, mock := newMockDB(t)
	repo := NewUserRepository(gdb)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE username = ? ORDER BY `users`.`id` LIMIT ?")).
		WithArgs("admin", 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "role"}).AddRow(1, "admin", "admin"))
	u, err := repo.FindByUsername("admin")
	if err != nil {
		t.Fatalf("FindByUsername error: %v", err)
	}
	if u.ID != 1 || u.Username != "admin" {
		t.Fatalf("unexpected user: %+v", u)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectations not met: %v", err)
	}
}

func TestUserRepositoryFindByUsernameNotFound(t *testing.T) {
	gdb, mock := newMockDB(t)
	repo := NewUserRepository(gdb)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE username = ? ORDER BY `users`.`id` LIMIT ?")).
		WithArgs("nobody", 1).
		WillReturnRows(sqlmock.NewRows([]string{"id"}))
	_, err := repo.FindByUsername("nobody")
	if err != ErrNotFound {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}
