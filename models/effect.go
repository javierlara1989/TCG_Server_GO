package models

import (
	"time"
)

// Effect represents an effect that can be applied to cards
type Effect struct {
	ID          int        `json:"id" db:"id"`
	Description string     `json:"description" db:"description"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
}

// CardEffect represents the relationship between cards and effects
type CardEffect struct {
	CardID   int `json:"card_id" db:"card_id"`
	EffectID int `json:"effect_id" db:"effect_id"`
}
