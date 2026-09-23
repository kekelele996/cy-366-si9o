package database

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/esportsbar/backend/internal/config"
)

// ConnectRedis 建立 Redis 连接并做连通性检查。
func ConnectRedis(cfg *config.Config, logger *slog.Logger) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisHost + ":" + cfg.RedisPort,
		Password: cfg.RedisPassword,
		DB:       0,
	})
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := client.Ping(ctx).Err(); err != nil {
		logger.Warn("redis ping failed, fallback to in-memory rate limit", "err", err)
		return client
	}
	logger.Info("redis connected", "addr", cfg.RedisHost+":"+cfg.RedisPort)
	return client
}

var _ = fmt.Sprintf
