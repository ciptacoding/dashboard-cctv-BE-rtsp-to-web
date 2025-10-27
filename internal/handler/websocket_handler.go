package handler

import (
	"cctv-monitoring-backend/internal/service"
	ws "cctv-monitoring-backend/internal/websocket"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

// WebSocketHandler handles WebSocket connections
type WebSocketHandler struct {
	hub         *ws.Hub
	authService service.AuthService
}

// NewWebSocketHandler creates a new WebSocket handler
func NewWebSocketHandler(hub *ws.Hub, authService service.AuthService) *WebSocketHandler {
	return &WebSocketHandler{
		hub:         hub,
		authService: authService,
	}
}

// Upgrade upgrades HTTP connection to WebSocket
func (h *WebSocketHandler) Upgrade(c *fiber.Ctx) error {
	// Check if this is a WebSocket upgrade request
	if !websocket.IsWebSocketUpgrade(c) {
		return fiber.ErrUpgradeRequired
	}

	return c.Next()
}

// HandleConnection handles WebSocket connections
func (h *WebSocketHandler) HandleConnection(c *websocket.Conn) {
	// Optional: Authenticate WebSocket connection
	// You can pass token in query params: ws://localhost:8080/ws?token=xxx
	token := c.Query("token")
	if token != "" {
		jwtSecret := c.Locals("jwt_secret").(string)
		_, err := h.authService.VerifyToken(token, jwtSecret)
		if err != nil {
			log.Printf("WebSocket authentication failed: %v", err)
			c.WriteJSON(fiber.Map{
				"error": "Authentication failed",
			})
			c.Close()
			return
		}
	}

	// Handle the connection via Hub
	h.hub.HandleConnection(c)
}
