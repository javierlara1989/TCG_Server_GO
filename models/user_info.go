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

// UserCard represents a user's card inventory
type UserCard struct {
	ID        int       `json:"id" db:"id"`
	UserID    int       `json:"user_id" db:"user_id"`
	CardID    int       `json:"card_id" db:"card_id"`
	Amount    int       `json:"amount" db:"amount"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	Card      *Card     `json:"card,omitempty"`
}

// UserCardResponse represents the response for user card operations
type UserCardResponse struct {
	UserCard *UserCard `json:"user_card"`
	Message  string    `json:"message"`
}

// UserCardsResponse represents the response for multiple user cards
type UserCardsResponse struct {
	UserCards []UserCard `json:"user_cards"`
	Message   string     `json:"message"`
}
