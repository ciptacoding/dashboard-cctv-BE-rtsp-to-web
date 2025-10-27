package websocket

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/gofiber/websocket/v2"
)

// WSMessage adalah struktur pesan WebSocket
type WSMessage struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

// CameraStatusUpdate adalah struktur untuk update status kamera
type CameraStatusUpdate struct {
	ID       string `json:"id"`
	Status   string `json:"status"`
	LastSeen string `json:"last_seen,omitempty"`
}

// CameraStreamUpdate untuk notifikasi stream issues
type CameraStreamUpdate struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Status  string `json:"status"` // "frozen", "offline", "online"
	Message string `json:"message"`
}

// Hub maintains the set of active clients and broadcasts messages
type Hub struct {
	// Registered clients
	clients map[*websocket.Conn]bool

	// Broadcast channel
	broadcast chan WSMessage

	// Register channel
	register chan *websocket.Conn

	// Unregister channel
	unregister chan *websocket.Conn

	// Mutex untuk thread-safe operations
	mu sync.RWMutex
}

// NewHub creates a new Hub
func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*websocket.Conn]bool),
		broadcast:  make(chan WSMessage, 256),
		register:   make(chan *websocket.Conn),
		unregister: make(chan *websocket.Conn),
	}
}

// Run starts the hub
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()
			log.Printf("✓ WebSocket client connected. Total clients: %d", len(h.clients))

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				client.Close()
			}
			h.mu.Unlock()
			log.Printf("✓ WebSocket client disconnected. Total clients: %d", len(h.clients))

		case message := <-h.broadcast:
			h.mu.RLock()
			for client := range h.clients {
				err := client.WriteJSON(message)
				if err != nil {
					log.Printf("Error broadcasting to client: %v", err)
					// Unregister client jika error
					go func(c *websocket.Conn) {
						h.unregister <- c
					}(client)
				}
			}
			h.mu.RUnlock()
		}
	}
}

// Register adds a client to the hub
func (h *Hub) Register(conn *websocket.Conn) {
	h.register <- conn
}

// Unregister removes a client from the hub
func (h *Hub) Unregister(conn *websocket.Conn) {
	h.unregister <- conn
}

// Broadcast sends a message to all connected clients
func (h *Hub) Broadcast(message WSMessage) {
	h.broadcast <- message
}

// BroadcastCameraStatus broadcasts camera status update
func (h *Hub) BroadcastCameraStatus(cameraID, status, lastSeen string) {
	message := WSMessage{
		Type: "camera_status",
		Data: CameraStatusUpdate{
			ID:       cameraID,
			Status:   status,
			LastSeen: lastSeen,
		},
	}
	h.Broadcast(message)
}

// BroadcastStreamUpdate broadcasts stream update (frozen, offline, etc)
func (h *Hub) BroadcastStreamUpdate(cameraID, name, status, msg string) {
	message := WSMessage{
		Type: "stream_update",
		Data: CameraStreamUpdate{
			ID:      cameraID,
			Name:    name,
			Status:  status,
			Message: msg,
		},
	}
	h.Broadcast(message)
}

// GetClientCount returns the number of connected clients
func (h *Hub) GetClientCount() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.clients)
}

// HandleConnection handles a WebSocket connection
func (h *Hub) HandleConnection(c *websocket.Conn) {
	// Register client
	h.Register(c)
	defer h.Unregister(c)

	// Send welcome message
	welcomeMsg := WSMessage{
		Type: "connected",
		Data: map[string]interface{}{
			"message": "Connected to CCTV Monitoring WebSocket",
			"clients": h.GetClientCount(),
		},
	}
	if err := c.WriteJSON(welcomeMsg); err != nil {
		log.Printf("Error sending welcome message: %v", err)
		return
	}

	// Listen for messages (mostly for ping/pong)
	for {
		var msg WSMessage
		if err := c.ReadJSON(&msg); err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

		// Handle ping
		if msg.Type == "ping" {
			pongMsg := WSMessage{
				Type: "pong",
				Data: map[string]interface{}{
					"timestamp": msg.Data,
				},
			}
			if err := c.WriteJSON(pongMsg); err != nil {
				log.Printf("Error sending pong: %v", err)
				break
			}
		}
	}
}

// BroadcastJSON broadcasts a custom JSON message
func (h *Hub) BroadcastJSON(messageType string, data interface{}) {
	message := WSMessage{
		Type: messageType,
		Data: data,
	}
	h.Broadcast(message)
}

// GetClientsInfo returns information about connected clients (for debugging)
func (h *Hub) GetClientsInfo() map[string]interface{} {
	h.mu.RLock()
	defer h.mu.RUnlock()

	return map[string]interface{}{
		"total_clients": len(h.clients),
		"timestamp":     json.Number("0"), // You can add actual timestamp if needed
	}
}
