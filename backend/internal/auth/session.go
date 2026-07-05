package auth

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"
)

func GenerateSecureToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func CreateSession(db *sql.DB, userID string) (string, error) {
	token, err := GenerateSecureToken()
	if err != nil {
		return "", fmt.Errorf("failed to generate secure token: %w", err)
	}

	expiration := time.Now().Add(24 * time.Hour) // Session valid for 24 hours

	query := `INSERT INTO sessions (user_id, token, expires_at) VALUES ($1, $2, $3)`
	if _, err := db.Exec(query, userID, token, expiration); err != nil {
		return "", fmt.Errorf("failed to create session: %w", err)
	}

	return token, nil
}

func SetSessionCookie(w http.ResponseWriter, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   86400, // 24 hours
	})
}