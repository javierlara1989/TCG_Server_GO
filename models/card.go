package models

import (
	"time"
)

// CardType represents the type of card
type CardType string

const (
	CardTypeMonster CardType = "Monster"
	CardTypeSpell   CardType = "Spell"
	CardTypeEnergy  CardType = "Energy"
)

// CardElement represents the element of the card
type CardElement string

const (
	CardElementFire    CardElement = "Fire"
	CardElementWater   CardElement = "Water"
	CardElementWind    CardElement = "Wind"
	CardElementEarth   CardElement = "Earth"
	CardElementNeutral CardElement = "Neutral"
	CardElementHoly    CardElement = "Holy"
	CardElementDark    CardElement = "Dark"
)

// Card represents a card in the game
type Card struct {
	ID        int         `json:"id" db:"id"`
	Name      string      `json:"name" db:"name"`
	Type      CardType    `json:"type" db:"type"`
	Legend    string      `json:"legend" db:"legend"`
	Element   CardElement `json:"element" db:"element"`
	CreatedAt time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt time.Time   `json:"updated_at" db:"updated_at"`
}

// CreateCardRequest represents the data needed to create a card
type CreateCardRequest struct {
	Name    string      `json:"name" validate:"required,min=1"`
	Type    CardType    `json:"type" validate:"required,oneof=Monster Spell Energy"`
	Legend  string      `json:"legend" validate:"required"`
	Element CardElement `json:"element" validate:"required,oneof=Fire Water Wind Earth Neutral Holy Dark"`
}

// UpdateCardRequest represents the data needed to update a card
type UpdateCardRequest struct {
	Name    *string      `json:"name,omitempty" validate:"omitempty,min=1"`
	Type    *CardType    `json:"type,omitempty" validate:"omitempty,oneof=Monster Spell Energy"`
	Legend  *string      `json:"legend,omitempty" validate:"omitempty"`
	Element *CardElement `json:"element,omitempty" validate:"omitempty,oneof=Fire Water Wind Earth Neutral Holy Dark"`
}

// CardResponse represents the response for card operations
type CardResponse struct {
	Card    *Card  `json:"card"`
	Message string `json:"message"`
}

// CardsResponse represents the response for multiple cards
type CardsResponse struct {
	Cards   []*Card `json:"cards"`
	Message string  `json:"message"`
}
