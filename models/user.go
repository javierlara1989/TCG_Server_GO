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
	ID        int        `json:"id" db:"id"`
	Nombre    string     `json:"nombre" db:"nombre" validate:"required,min=6,alpha"`
	Email     string     `json:"email" db:"email" validate:"required,email"`
	Password  string     `json:"-" db:"password" validate:"required,min=6,alphanum"` // Excluded from JSON serialization
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
}

// CreateUserRequest represents the data needed to create a new user
type CreateUserRequest struct {
	Nombre   string `json:"nombre" validate:"required,min=6,alpha"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6,alphanum"`
}
