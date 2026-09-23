package router

import (
	"github.com/gin-gonic/gin"

	"github.com/esportsbar/backend/internal/constants"
	"github.com/esportsbar/backend/internal/handler"
	"github.com/esportsbar/backend/internal/middleware"
)

// RegisterRecharge 注册充值路由。
func RegisterRecharge(rg *gin.RouterGroup, h *handler.RechargeHandler, jwtSecret string) {
	recharges := rg.Group("/recharges", middleware.Auth(jwtSecret))
	{
		recharges.POST("", middleware.RBAC(constants.RoleAdmin, constants.RoleStaff), h.Recharge)
		recharges.GET("/mine", h.ListRecharges)
		recharges.POST("/packages", h.BuyPackage)
		recharges.GET("/packages/orders", h.ListOrders)
	}
}
