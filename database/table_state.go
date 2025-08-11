package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	"tcg-server-go/models"
)

// CreateTableState creates a new table state record (internal use only)
func CreateTableState(tableState *models.TableState) error {
	query := `
		INSERT INTO table_state (
			table_id, log, owners_deck_id, rivals_deck_id,
			owners_active_monster, owners_bench_monster_1, owners_bench_monster_2, owners_bench_monster_3,
			owners_active_monster_hp, owners_bench_monster_1_hp, owners_bench_monster_2_hp, owners_bench_monster_3_hp,
			owners_graveyard, rivals_active_monster, rivals_bench_monster_1, rivals_bench_monster_2, rivals_bench_monster_3,
			rivals_active_monster_hp, rivals_bench_monster_1_hp, rivals_bench_monster_2_hp, rivals_bench_monster_3_hp,
			rivals_graveyard
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	// Convert slices to JSON strings
	ownersActiveMonsterJSON, _ := json.Marshal(tableState.OwnersActiveMonster)
	ownersBenchMonster1JSON, _ := json.Marshal(tableState.OwnersBenchMonster1)
	ownersBenchMonster2JSON, _ := json.Marshal(tableState.OwnersBenchMonster2)
	ownersBenchMonster3JSON, _ := json.Marshal(tableState.OwnersBenchMonster3)
	ownersGraveyardJSON, _ := json.Marshal(tableState.OwnersGraveyard)
	rivalsActiveMonsterJSON, _ := json.Marshal(tableState.RivalsActiveMonster)
	rivalsBenchMonster1JSON, _ := json.Marshal(tableState.RivalsBenchMonster1)
	rivalsBenchMonster2JSON, _ := json.Marshal(tableState.RivalsBenchMonster2)
	rivalsBenchMonster3JSON, _ := json.Marshal(tableState.RivalsBenchMonster3)
	rivalsGraveyardJSON, _ := json.Marshal(tableState.RivalsGraveyard)

	result, err := DB.Exec(query,
		tableState.TableID, tableState.Log, tableState.OwnersDeckID, tableState.RivalsDeckID,
		ownersActiveMonsterJSON, ownersBenchMonster1JSON, ownersBenchMonster2JSON, ownersBenchMonster3JSON,
		tableState.OwnersActiveMonsterHP, tableState.OwnersBenchMonster1HP, tableState.OwnersBenchMonster2HP, tableState.OwnersBenchMonster3HP,
		ownersGraveyardJSON, rivalsActiveMonsterJSON, rivalsBenchMonster1JSON, rivalsBenchMonster2JSON, rivalsBenchMonster3JSON,
		tableState.RivalsActiveMonsterHP, tableState.RivalsBenchMonster1HP, tableState.RivalsBenchMonster2HP, tableState.RivalsBenchMonster3HP,
		rivalsGraveyardJSON,
	)
	if err != nil {
		return fmt.Errorf("error creating table state: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("error getting last insert id: %v", err)
	}

	tableState.ID = uint(id)
	return nil
}

// GetTableStateByTableID retrieves the current state of a table (internal use only)
func GetTableStateByTableID(tableID uint) (*models.TableState, error) {
	query := `
		SELECT id, table_id, log, owners_deck_id, rivals_deck_id,
			owners_active_monster, owners_bench_monster_1, owners_bench_monster_2, owners_bench_monster_3,
			owners_active_monster_hp, owners_bench_monster_1_hp, owners_bench_monster_2_hp, owners_bench_monster_3_hp,
			owners_graveyard, rivals_active_monster, rivals_bench_monster_1, rivals_bench_monster_2, rivals_bench_monster_3,
			rivals_active_monster_hp, rivals_bench_monster_1_hp, rivals_bench_monster_2_hp, rivals_bench_monster_3_hp,
			rivals_graveyard, created_at, updated_at
		FROM table_state
		WHERE table_id = ?
		ORDER BY created_at DESC
		LIMIT 1
	`

	var tableState models.TableState
	var ownersActiveMonsterJSON, ownersBenchMonster1JSON, ownersBenchMonster2JSON, ownersBenchMonster3JSON,
		ownersGraveyardJSON, rivalsActiveMonsterJSON, rivalsBenchMonster1JSON, rivalsBenchMonster2JSON,
		rivalsBenchMonster3JSON, rivalsGraveyardJSON sql.NullString

	err := DB.QueryRow(query, tableID).Scan(
		&tableState.ID, &tableState.TableID, &tableState.Log, &tableState.OwnersDeckID, &tableState.RivalsDeckID,
		&ownersActiveMonsterJSON, &ownersBenchMonster1JSON, &ownersBenchMonster2JSON, &ownersBenchMonster3JSON,
		&tableState.OwnersActiveMonsterHP, &tableState.OwnersBenchMonster1HP, &tableState.OwnersBenchMonster2HP, &tableState.OwnersBenchMonster3HP,
		&ownersGraveyardJSON, &rivalsActiveMonsterJSON, &rivalsBenchMonster1JSON, &rivalsBenchMonster2JSON, &rivalsBenchMonster3JSON,
		&tableState.RivalsActiveMonsterHP, &tableState.RivalsBenchMonster1HP, &tableState.RivalsBenchMonster2HP, &tableState.RivalsBenchMonster3HP,
		&rivalsGraveyardJSON, &tableState.CreatedAt, &tableState.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("error getting table state: %v", err)
	}

	// Parse JSON arrays
	if ownersActiveMonsterJSON.Valid {
		json.Unmarshal([]byte(ownersActiveMonsterJSON.String), &tableState.OwnersActiveMonster)
	}
	if ownersBenchMonster1JSON.Valid {
		json.Unmarshal([]byte(ownersBenchMonster1JSON.String), &tableState.OwnersBenchMonster1)
	}
	if ownersBenchMonster2JSON.Valid {
		json.Unmarshal([]byte(ownersBenchMonster2JSON.String), &tableState.OwnersBenchMonster2)
	}
	if ownersBenchMonster3JSON.Valid {
		json.Unmarshal([]byte(ownersBenchMonster3JSON.String), &tableState.OwnersBenchMonster3)
	}
	if ownersGraveyardJSON.Valid {
		json.Unmarshal([]byte(ownersGraveyardJSON.String), &tableState.OwnersGraveyard)
	}
	if rivalsActiveMonsterJSON.Valid {
		json.Unmarshal([]byte(rivalsActiveMonsterJSON.String), &tableState.RivalsActiveMonster)
	}
	if rivalsBenchMonster1JSON.Valid {
		json.Unmarshal([]byte(rivalsBenchMonster1JSON.String), &tableState.RivalsBenchMonster1)
	}
	if rivalsBenchMonster2JSON.Valid {
		json.Unmarshal([]byte(rivalsBenchMonster2JSON.String), &tableState.RivalsBenchMonster2)
	}
	if rivalsBenchMonster3JSON.Valid {
		json.Unmarshal([]byte(rivalsBenchMonster3JSON.String), &tableState.RivalsBenchMonster3)
	}
	if rivalsGraveyardJSON.Valid {
		json.Unmarshal([]byte(rivalsGraveyardJSON.String), &tableState.RivalsGraveyard)
	}

	return &tableState, nil
}

// UpdateTableState updates an existing table state (internal use only)
func UpdateTableState(tableState *models.TableState) error {
	query := `
		UPDATE table_state SET
			log = ?, owners_deck_id = ?, rivals_deck_id = ?,
			owners_active_monster = ?, owners_bench_monster_1 = ?, owners_bench_monster_2 = ?, owners_bench_monster_3 = ?,
			owners_active_monster_hp = ?, owners_bench_monster_1_hp = ?, owners_bench_monster_2_hp = ?, owners_bench_monster_3_hp = ?,
			owners_graveyard = ?, rivals_active_monster = ?, rivals_bench_monster_1 = ?, rivals_bench_monster_2 = ?, rivals_bench_monster_3 = ?,
			rivals_active_monster_hp = ?, rivals_bench_monster_1_hp = ?, rivals_bench_monster_2_hp = ?, rivals_bench_monster_3_hp = ?,
			rivals_graveyard = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`

	// Convert slices to JSON strings
	ownersActiveMonsterJSON, _ := json.Marshal(tableState.OwnersActiveMonster)
	ownersBenchMonster1JSON, _ := json.Marshal(tableState.OwnersBenchMonster1)
	ownersBenchMonster2JSON, _ := json.Marshal(tableState.OwnersBenchMonster2)
	ownersBenchMonster3JSON, _ := json.Marshal(tableState.OwnersBenchMonster3)
	ownersGraveyardJSON, _ := json.Marshal(tableState.OwnersGraveyard)
	rivalsActiveMonsterJSON, _ := json.Marshal(tableState.RivalsActiveMonster)
	rivalsBenchMonster1JSON, _ := json.Marshal(tableState.RivalsBenchMonster1)
	rivalsBenchMonster2JSON, _ := json.Marshal(tableState.RivalsBenchMonster2)
	rivalsBenchMonster3JSON, _ := json.Marshal(tableState.RivalsBenchMonster3)
	rivalsGraveyardJSON, _ := json.Marshal(tableState.RivalsGraveyard)

	_, err := DB.Exec(query,
		tableState.Log, tableState.OwnersDeckID, tableState.RivalsDeckID,
		ownersActiveMonsterJSON, ownersBenchMonster1JSON, ownersBenchMonster2JSON, ownersBenchMonster3JSON,
		tableState.OwnersActiveMonsterHP, tableState.OwnersBenchMonster1HP, tableState.OwnersBenchMonster2HP, tableState.OwnersBenchMonster3HP,
		ownersGraveyardJSON, rivalsActiveMonsterJSON, rivalsBenchMonster1JSON, rivalsBenchMonster2JSON, rivalsBenchMonster3JSON,
		tableState.RivalsActiveMonsterHP, tableState.RivalsBenchMonster1HP, tableState.RivalsBenchMonster2HP, tableState.RivalsBenchMonster3HP,
		rivalsGraveyardJSON, tableState.ID,
	)
	if err != nil {
		return fmt.Errorf("error updating table state: %v", err)
	}

	return nil
}

// DeleteTableState deletes a table state record (internal use only)
func DeleteTableState(id uint) error {
	query := "DELETE FROM table_state WHERE id = ?"
	_, err := DB.Exec(query, id)
	if err != nil {
		return fmt.Errorf("error deleting table state: %v", err)
	}
	return nil
}

// GetTableStateHistory retrieves the history of table states for a specific table (internal use only)
func GetTableStateHistory(tableID uint, limit int) ([]models.TableState, error) {
	query := `
		SELECT id, table_id, log, owners_deck_id, rivals_deck_id,
			owners_active_monster, owners_bench_monster_1, owners_bench_monster_2, owners_bench_monster_3,
			owners_active_monster_hp, owners_bench_monster_1_hp, owners_bench_monster_2_hp, owners_bench_monster_3_hp,
			owners_graveyard, rivals_active_monster, rivals_bench_monster_1, rivals_bench_monster_2, rivals_bench_monster_3,
			rivals_active_monster_hp, rivals_bench_monster_1_hp, rivals_bench_monster_2_hp, rivals_bench_monster_3_hp,
			rivals_graveyard, created_at, updated_at
		FROM table_state
		WHERE table_id = ?
		ORDER BY created_at DESC
		LIMIT ?
	`

	rows, err := DB.Query(query, tableID, limit)
	if err != nil {
		return nil, fmt.Errorf("error getting table state history: %v", err)
	}
	defer rows.Close()

	var tableStates []models.TableState
	for rows.Next() {
		var tableState models.TableState
		var ownersActiveMonsterJSON, ownersBenchMonster1JSON, ownersBenchMonster2JSON, ownersBenchMonster3JSON,
			ownersGraveyardJSON, rivalsActiveMonsterJSON, rivalsBenchMonster1JSON, rivalsBenchMonster2JSON,
			rivalsBenchMonster3JSON, rivalsGraveyardJSON sql.NullString

		err := rows.Scan(
			&tableState.ID, &tableState.TableID, &tableState.Log, &tableState.OwnersDeckID, &tableState.RivalsDeckID,
			&ownersActiveMonsterJSON, &ownersBenchMonster1JSON, &ownersBenchMonster2JSON, &ownersBenchMonster3JSON,
			&tableState.OwnersActiveMonsterHP, &tableState.OwnersBenchMonster1HP, &tableState.OwnersBenchMonster2HP, &tableState.OwnersBenchMonster3HP,
			&ownersGraveyardJSON, &rivalsActiveMonsterJSON, &rivalsBenchMonster1JSON, &rivalsBenchMonster2JSON, &rivalsBenchMonster3JSON,
			&tableState.RivalsActiveMonsterHP, &tableState.RivalsBenchMonster1HP, &tableState.RivalsBenchMonster2HP, &tableState.RivalsBenchMonster3HP,
			&rivalsGraveyardJSON, &tableState.CreatedAt, &tableState.UpdatedAt,
		)
		if err != nil {
			log.Printf("Error scanning table state row: %v", err)
			continue
		}

		// Parse JSON arrays
		if ownersActiveMonsterJSON.Valid {
			json.Unmarshal([]byte(ownersActiveMonsterJSON.String), &tableState.OwnersActiveMonster)
		}
		if ownersBenchMonster1JSON.Valid {
			json.Unmarshal([]byte(ownersBenchMonster1JSON.String), &tableState.OwnersBenchMonster1)
		}
		if ownersBenchMonster2JSON.Valid {
			json.Unmarshal([]byte(ownersBenchMonster2JSON.String), &tableState.OwnersBenchMonster2)
		}
		if ownersBenchMonster3JSON.Valid {
			json.Unmarshal([]byte(ownersBenchMonster3JSON.String), &tableState.OwnersBenchMonster3)
		}
		if ownersGraveyardJSON.Valid {
			json.Unmarshal([]byte(ownersGraveyardJSON.String), &tableState.OwnersGraveyard)
		}
		if rivalsActiveMonsterJSON.Valid {
			json.Unmarshal([]byte(rivalsActiveMonsterJSON.String), &tableState.RivalsActiveMonster)
		}
		if rivalsBenchMonster1JSON.Valid {
			json.Unmarshal([]byte(rivalsBenchMonster1JSON.String), &tableState.RivalsBenchMonster1)
		}
		if rivalsBenchMonster2JSON.Valid {
			json.Unmarshal([]byte(rivalsBenchMonster2JSON.String), &tableState.RivalsBenchMonster2)
		}
		if rivalsBenchMonster3JSON.Valid {
			json.Unmarshal([]byte(rivalsBenchMonster3JSON.String), &tableState.RivalsBenchMonster3)
		}
		if rivalsGraveyardJSON.Valid {
			json.Unmarshal([]byte(rivalsGraveyardJSON.String), &tableState.RivalsGraveyard)
		}

		tableStates = append(tableStates, tableState)
	}

	return tableStates, nil
}
