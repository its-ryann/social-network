package websocket

import (
	"log"
	"net/http"
	"social-network/backend/internal/auth"
)

func ServeWebSocket(hub *Hub, w http.ResponseWriter, r *http.Request) {
	// 1. Authenticate the incoming request
	userID, ok := auth.GetUserIDFromContext(r.Context())
	if !ok {
		log.Println("[WEBSOCKET] Failed to authenticate WebSocket request")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	conn, err = upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("[WEBSOCKET] Upgrade failed: %v", err)
		return
	}

	client := &Client{
		Hub:    hub,
		Conn:   conn,
		UserID: userID,
		Send:   make(chan []byte, 256),
	}

	hub.Register <- client

	go client.WritePump()
	go client.ReadPump()
}