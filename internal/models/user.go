package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID                  uuid.UUID  `json:"id"`
	Email               string     `json:"email"`
	PasswordHash        string     `json:"-"`
	Username            string     `json:"username"`
	IsActive            bool       `json:"is_active"`
	EmailVerified       bool       `json:"email_verified"`
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at"`
	LastLoginAt         *time.Time `json:"last_login_at,omitempty"`
	FailedLoginAttempts int        `json:"-"`
	LockedUntil         *time.Time `json:"-"`
}

type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email,max=100"`
	Username string `json:"username" validate:"required,min=3,max=50"`
	Password string `json:"password" validate:"required,min=6,max=128"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email,max=100"`
	Password string `json:"password" validate:"required,max=128"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type AuthResponse struct {
	UserID       uuid.UUID `json:"user_id"`
	AccessToken  string    `json:"token"`
	RefreshToken string    `json:"refresh_token,omitempty"`
}

type ErrorResponse struct {
	Error string `json:"error" example:"invalid email or password"`
}

func (u *User) ToResponse() *User {
	return &User{
		ID:            u.ID,
		Email:         u.Email,
		Username:      u.Username,
		IsActive:      u.IsActive,
		EmailVerified: u.EmailVerified,
		CreatedAt:     u.CreatedAt,
		UpdatedAt:     u.UpdatedAt,
		LastLoginAt:   u.LastLoginAt,
	}
}
