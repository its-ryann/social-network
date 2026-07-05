package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"social-network/backend/internal/auth"
	"time"
)

type UserHandler struct {
	DB *sql.DB
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	// Parse multipart payload bounding max buffer overhead to 10MB to mitigate DoS risks
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, `{"error":"Payload sizing threshold exceeded"}`, http.StatusBadRequest)
		return
	}

	email := r.FormValue("email")
	password := r.FormValue("password")
	firstName := r.FormValue("first_name")
	lastName := r.FormValue("last_name")
	dob := r.FormValue("dob")

	// Validate required fields
	if email == "" || password == "" || firstName == "" || lastName == "" || dob == "" {
		http.Error(w, `{"error":"Missing mandatory registration fields"}`, http.StatusUnprocessableEntity)
		return
	}

	// Encrypt raw user input password credentials
	hashedPassword, err := auth.HashPassword(password)
	if err != nil {
		http.Error(w, `{"error":"Internal cryptographic processing failure"}`, http.StatusInternalServerError)
		return
	}

	// Generate a unique user id string prefix
	userID := fmt.Sprintf("usr_%d", time.Now().UnixNano())

	// Persist data record safely using standard parameterized queries
	query := `INSERT INTO users (id, email, password_hash, first_name, last_name, date_of_birth) VALUES (?, ?, ?, ?, ?, ?)`
	_, err = h.DB.Exec(query, userID, email, hashedPassword, firstName, lastName, dob)
	if err != nil {
		http.Error(w, `{"error":"Account constraints collision"}`, http.StatusConflict)
		return
	}

	// Automatically establish user login session tracking
	token, err := auth.CreateSession(h.DB, userID)
	if err != nil {
		http.Error(w, `{"error":"State tracking engine generation error"}`, http.StatusInternalServerError)
		return
	}

	auth.SetSessionCookie(w, token)

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"message":"User account established successfully"}`))
}