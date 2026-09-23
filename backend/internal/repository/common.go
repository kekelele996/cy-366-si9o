package repository

import "gorm.io/gorm/clause"

// clauseLocking 返回 MySQL 行锁子句，并发扣款/状态流转时使用 SELECT ... FOR UPDATE。
func clauseLocking() clause.Locking {
	return clause.Locking{Strength: "UPDATE"}
}
