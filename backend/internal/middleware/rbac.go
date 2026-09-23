package middleware

import (
	"github.com/gin-gonic/gin"

	"github.com/esportsbar/backend/internal/constants"
	"github.com/esportsbar/backend/internal/util"
)

// RBAC 权限中间件：仅允许指定角色访问。
func RBAC(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, ok := c.Get("role")
		if !ok {
			util.AbortJSON(c, 401, constants.CodeUnauthorized, "未登录或登录已过期")
			return
		}
		roleStr, _ := role.(string)
		for _, r := range roles {
			if roleStr == r {
				c.Next()
				return
			}
		}
		util.AbortJSON(c, 403, constants.CodeForbidden, "没有操作权限，当前角色不允许该操作")
	}
}
