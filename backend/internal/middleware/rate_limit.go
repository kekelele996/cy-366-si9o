package middleware

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"

	"github.com/esportsbar/backend/internal/util"
)

// 简易限流：每 IP 每分钟最多 perMinute 次请求。
// 优先使用 Redis（INCR+EXPIRE 原子计数），Redis 不可用时回退到进程内计数。
type rateBucket struct {
	mu    sync.Mutex
	count map[string]*bucketItem
}

type bucketItem struct {
	count   int
	resetAt time.Time
}

var fallbackBucket = &rateBucket{count: make(map[string]*bucketItem)}

// RateLimit 限流中间件。
func RateLimit(perMinute int, rdb *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		if rdb != nil {
			ctx := context.Background()
			key := "ratelimit:" + ip
			n, err := rdb.Incr(ctx, key).Result()
			if err == nil {
				if n == 1 {
					rdb.Expire(ctx, key, time.Minute)
				}
				if n > int64(perMinute) {
					util.AbortJSON(c, http.StatusTooManyRequests, 42900, "请求过于频繁，请稍后再试")
					return
				}
				c.Next()
				return
			}
		}
		// 回退：进程内计数
		fallbackBucket.mu.Lock()
		now := time.Now()
		item, ok := fallbackBucket.count[ip]
		if !ok || now.After(item.resetAt) {
			item = &bucketItem{count: 0, resetAt: now.Add(time.Minute)}
			fallbackBucket.count[ip] = item
		}
		item.count++
		fallbackBucket.mu.Unlock()
		if item.count > perMinute {
			util.AbortJSON(c, http.StatusTooManyRequests, 42900, "请求过于频繁，请稍后再试")
			return
		}
		c.Next()
	}
}
