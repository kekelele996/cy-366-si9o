package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/esportsbar/backend/internal/constants"
	"github.com/esportsbar/backend/internal/util"
)

// Auth JWT 认证中间件。
func Auth(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" || !strings.HasPrefix(header, "Bearer ") {
			util.AbortJSON(c, 401, constants.CodeUnauthorized, "未登录或登录已过期")
			return
		}
		tokenStr := strings.TrimPrefix(header, "Bearer ")
		claims, err := util.ParseToken(jwtSecret, tokenStr)
		if err != nil {
			util.AbortJSON(c, 401, constants.CodeUnauthorized, "登录凭证无效，请重新登录")
			return
		}
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)
		c.Next()
	}
}
