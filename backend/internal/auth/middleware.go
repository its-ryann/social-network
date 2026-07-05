package auth

import (
	"context"
	"database/sql"
	"net/http"
	"time"
)

type contextKey string
const UserIDKey contextKey = "userID"

// Authenticate intercepts incoming requests to enforce token evaluation routines
func Authenticate(db *sql.DB) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("session_token")
			if err != nil {
				http.Error(w, `{"error":"Authentication state missing"}`, http.StatusUnauthorized)
				return
			}

			var userID string
			var expiresAt time.Time

			// Query SQLite to verify the session exists and is still valid
			query := `SELECT user_id, expires_at FROM sessions WHERE id = ?`
			err = db.QueryRow(query, cookie.Value).Scan(&userID, &expiresAt)
			if err != nil {
				if err == sql.ErrNoRows {
					http.Error(w, `{"error":"Invalid authentication criteria"}`, http.StatusUnauthorized)
					return
				}
				http.Error(w, `{"error":"Database integrity connection error"}`, http.StatusInternalServerError)
				return
			}

			// Validate token lifetime expiration window
			if time.Now().After(expiresAt) {
				http.Error(w, `{"error":"Authentication state expired"}`, http.StatusUnauthorized)
				return
			}

			// Inject securely verified data directly into the active HTTP request context
			ctx := context.WithValue(r.Context(), UserIDKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetUserIDFromContext abstracts context parsing structures defensively
func GetUserIDFromContext(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(UserIDKey).(string)
	return userID, ok
}