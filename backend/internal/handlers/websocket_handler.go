package handlers

import (
	"crypto-orderbook/internal/websocket"

	"github.com/gofiber/fiber/v2"
	ws "github.com/gofiber/websocket/v2"
)

type WebSocketHandler struct {
	hub *websocket.Hub
}

func NewWebSocketHandler(hub *websocket.Hub) *WebSocketHandler {
	return &WebSocketHandler{hub: hub}
}

func (h *WebSocketHandler) HandleWebSocket(c *ws.Conn) {
	client := websocket.NewClient(h.hub, c)
	h.hub.Register(client)

	go client.WritePump()
	client.ReadPump()
}

func (h *WebSocketHandler) UpgradeMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		if ws.IsWebSocketUpgrade(c) {
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	}
}
