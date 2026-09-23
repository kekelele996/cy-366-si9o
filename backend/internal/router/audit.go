package router

import (
	"github.com/gin-gonic/gin"

	"github.com/esportsbar/backend/internal/constants"
	"github.com/esportsbar/backend/internal/handler"
	"github.com/esportsbar/backend/internal/middleware"
)

// RegisterAudit 注册审计路由。
func RegisterAudit(rg *gin.RouterGroup, h *handler.AuditHandler, jwtSecret string) {
	audits := rg.Group("/audits", middleware.Auth(jwtSecret), middleware.RBAC(constants.RoleAdmin, constants.RoleStaff))
	{
		audits.GET("", h.List)
	}
}
