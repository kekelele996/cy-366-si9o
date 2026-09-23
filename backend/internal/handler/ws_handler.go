package handler

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// StationHub 机位状态广播中心。
type StationHub struct {
	clients   map[*websocket.Conn]bool
	broadcast chan []byte
	logger    *slog.Logger
}

// NewStationHub 构造广播中心。
func NewStationHub(logger *slog.Logger) *StationHub {
	return &StationHub{
		clients:   make(map[*websocket.Conn]bool),
		broadcast: make(chan []byte, 64),
		logger:    logger,
	}
}

// Run 启动广播协程。
func (h *StationHub) Run() {
	for {
		msg := <-h.broadcast
		for conn := range h.clients {
			if err := conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				h.logger.Warn("ws write failed", "err", err)
				conn.Close()
				delete(h.clients, conn)
			}
		}
	}
}

// Publish 广播机位状态消息。
func (h *StationHub) Publish(msg []byte) {
	h.broadcast <- msg
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// WSHandler WebSocket 接口处理器。
type WSHandler struct {
	hub    *StationHub
	logger *slog.Logger
}

// NewWSHandler 构造 WebSocket 接口处理器。
func NewWSHandler(hub *StationHub, logger *slog.Logger) *WSHandler {
	return &WSHandler{hub: hub, logger: logger}
}

// Stations 机位状态实时推送。
func (h *WSHandler) Stations(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		h.logger.Warn("ws upgrade failed", "err", err)
		return
	}
	h.hub.clients[conn] = true
	defer func() {
		delete(h.hub.clients, conn)
		conn.Close()
	}()
	for {
		if _, _, err := conn.ReadMessage(); err != nil {
			break
		}
	}
}
