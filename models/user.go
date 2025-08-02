package models

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type ValidateResponse struct {
	Message string `json:"message"`
}

type Claims struct {
	UserID int    `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

type User struct {
	ID                      int        `json:"id" db:"id"`
	Nombre                  string     `json:"nombre" db:"nombre" validate:"required,min=6,alpha"`
	Email                   string     `json:"email" db:"email" validate:"required,email"`
	Password                string     `json:"-" db:"password" validate:"required,min=6,alphanum"` // Excluded from JSON serialization
	ValidationCode          *string    `json:"-" db:"validation_code"`
	ValidationCodeExpiresAt *time.Time `json:"-" db:"validation_code_expires_at"`
	ValidatedAt             *time.Time `json:"validated_at" db:"validated_at"`
	CreatedAt               time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt               time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt               *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
}

// CreateUserRequest represents the data needed to create a new user
type CreateUserRequest struct {
	Nombre   string `json:"nombre" validate:"required,min=6,alpha"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6,alphanum"`
}

// VerifyEmailRequest represents the data needed to verify an email
type VerifyEmailRequest struct {
	Email          string `json:"email" validate:"required,email"`
	ValidationCode string `json:"validation_code" validate:"required"`
}

// VerifyEmailResponse represents the response for email verification
type VerifyEmailResponse struct {
	Message string `json:"message"`
	UserID  int    `json:"user_id"`
}

// ResendCodeRequest represents the request to resend validation code
type ResendCodeRequest struct {
	Email string `json:"email" validate:"required,email"`
}

// ResendCodeResponse represents the response for resending validation code
type ResendCodeResponse struct {
	Message string `json:"message"`
}
