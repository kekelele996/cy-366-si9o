-- 电竞馆上机管理系统数据库初始化脚本
-- MySQL 8.0；容器启动时自动执行（幂等）

CREATE DATABASE IF NOT EXISTS esportsbar_db
  DEFAULT CHARACTER SET utf8mb4
  DEFAULT COLLATE utf8mb4_unicode_ci;

CREATE USER IF NOT EXISTS 'esportsbar_user'@'%' IDENTIFIED BY 'esportsbar_pwd';
GRANT ALL PRIVILEGES ON esportsbar_db.* TO 'esportsbar_user'@'%';
FLUSH PRIVILEGES;

-- 具体表结构由后端 GORM AutoMigrate 自动创建（backend/internal/database/database.go），
-- 也可见 backend/migrations/001_init.sql 的参考建表语句。
