package middleware

import (
	"github.com/gin-gonic/gin"

	"github.com/esportsbar/backend/internal/util"
)

// RequestID 为每个请求生成追踪 ID，并注入响应头。
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.GetHeader("X-Request-ID")
		if id == "" {
			id = util.NewRequestID()
		}
		c.Set("request_id", id)
		c.Header("X-Request-ID", id)
		util.WithRequestID(c.Request.Context(), id)
		c.Next()
	}
}
