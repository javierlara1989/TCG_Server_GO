package models

import (
	"time"
)

type Table struct {
	ID         uint       `json:"id"`
	Category   string     `json:"category"`
	Privacy    string     `json:"privacy"`
	Password   *string    `json:"password,omitempty"`
	Prize      string     `json:"prize"`
	Amount     *int       `json:"amount,omitempty"`
	Winner     *bool      `json:"winner,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	FinishedAt *time.Time `json:"finished_at,omitempty"`
}

type UserTable struct {
	ID      uint  `json:"id"`
	UserID  uint  `json:"user_id"`
	RivalID *uint `json:"rival_id,omitempty"`
	TableID uint  `json:"table_id"`
	User    User  `json:"user,omitempty"`
	Rival   *User `json:"rival,omitempty"`
	Table   Table `json:"table,omitempty"`
}
