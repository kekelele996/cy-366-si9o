package router

import (
	"github.com/gin-gonic/gin"

	"github.com/esportsbar/backend/internal/constants"
	"github.com/esportsbar/backend/internal/handler"
	"github.com/esportsbar/backend/internal/middleware"
)

// RegisterTimePackage 注册时长包路由。
func RegisterTimePackage(rg *gin.RouterGroup, h *handler.TimePackageHandler, jwtSecret string) {
	pkgs := rg.Group("/packages", middleware.Auth(jwtSecret))
	{
		pkgs.GET("", h.List)
		pkgs.GET("/active", h.ListActive)
		pkgs.GET("/:id", h.Get)
		pkgs.POST("", middleware.RBAC(constants.RoleAdmin, constants.RoleStaff), h.Create)
		pkgs.PUT("/:id", middleware.RBAC(constants.RoleAdmin, constants.RoleStaff), h.Update)
		pkgs.DELETE("/:id", middleware.RBAC(constants.RoleAdmin), h.Delete)
	}
}
