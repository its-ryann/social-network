package auth

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"
)

// GenerateSecureToken creates an unguessable 32-byte cryptographic session tracking string
func GenerateSecureToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// CreateSession generates a new tracking token and stores it inside SQLite
func CreateSession(db *sql.DB, userID string) (string, error) {
	token, err := GenerateSecureToken()
	if err != nil {
		return "", err
	}

	// Session expires explicitly in 24 hours
	expiresAt := time.Now().Add(24 * time.Hour)

	query := `INSERT INTO sessions (id, user_id, expires_at) VALUES (?, ?, ?)`
	_, err = db.Exec(query, token, userID, expiresAt)
	if err != nil {
		return "", fmt.Errorf("session persistence error: %w", err)
	}

	return token, nil
}

// SetSessionCookie applies rigid HTTP attributes to prevent common web attacks
func SetSessionCookie(w http.ResponseWriter, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    token,
		Path:     "/",
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,                  // Blocks XSS execution engines from extracting the token
		Secure:   false,                 // Toggle to true in production environments enforcing HTTPS
		SameSite: http.SameSiteLaxMode,  // Mitigates standard Cross-Site Request Forgery vectors
	})
}