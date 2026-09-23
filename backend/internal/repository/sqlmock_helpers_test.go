package repository

import (
	"github.com/DATA-DOG/go-sqlmock"
)

// sqlmockRowsCount 生成 count 查询结果行。
func sqlmockRowsCount(n int) *sqlmock.Rows {
	return sqlmock.NewRows([]string{"count"}).AddRow(n)
}
