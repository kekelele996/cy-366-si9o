package router

import (
	"github.com/gin-gonic/gin"

	"github.com/esportsbar/backend/internal/handler"
	"github.com/esportsbar/backend/internal/middleware"
)

// RegisterSession 注册上机记录路由。
func RegisterSession(rg *gin.RouterGroup, h *handler.SessionHandler, jwtSecret string) {
	sessions := rg.Group("/sessions", middleware.Auth(jwtSecret))
	{
		sessions.GET("", h.List)
		sessions.GET("/rank", h.Rank)
		sessions.POST("", h.Start)
		sessions.POST("/:id/renew", h.Renew)
		sessions.POST("/:id/end", h.End)
	}
}
