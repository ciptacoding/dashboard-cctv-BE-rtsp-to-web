package models

import "time"

// User merepresentasikan struktur data user dalam database
type User struct {
	ID           string    `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"` // Tidak di-serialize ke JSON
	Role         string    `json:"role"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// LoginRequest adalah struktur untuk request login
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginResponse adalah struktur untuk response login
type LoginResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

// CreateUserRequest adalah struktur untuk membuat user baru
type CreateUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}
