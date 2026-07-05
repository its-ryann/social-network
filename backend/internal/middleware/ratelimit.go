package middleware

import (
	"net/http"
	"net"
	"time"
	"sync"
)

type visitor struct {
	lastSeen time.Time
	tokens   int
}

var (
	visitors = make(map[string]*visitor)
	mu       sync.Mutex
)


// RateLimit is a middleware that limits the number of requests from a single IP address
func RateLimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			http.Error(w, "Unable to determine IP address", http.StatusInternalServerError)
			return
		}

		mu.Lock()
		v, exists := visitors[ip]
		if !exists {
			v = &visitor{lastSeen: time.Now(), tokens: 10} // Allow 10 requests initially
			visitors[ip] = v
		}

		// Refill tokens based on time elapsed
		elapsed := time.Since(v.lastSeen)
		v.tokens += int(elapsed.Seconds()) // 1 token per second
		if v.tokens > 10 {
			v.tokens = 10 // Cap at 10 tokens
		}
		v.lastSeen = time.Now()

		if v.tokens <= 0 {
			mu.Unlock()
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}

		v.tokens--
		mu.Unlock()

		next.ServeHTTP(w, r)
	})
}