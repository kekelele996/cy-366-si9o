package router

import (
	"github.com/gin-gonic/gin"

	"github.com/esportsbar/backend/internal/handler"
	"github.com/esportsbar/backend/internal/middleware"
)

// RegisterDashboard 注册看板路由。
func RegisterDashboard(rg *gin.RouterGroup, h *handler.DashboardHandler, jwtSecret string) {
	dash := rg.Group("/dashboard", middleware.Auth(jwtSecret))
	{
		dash.GET("/summary", h.Summary)
	}
}
