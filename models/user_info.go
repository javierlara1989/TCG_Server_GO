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

// Deck represents a user's card deck
type Deck struct {
	ID     int    `json:"id" db:"id"`
	UserID int    `json:"user_id" db:"user_id"`
	Name   string `json:"name" db:"name"`
	Valid  bool   `json:"valid" db:"valid"`
}

// CreateDeckRequest represents the data needed to create a deck
type CreateDeckRequest struct {
	Name      string `json:"name" validate:"required,min=1,max=100"`
	CardIDs   []int  `json:"card_ids" validate:"required,min=1"`
	CardCount []int  `json:"card_count" validate:"required,min=1"`
}

// DeckResponse represents the response for deck operations
type DeckResponse struct {
	Deck    *Deck  `json:"deck"`
	Message string `json:"message"`
}

// DecksResponse represents the response for multiple decks
type DecksResponse struct {
	Decks   []Deck `json:"decks"`
	Message string `json:"message"`
}

// DeckCard represents a card in a deck
type DeckCard struct {
	DeckID int   `json:"deck_id" db:"deck_id"`
	CardID int   `json:"card_id" db:"card_id"`
	Number int   `json:"number" db:"number"`
	Card   *Card `json:"card,omitempty"`
}

// DeckWithCards represents a deck with its cards
type DeckWithCards struct {
	Deck  *Deck      `json:"deck"`
	Cards []DeckCard `json:"cards"`
}

// DeckWithCardsResponse represents the response for deck with cards
type DeckWithCardsResponse struct {
	DeckWithCards *DeckWithCards `json:"deck_with_cards"`
	Message       string         `json:"message"`
}
