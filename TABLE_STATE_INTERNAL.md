# Table State - Internal Usage

This document describes the internal usage of the `table_state` table and its associated functions.

## Overview

The `table_state` table is designed to store the current state of game tables internally. It is not exposed through public API endpoints and is meant to be used only by the server's internal game logic.

## Table Structure

The `table_state` table contains the following fields:

- `id`: Primary key
- `table_id`: Foreign key to the tables table
- `log`: Long text description of the current game state
- `owners_deck_id`: ID of the owner's deck (nullable)
- `rivals_deck_id`: ID of the rival's deck (nullable)
- `owners_active_monster`: JSON array of card IDs for owner's active monster(s)
- `owners_bench_monster_1`: JSON array of card IDs for owner's first bench monster
- `owners_bench_monster_2`: JSON array of card IDs for owner's second bench monster
- `owners_bench_monster_3`: JSON array of card IDs for owner's third bench monster
- `owners_active_monster_hp`: HP value for owner's active monster (nullable)
- `owners_bench_monster_1_hp`: HP value for owner's first bench monster (nullable)
- `owners_bench_monster_2_hp`: HP value for owner's second bench monster (nullable)
- `owners_bench_monster_3_hp`: HP value for owner's third bench monster (nullable)
- `owners_graveyard`: JSON array of card IDs in owner's graveyard
- `rivals_active_monster`: JSON array of card IDs for rival's active monster(s)
- `rivals_bench_monster_1`: JSON array of card IDs for rival's first bench monster
- `rivals_bench_monster_2`: JSON array of card IDs for rival's second bench monster
- `rivals_bench_monster_3`: JSON array of card IDs for rival's third bench monster
- `rivals_active_monster_hp`: HP value for rival's active monster (nullable)
- `rivals_bench_monster_1_hp`: HP value for rival's first bench monster (nullable)
- `rivals_bench_monster_2_hp`: HP value for rival's second bench monster (nullable)
- `rivals_bench_monster_3_hp`: HP value for rival's third bench monster (nullable)
- `rivals_graveyard`: JSON array of card IDs in rival's graveyard
- `created_at`: Timestamp when the state was created
- `updated_at`: Timestamp when the state was last updated

## Internal Functions

### CreateTableState(tableState *models.TableState) error

Creates a new table state record.

**Usage:**
```go
tableState := &models.TableState{
    TableID: 1,
    Log: "Game started. Player 1 drew 5 cards.",
    OwnersDeckID: &deckID1,
    RivalsDeckID: &deckID2,
    OwnersActiveMonster: []uint{101, 102},
    OwnersBenchMonster1: []uint{103},
    OwnersBenchMonster2: []uint{},
    OwnersBenchMonster3: []uint{},
    OwnersActiveMonsterHP: &hp1,
    OwnersBenchMonster1HP: &hp2,
    OwnersBenchMonster2HP: nil,
    OwnersBenchMonster3HP: nil,
    OwnersGraveyard: []uint{},
    RivalsActiveMonster: []uint{201},
    RivalsBenchMonster1: []uint{202, 203},
    RivalsBenchMonster2: []uint{},
    RivalsBenchMonster3: []uint{},
    RivalsActiveMonsterHP: &hp3,
    RivalsBenchMonster1HP: &hp4,
    RivalsBenchMonster2HP: nil,
    RivalsBenchMonster3HP: nil,
    RivalsGraveyard: []uint{204},
}

err := database.CreateTableState(tableState)
if err != nil {
    // Handle error
}
```

### GetTableStateByTableID(tableID uint) (*models.TableState, error)

Retrieves the current state of a table.

**Usage:**
```go
tableState, err := database.GetTableStateByTableID(1)
if err != nil {
    // Handle error
}
if tableState == nil {
    // No state found for this table
}
```

### UpdateTableState(tableState *models.TableState) error

Updates an existing table state.

**Usage:**
```go
// First get the current state
tableState, err := database.GetTableStateByTableID(1)
if err != nil {
    // Handle error
}

// Modify the state
tableState.Log = "Player 1 attacked with monster 101. Damage dealt: 500."
tableState.OwnersActiveMonsterHP = &newHP

// Update the state
err = database.UpdateTableState(tableState)
if err != nil {
    // Handle error
}
```

### DeleteTableState(id uint) error

Deletes a table state record.

**Usage:**
```go
err := database.DeleteTableState(1)
if err != nil {
    // Handle error
}
```

### GetTableStateHistory(tableID uint, limit int) ([]models.TableState, error)

Retrieves the history of table states for a specific table.

**Usage:**
```go
tableStates, err := database.GetTableStateHistory(1, 10)
if err != nil {
    // Handle error
}

for _, state := range tableStates {
    fmt.Printf("State ID: %d, Log: %s\n", state.ID, state.Log)
}
```

## Model Structure

The `TableState` model is defined in `models/table.go`:

```go
type TableState struct {
    ID                      uint      `json:"id"`
    TableID                 uint      `json:"table_id"`
    Log                     string    `json:"log"`
    OwnersDeckID            *uint     `json:"owners_deck_id,omitempty"`
    RivalsDeckID            *uint     `json:"rivals_deck_id,omitempty"`
    OwnersActiveMonster     []uint    `json:"owners_active_monster"`
    OwnersBenchMonster1     []uint    `json:"owners_bench_monster_1"`
    OwnersBenchMonster2     []uint    `json:"owners_bench_monster_2"`
    OwnersBenchMonster3     []uint    `json:"owners_bench_monster_3"`
    OwnersActiveMonsterHP   *int      `json:"owners_active_monster_hp,omitempty"`
    OwnersBenchMonster1HP   *int      `json:"owners_bench_monster_1_hp,omitempty"`
    OwnersBenchMonster2HP   *int      `json:"owners_bench_monster_2_hp,omitempty"`
    OwnersBenchMonster3HP   *int      `json:"owners_bench_monster_3_hp,omitempty"`
    OwnersGraveyard         []uint    `json:"owners_graveyard"`
    RivalsActiveMonster     []uint    `json:"rivals_active_monster"`
    RivalsBenchMonster1     []uint    `json:"rivals_bench_monster_1"`
    RivalsBenchMonster2     []uint    `json:"rivals_bench_monster_2"`
    RivalsBenchMonster3     []uint    `json:"rivals_bench_monster_3"`
    RivalsActiveMonsterHP   *int      `json:"rivals_active_monster_hp,omitempty"`
    RivalsBenchMonster1HP   *int      `json:"rivals_bench_monster_1_hp,omitempty"`
    RivalsBenchMonster2HP   *int      `json:"rivals_bench_monster_2_hp,omitempty"`
    RivalsBenchMonster3HP   *int      `json:"rivals_bench_monster_3_hp,omitempty"`
    RivalsGraveyard         []uint    `json:"rivals_graveyard"`
    CreatedAt               time.Time `json:"created_at"`
    UpdatedAt               time.Time `json:"updated_at"`
}
```

## Notes

1. **Internal Use Only**: These functions are designed for internal server use and are not exposed through API endpoints.

2. **JSON Arrays**: Monster arrays and graveyard arrays are stored as JSON in the database and automatically converted to/from Go slices.

3. **Nullable Fields**: HP values and deck IDs can be null when not applicable.

4. **Foreign Keys**: The table has foreign key constraints to ensure data integrity with the `tables` and `decks` tables.

5. **Indexing**: The table is indexed on `table_id`, `owners_deck_id`, and `rivals_deck_id` for efficient queries.

6. **Cascade Deletion**: When a table is deleted, all associated table states are automatically deleted.

7. **History**: The `GetTableStateHistory` function returns states in descending order by creation time (most recent first).
