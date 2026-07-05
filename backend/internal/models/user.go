package models

import "time"

type User struct {
	ID           string    `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"` // Never expose the hash in JSON payloads
	FirstName    string    `json:"first_name"`
	LastName     string    `json:"last_name"`
	DateOfBirth  string    `json:"date_of_birth"`
	AvatarURL    string    `json:"avatar_url,omitempty"`
	Nickname     string    `json:"nickname,omitempty"`
	AboutMe      string    `json:"about_me,omitempty"`
	IsPublic     bool      `json:"is_public"`
	CreatedAt    time.Time `json:"created_at"`
}