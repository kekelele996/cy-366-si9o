// Package database 数据库连接与初始化。
package database

import (
	"fmt"
	"log/slog"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/esportsbar/backend/internal/config"
	"github.com/esportsbar/backend/internal/model"
	"github.com/esportsbar/backend/internal/util"
)

// Connect 建立 MySQL 连接并自动迁移。
func Connect(cfg *config.Config) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(cfg.DSN()), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		return nil, fmt.Errorf("database connect: %w", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("database pool: %w", err)
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	if err := migrate(db); err != nil {
		return nil, fmt.Errorf("database migrate: %w", err)
	}
	return db, nil
}

// migrate 自动建表。
func migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&model.User{},
		&model.Station{},
		&model.TimePackage{},
		&model.UserPackage{},
		&model.Recharge{},
		&model.PackageOrder{},
		&model.Reservation{},
		&model.Session{},
		&model.Tournament{},
		&model.Team{},
		&model.Registration{},
		&model.Match{},
		&model.AuditLog{},
	)
}

// Seed 初始化种子数据（管理员、示例机位、时长包）。
func Seed(db *gorm.DB, logger *slog.Logger) error {
	var count int64
	if err := db.Model(&model.User{}).Count(&count).Error; err != nil {
		return fmt.Errorf("seed count user: %w", err)
	}
	if count == 0 {
		hash, err := util.HashPassword("admin123456")
		if err != nil {
			return fmt.Errorf("seed hash admin: %w", err)
		}
		admin := &model.User{
			Username: "admin",
			Password: hash,
			Nickname: "管理员",
			Role:     "admin",
			Status:   "active",
			Balance:  0,
		}
		if err := db.Create(admin).Error; err != nil {
			return fmt.Errorf("seed create admin: %w", err)
		}
		hashMember, err := util.HashPassword("member123456")
		if err != nil {
			return fmt.Errorf("seed hash member: %w", err)
		}
		member := &model.User{
			Username: "member",
			Password: hashMember,
			Nickname: "体验会员",
			Role:     "member",
			Status:   "active",
			Balance:  100,
		}
		if err := db.Create(member).Error; err != nil {
			return fmt.Errorf("seed create member: %w", err)
		}
		logger.Info("seed admin/member created")
	}

	var stationCount int64
	db.Model(&model.Station{}).Count(&stationCount)
	if stationCount == 0 {
		areas := []string{"A区", "B区", "包厢区"}
		for i := 1; i <= 12; i++ {
			area := areas[(i-1)%3]
			stationType := "seat"
			if area == "包厢区" {
				stationType = "box"
			}
			station := &model.Station{
				Name:         fmt.Sprintf("%s-%02d", area, i),
				Area:         area,
				StationType:  stationType,
				PricePerHour: 8,
				Status:       "idle",
				Description:  "电竞专用机位",
			}
			if stationType == "box" {
				station.PricePerHour = 20
			}
			if err := db.Create(station).Error; err != nil {
				return fmt.Errorf("seed create station: %w", err)
			}
		}
		logger.Info("seed stations created")
	}

	var pkgCount int64
	db.Model(&model.TimePackage{}).Count(&pkgCount)
	if pkgCount == 0 {
		pkgs := []model.TimePackage{
			{Name: "10小时时长包", Hours: 10, Price: 60, ValidDays: 30, Status: "active"},
			{Name: "30小时时长包", Hours: 30, Price: 150, ValidDays: 60, Status: "active"},
			{Name: "电竞月卡", Hours: 60, Price: 288, ValidDays: 30, Status: "active"},
		}
		if err := db.Create(&pkgs).Error; err != nil {
			return fmt.Errorf("seed create packages: %w", err)
		}
		logger.Info("seed time packages created")
	}
	return nil
}
