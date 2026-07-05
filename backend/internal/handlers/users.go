package handlers

import (
	"fmt"
	"net/http"
	"social-network/backend/internal/auth"
	"time"
)

type UserHandler struct {
	DB *Database
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(10 << 20); err != nil { // Limit to 10MB
		http.Error(w, "Failed to parse form data", http.StatusBadRequest)
		return
	}

	email := r.FormValue("email")
	password := r.FormValue("password")
	firstName := r.FormValue("first_name")
	lastName := r.FormValue("last_name")
	dateOfBirth := r.FormValue("date_of_birth")

	if email == "" || password == "" || firstName == "" || lastName == "" || dateOfBirth == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	passwordHash, err := auth.HashPassword(password)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	userID := fmt.Sprintf("%d", time.Now().UnixNano()) // Simple unique ID generation

	user := &User{
		ID:           userID,
		Email:        email,
		PasswordHash: passwordHash,
		FirstName:    firstName,
		LastName:     lastName,
		DateOfBirth:  dateOfBirth,
		IsPublic:     true, // Default to public profile
		CreatedAt:    time.Now(),
	}

	if err := h.DB.CreateUser(user); err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	token, err := auth.CreateSession(h.DB.DB, userID)
	if err != nil {
		http.Error(w, "Failed to create session", http.StatusInternalServerError)
		return
	}

	auth.SetSessionCookie(w, token)

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "User registered successfully")
}