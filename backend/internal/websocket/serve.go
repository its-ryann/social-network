package websocket

import (
	"database/sql"
	"log"
	"net/http"
	"social-network/backend/internal/auth"
)

// ServeWebSocket performs the upgrade handshake using the pre-authenticated context userID
func ServeWebSocket(hub *Hub, db *sql.DB, w http.ResponseWriter, r *http.Request) {
	// Extract the userID injected into the context by auth.Authenticate middleware
	userID, ok := auth.GetUserIDFromContext(r.Context())
	if !ok {
		log.Println("[WEBSOCKET ERROR] ServeWebSocket invoked without authenticated context")
		http.Error(w, `{"error":"Unauthorized edge entry"}`, http.StatusUnauthorized)
		return
	}

	// Validate session token from cookie to ensure the user is still authenticated
	cookie, err := r.Cookie("session_token")
	if err != nil {
		http.Error(w, `{"error":"Session token missing"}`, http.StatusUnauthorized)
		return
	}

	var dbUserID string
	// Simple time query to prevent SQLite parsing variations and ensure consistent session validation
	query := `SELECT user_id FROM sessions WHERE id = ?`
	err = db.QueryRow(query, cookie.Value).Scan(&dbUserID)
	if err != nil {
		log.Printf("[WEBSOCKET CRITICAL] DB verification failed: %v", err)
		http.Error(w, `{"error":"Database lookup failure"}`, http.StatusInternalServerError)
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