package middleware

import (
	"github.com/gin-gonic/gin"

	"github.com/esportsbar/backend/internal/model"
	"github.com/esportsbar/backend/internal/service"
)

// Audit 操作审计中间件：为写操作记录审计日志。
func Audit(auditService *service.AuditService) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		method := c.Request.Method
		if method == "GET" || method == "OPTIONS" || method == "HEAD" {
			return
		}
		userID, _ := c.Get("user_id")
		username, _ := c.Get("username")
		uid, _ := userID.(uint)
		uname, _ := username.(string)
		if uid == 0 {
			return
		}
		entry := &model.AuditLog{
			UserID:   uid,
			Username: uname,
			Action:   c.Request.Method + " " + c.FullPath(),
			Module:   moduleFromPath(c.FullPath()),
			Detail:   c.Request.URL.String(),
			IP:       c.ClientIP(),
		}
		_ = auditService.Record(entry)
	}
}

// moduleFromPath 从路由路径提取模块名。
func moduleFromPath(path string) string {
	// /api/v1/users/:id -> users
	parts := splitPath(path)
	for i, p := range parts {
		if p == "api" && i+2 < len(parts) {
			return parts[i+2]
		}
	}
	return "unknown"
}

// splitPath 简易路径切分。
func splitPath(path string) []string {
	var parts []string
	cur := ""
	for _, ch := range path {
		if ch == '/' {
			if cur != "" {
				parts = append(parts, cur)
				cur = ""
			}
			continue
		}
		cur += string(ch)
	}
	if cur != "" {
		parts = append(parts, cur)
	}
	return parts
}
