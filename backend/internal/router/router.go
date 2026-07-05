package router

import (
	"net/http"
)

// NewRouter initializes an isolated mux using native method-matching (Go 1.22+)
func NewRouter() *http.ServeMux {
	mux := http.NewServeMux()
	
	// Register the core infrastructure health endpoint
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"healthy"}`))
	})
	
	return mux
}