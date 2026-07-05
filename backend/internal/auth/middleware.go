package auth

import (
	"context"
	"database/sql"
	"net"
	"net/http"
	"time"
)

type contextKey string

const userIDKey contextKey = "userID"

func Authenticate(db *sql.DB, next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("session_token")
			if err != nil {
				http.Error(w, `{"error":"Authentication required"}`, http.StatusUnauthorized)
				return
			}

			var userID string
			var expiresAt time.Time

			query := `SELECT user_id, expires_at FROM sessions WHERE id = ?`
			err = db.QueryRow(query, cookie.Value).Scan(&userID, &expiresAt)
			if err != nil {
				if err == sql.ErrNoRows {
					http.Error(w, `{"error":"Invalid session token"}`, http.StatusUnauthorized)
				} else {
					http.Error(w, `{"error":"Database query failure"}`, http.StatusInternalServerError)
				}
				return
			}

			if time.Now().After(expiresAt) {
				http.Error(w, `{"error":"Session token expired"}`, http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), userIDKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}(next)
}

func GetUserIDFromContext(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(userIDKey).(string)
	return userID, ok
}