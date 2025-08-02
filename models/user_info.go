package models

import (
	"time"
)

// UserInfo represents game information for a user account
type UserInfo struct {
	ID         int       `json:"id" db:"id"`
	UserID     int       `json:"user_id" db:"user_id"`
	Level      int       `json:"level" db:"level"`
	Experience int       `json:"experience" db:"experience"`
	Money      int       `json:"money" db:"money"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}

// CreateUserInfoRequest represents the data needed to create user info
type CreateUserInfoRequest struct {
	UserID     int `json:"user_id" validate:"required"`
	Level      int `json:"level" validate:"min=1"`
	Experience int `json:"experience" validate:"min=0"`
	Money      int `json:"money" validate:"min=0"`
}

// UpdateUserInfoRequest represents the data needed to update user info
type UpdateUserInfoRequest struct {
	Level      *int `json:"level,omitempty" validate:"omitempty,min=1"`
	Experience *int `json:"experience,omitempty" validate:"omitempty,min=0"`
	Money      *int `json:"money,omitempty" validate:"omitempty,min=0"`
}

// UserInfoResponse represents the response for user info operations
type UserInfoResponse struct {
	UserInfo *UserInfo `json:"user_info"`
	Message  string    `json:"message"`
}
