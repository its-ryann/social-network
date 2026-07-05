package router

import (
	"database/sql"
	"net/http"
	"social-network/backend/internal/auth"
	"social-network/backend/internal/handlers"
	"social-network/backend/internal/websocket"
)

func NewRouter(db *sql.DB, hub *websocket.Hub) *http.ServeMux {
	mux := http.NewServeMux()
	
	userHandler := &handlers.UserHandler{DB: db}
	authMiddleware := auth.Authenticate(db)

	// Infrastructure verification
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"healthy"}`))
	})

	// Public Registration
	mux.HandleFunc("POST /register", userHandler.Register)

	// Protected Endpoint Verification Routine
	meHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, _ := auth.GetUserIDFromContext(r.Context())
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"authenticated_user_id":"` + userID + `"}`))
	})
	
	// Chain verification layer boundaries
	mux.Handle("GET /me", authMiddleware(meHandler))

	// WebSocket Endpoint wrapped with authentication middleware
	wsHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		websocket.ServeWebSocket(hub, db, w, r)
	})
	mux.Handle("GET /ws", authMiddleware(wsHandler))
	
	return mux
}