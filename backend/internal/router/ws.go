package router

import (
	"github.com/gin-gonic/gin"

	"github.com/esportsbar/backend/internal/handler"
)

// RegisterWS 注册 WebSocket 路由。
func RegisterWS(rg *gin.RouterGroup, h *handler.WSHandler) {
	rg.GET("/ws/stations", h.Stations)
}
