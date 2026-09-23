-- 电竞馆上机管理系统初始表结构（GORM AutoMigrate 会自动执行，此脚本供参考与初始化）
CREATE TABLE IF NOT EXISTS users (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(64) NOT NULL UNIQUE,
    password VARCHAR(128) NOT NULL,
    nickname VARCHAR(64) DEFAULT '',
    phone VARCHAR(20) DEFAULT '',
    role VARCHAR(16) DEFAULT 'member',
    balance DECIMAL(12,2) DEFAULT 0,
    status VARCHAR(16) DEFAULT 'active',
    created_at DATETIME(3),
    updated_at DATETIME(3)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS stations (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(64) NOT NULL,
    area VARCHAR(64) NOT NULL,
    station_type VARCHAR(16) DEFAULT 'seat',
    price_per_hour DECIMAL(10,2) DEFAULT 0,
    status VARCHAR(16) DEFAULT 'idle',
    description VARCHAR(255) DEFAULT '',
    created_at DATETIME(3),
    updated_at DATETIME(3)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS time_packages (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(64) NOT NULL,
    hours DECIMAL(10,2) NOT NULL,
    price DECIMAL(10,2) NOT NULL,
    valid_days INT DEFAULT 30,
    status VARCHAR(16) DEFAULT 'active',
    created_at DATETIME(3),
    updated_at DATETIME(3)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS user_packages (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT UNSIGNED NOT NULL,
    package_id BIGINT UNSIGNED NOT NULL,
    package_name VARCHAR(64) DEFAULT '',
    total_hours DECIMAL(10,2) DEFAULT 0,
    remaining_hours DECIMAL(10,2) DEFAULT 0,
    expire_at DATETIME(3),
    status VARCHAR(16) DEFAULT 'active',
    created_at DATETIME(3),
    updated_at DATETIME(3),
    KEY idx_user_packages_user (user_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS recharges (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT UNSIGNED NOT NULL,
    amount DECIMAL(12,2) NOT NULL,
    payment_method VARCHAR(16) DEFAULT 'balance',
    operator_id BIGINT UNSIGNED DEFAULT 0,
    remark VARCHAR(255) DEFAULT '',
    created_at DATETIME(3),
    KEY idx_recharges_user (user_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS package_orders (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    order_no VARCHAR(64) NOT NULL UNIQUE,
    user_id BIGINT UNSIGNED NOT NULL,
    package_id BIGINT UNSIGNED NOT NULL,
    package_name VARCHAR(64) DEFAULT '',
    amount DECIMAL(12,2) NOT NULL,
    hours DECIMAL(10,2) DEFAULT 0,
    payment_method VARCHAR(16) DEFAULT 'balance',
    status VARCHAR(16) DEFAULT 'pending',
    paid_at DATETIME(3),
    created_at DATETIME(3),
    updated_at DATETIME(3),
    KEY idx_package_orders_user (user_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS reservations (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT UNSIGNED NOT NULL,
    station_id BIGINT UNSIGNED NOT NULL,
    start_time DATETIME(3) NOT NULL,
    end_time DATETIME(3) NOT NULL,
    status VARCHAR(16) DEFAULT 'pending',
    remark VARCHAR(255) DEFAULT '',
    created_at DATETIME(3),
    updated_at DATETIME(3),
    KEY idx_reservations_user (user_id),
    KEY idx_reservations_station (station_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS sessions (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT UNSIGNED NOT NULL,
    station_id BIGINT UNSIGNED NOT NULL,
    reservation_id BIGINT UNSIGNED DEFAULT 0,
    start_time DATETIME(3) NOT NULL,
    end_time DATETIME(3),
    duration_minutes INT DEFAULT 0,
    game_type VARCHAR(16) DEFAULT 'other',
    amount DECIMAL(12,2) DEFAULT 0,
    status VARCHAR(16) DEFAULT 'active',
    created_at DATETIME(3),
    updated_at DATETIME(3),
    KEY idx_sessions_user (user_id),
    KEY idx_sessions_station (station_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS tournaments (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(128) NOT NULL,
    game_type VARCHAR(16) DEFAULT 'lol',
    description VARCHAR(512) DEFAULT '',
    register_start DATETIME(3),
    register_end DATETIME(3),
    start_time DATETIME(3),
    max_teams INT DEFAULT 16,
    status VARCHAR(16) DEFAULT 'draft',
    created_by BIGINT UNSIGNED DEFAULT 0,
    created_at DATETIME(3),
    updated_at DATETIME(3)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS teams (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    tournament_id BIGINT UNSIGNED NOT NULL,
    name VARCHAR(64) NOT NULL,
    leader_id BIGINT UNSIGNED DEFAULT 0,
    member_count INT DEFAULT 1,
    created_at DATETIME(3),
    updated_at DATETIME(3),
    KEY idx_teams_tournament (tournament_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS registrations (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    tournament_id BIGINT UNSIGNED NOT NULL,
    user_id BIGINT UNSIGNED NOT NULL,
    team_id BIGINT UNSIGNED DEFAULT 0,
    mode VARCHAR(16) DEFAULT 'solo',
    group_no INT DEFAULT 0,
    status VARCHAR(16) DEFAULT 'pending',
    created_at DATETIME(3),
    updated_at DATETIME(3),
    KEY idx_registrations_tournament (tournament_id),
    KEY idx_registrations_user (user_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS matches (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    tournament_id BIGINT UNSIGNED NOT NULL,
    round INT DEFAULT 1,
    group_no INT DEFAULT 0,
    team_a_id BIGINT UNSIGNED DEFAULT 0,
    team_b_id BIGINT UNSIGNED DEFAULT 0,
    score_a INT DEFAULT 0,
    score_b INT DEFAULT 0,
    winner_id BIGINT UNSIGNED DEFAULT 0,
    status VARCHAR(16) DEFAULT 'pending',
    scheduled_at DATETIME(3),
    created_at DATETIME(3),
    updated_at DATETIME(3),
    KEY idx_matches_tournament (tournament_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS audit_logs (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT UNSIGNED DEFAULT 0,
    username VARCHAR(64) DEFAULT '',
    action VARCHAR(64) DEFAULT '',
    module VARCHAR(64) DEFAULT '',
    target_type VARCHAR(64) DEFAULT '',
    target_id BIGINT UNSIGNED DEFAULT 0,
    detail VARCHAR(1024) DEFAULT '',
    ip VARCHAR(64) DEFAULT '',
    created_at DATETIME(3),
    KEY idx_audit_logs_user (user_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
