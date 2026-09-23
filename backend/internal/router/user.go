package router

import (
	"github.com/gin-gonic/gin"

	"github.com/esportsbar/backend/internal/constants"
	"github.com/esportsbar/backend/internal/handler"
	"github.com/esportsbar/backend/internal/middleware"
)

// RegisterUser 注册用户管理路由。
func RegisterUser(rg *gin.RouterGroup, h *handler.UserHandler, jwtSecret string) {
	users := rg.Group("/users", middleware.Auth(jwtSecret))
	{
		users.GET("", middleware.RBAC(constants.RoleAdmin, constants.RoleStaff), h.List)
		users.GET("/:id", middleware.RBAC(constants.RoleAdmin, constants.RoleStaff), h.Get)
		users.POST("", middleware.RBAC(constants.RoleAdmin), h.Create)
		users.PUT("/:id", middleware.RBAC(constants.RoleAdmin), h.Update)
		users.DELETE("/:id", middleware.RBAC(constants.RoleAdmin), h.Delete)
	}
}
