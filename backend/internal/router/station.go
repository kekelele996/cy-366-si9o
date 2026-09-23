package router

import (
	"github.com/gin-gonic/gin"

	"github.com/esportsbar/backend/internal/constants"
	"github.com/esportsbar/backend/internal/handler"
	"github.com/esportsbar/backend/internal/middleware"
)

// RegisterStation 注册机位路由。
func RegisterStation(rg *gin.RouterGroup, h *handler.StationHandler, jwtSecret string) {
	stations := rg.Group("/stations", middleware.Auth(jwtSecret))
	{
		stations.GET("", h.List)
		stations.GET("/all", h.ListAll)
		stations.GET("/:id", h.Get)
		stations.POST("", middleware.RBAC(constants.RoleAdmin, constants.RoleStaff), h.Create)
		stations.PUT("/:id", middleware.RBAC(constants.RoleAdmin, constants.RoleStaff), h.Update)
		stations.PUT("/:id/status", middleware.RBAC(constants.RoleAdmin, constants.RoleStaff), h.UpdateStatus)
		stations.DELETE("/:id", middleware.RBAC(constants.RoleAdmin), h.Delete)
	}
}
