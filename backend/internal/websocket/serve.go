package websocket

import (
	"log"
	"net/http"
	"social-network/backend/internal/auth"
)

// ServeWebSocket performs the upgrade handshake using the pre-authenticated context userID
func ServeWebSocket(hub *Hub, w http.ResponseWriter, r *http.Request) {
	// Extract the userID injected into the context by auth.Authenticate middleware
	userID, ok := auth.GetUserIDFromContext(r.Context())
	if !ok {
		log.Println("[WEBSOCKET ERROR] ServeWebSocket invoked without authenticated context")
		http.Error(w, `{"error":"Unauthorized edge entry"}`, http.StatusUnauthorized)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("[WEBSOCKET ERROR] Upgrade handshake failed: %v", err)
		return
	}

	client := &Client{
		Hub:    hub,
		Conn:   conn,
		UserID: userID,
		Send:   make(chan []byte, 256),
	}

	client.Hub.Register <- client

	// Start concurrent asynchronous processing goroutines
	go client.WritePump()
	go client.ReadPump()
}