package middleware

import (
	"log"
	"net/http"
	"time"
)

// responseWriterInterceptor captures the status code sent to the client
type responseWriterInterceptor struct {
	http.ResponseWriter
	statusCode int
}

func (rwi *responseWriterInterceptor) WriteHeader(statusCode int) {
	rwi.statusCode = statusCode
	rwi.ResponseWriter.WriteHeader(statusCode)
}

// Logger intercepts requests to track duration, HTTP method, path, and response state
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		
		// Wrap the writer to capture the response code (default to 200 if not explicitly called)
		rwi := &responseWriterInterceptor{ResponseWriter: w, statusCode: http.StatusOK}
		
		next.ServeHTTP(rwi, r)
		
		log.Printf("[NETEDGE] %s %s | Status: %d | Latency: %s", 
			r.Method, r.URL.Path, rwi.statusCode, time.Since(start))
	})
}