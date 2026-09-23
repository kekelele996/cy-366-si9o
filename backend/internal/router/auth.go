package router

import (
	"github.com/gin-gonic/gin"

	"github.com/esportsbar/backend/internal/handler"
	"github.com/esportsbar/backend/internal/middleware"
)

// RegisterAuth 注册认证路由。
func RegisterAuth(rg *gin.RouterGroup, h *handler.AuthHandler, jwtSecret string) {
	auth := rg.Group("/auth")
	{
		auth.POST("/register", h.Register)
		auth.POST("/login", h.Login)
		auth.GET("/profile", middleware.Auth(jwtSecret), h.Profile)
	}
}
