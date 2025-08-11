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
	Time    int   `json:"time"`
	User    User  `json:"user,omitempty"`
	Rival   *User `json:"rival,omitempty"`
	Table   Table `json:"table,omitempty"`
}

type TableState struct {
	ID                    uint      `json:"id"`
	TableID               uint      `json:"table_id"`
	Log                   string    `json:"log"`
	OwnersDeckID          *uint     `json:"owners_deck_id,omitempty"`
	RivalsDeckID          *uint     `json:"rivals_deck_id,omitempty"`
	OwnersActiveMonster   []uint    `json:"owners_active_monster"`
	OwnersBenchMonster1   []uint    `json:"owners_bench_monster_1"`
	OwnersBenchMonster2   []uint    `json:"owners_bench_monster_2"`
	OwnersBenchMonster3   []uint    `json:"owners_bench_monster_3"`
	OwnersActiveMonsterHP *int      `json:"owners_active_monster_hp,omitempty"`
	OwnersBenchMonster1HP *int      `json:"owners_bench_monster_1_hp,omitempty"`
	OwnersBenchMonster2HP *int      `json:"owners_bench_monster_2_hp,omitempty"`
	OwnersBenchMonster3HP *int      `json:"owners_bench_monster_3_hp,omitempty"`
	OwnersGraveyard       []uint    `json:"owners_graveyard"`
	RivalsActiveMonster   []uint    `json:"rivals_active_monster"`
	RivalsBenchMonster1   []uint    `json:"rivals_bench_monster_1"`
	RivalsBenchMonster2   []uint    `json:"rivals_bench_monster_2"`
	RivalsBenchMonster3   []uint    `json:"rivals_bench_monster_3"`
	RivalsActiveMonsterHP *int      `json:"rivals_active_monster_hp,omitempty"`
	RivalsBenchMonster1HP *int      `json:"rivals_bench_monster_1_hp,omitempty"`
	RivalsBenchMonster2HP *int      `json:"rivals_bench_monster_2_hp,omitempty"`
	RivalsBenchMonster3HP *int      `json:"rivals_bench_monster_3_hp,omitempty"`
	RivalsGraveyard       []uint    `json:"rivals_graveyard"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
}
