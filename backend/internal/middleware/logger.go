package middleware

import (
	"bufio"
	"fmt"
	"log"
	"net"
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

// Hijack allows the response writer to support connection hijacking for WebSocket upgrades
func (rwi *responseWriterInterceptor) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	hijacker, ok := rwi.ResponseWriter.(http.Hijacker)
	if !ok {
		return nil, nil, fmt.Errorf("response writer does not support hijacking")
	}
	return hijacker.Hijack()
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