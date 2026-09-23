package router

import (
	"github.com/gin-gonic/gin"

	"github.com/esportsbar/backend/internal/constants"
	"github.com/esportsbar/backend/internal/handler"
	"github.com/esportsbar/backend/internal/middleware"
)

// RegisterReservation 注册预约路由。
func RegisterReservation(rg *gin.RouterGroup, h *handler.ReservationHandler, jwtSecret string) {
	reservations := rg.Group("/reservations", middleware.Auth(jwtSecret))
	{
		reservations.GET("", h.List)
		reservations.POST("", h.Create)
		reservations.POST("/:id/confirm", middleware.RBAC(constants.RoleAdmin, constants.RoleStaff), h.Confirm)
		reservations.POST("/:id/cancel", h.Cancel)
		reservations.POST("/:id/checkin", middleware.RBAC(constants.RoleAdmin, constants.RoleStaff), h.CheckIn)
	}
}
