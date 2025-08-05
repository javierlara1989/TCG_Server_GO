package database

import (
	"database/sql"
	"fmt"
)

// GetEffectByID retrieves an effect by its ID
func GetEffectByID(id int) (*sql.Row, error) {
	query := `
		SELECT id, description, created_at, updated_at, deleted_at
		FROM effects WHERE id = ? AND deleted_at IS NULL
	`

	row := DB.QueryRow(query, id)
	return row, nil
}

// GetAllEffects retrieves all non-deleted effects
func GetAllEffects() (*sql.Rows, error) {
	query := `
		SELECT id, description, created_at, updated_at, deleted_at
		FROM effects WHERE deleted_at IS NULL
		ORDER BY id
	`

	rows, err := DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error querying effects: %v", err)
	}

	return rows, nil
}

// SoftDeleteEffect soft deletes an effect by setting deleted_at
func SoftDeleteEffect(id int) error {
	query := `
		UPDATE effects 
		SET deleted_at = NOW(), updated_at = NOW()
		WHERE id = ? AND deleted_at IS NULL
	`

	_, err := DB.Exec(query, id)
	if err != nil {
		return fmt.Errorf("error soft deleting effect: %v", err)
	}

	return nil
}

// HardDeleteEffect permanently deletes an effect and its card associations
func HardDeleteEffect(id int) error {
	// Start a transaction
	tx, err := DB.Begin()
	if err != nil {
		return fmt.Errorf("error starting transaction: %v", err)
	}

	// Delete associated card_effects first
	_, err = tx.Exec("DELETE FROM card_effects WHERE effect_id = ?", id)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("error deleting card effects: %v", err)
	}

	// Delete the effect
	_, err = tx.Exec("DELETE FROM effects WHERE id = ?", id)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("error deleting effect: %v", err)
	}

	return tx.Commit()
}

// CreateCardEffect creates a relationship between a card and an effect
func CreateCardEffect(cardID, effectID int) error {
	query := `
		INSERT INTO card_effects (card_id, effect_id)
		VALUES (?, ?)
	`

	_, err := DB.Exec(query, cardID, effectID)
	if err != nil {
		return fmt.Errorf("error creating card effect: %v", err)
	}

	return nil
}

// GetEffectsByCardID retrieves all effects for a specific card
func GetEffectsByCardID(cardID int) (*sql.Rows, error) {
	query := `
		SELECT e.id, e.description, e.created_at, e.updated_at, e.deleted_at
		FROM effects e
		JOIN card_effects ce ON e.id = ce.effect_id
		WHERE ce.card_id = ? AND e.deleted_at IS NULL
		ORDER BY e.id
	`

	rows, err := DB.Query(query, cardID)
	if err != nil {
		return nil, fmt.Errorf("error querying card effects: %v", err)
	}

	return rows, nil
}

// GetCardsByEffectID retrieves all cards that have a specific effect
func GetCardsByEffectID(effectID int) (*sql.Rows, error) {
	query := `
		SELECT c.id, c.name, c.type, c.legend, c.element, c.created_at, c.updated_at
		FROM cards c
		JOIN card_effects ce ON c.id = ce.card_id
		WHERE ce.effect_id = ?
		ORDER BY c.id
	`

	rows, err := DB.Query(query, effectID)
	if err != nil {
		return nil, fmt.Errorf("error querying effect cards: %v", err)
	}

	return rows, nil
}

// DeleteCardEffect removes a relationship between a card and an effect
func DeleteCardEffect(cardID, effectID int) error {
	query := `
		DELETE FROM card_effects 
		WHERE card_id = ? AND effect_id = ?
	`

	_, err := DB.Exec(query, cardID, effectID)
	if err != nil {
		return fmt.Errorf("error deleting card effect: %v", err)
	}

	return nil
}

// DeleteAllCardEffects removes all effects for a specific card
func DeleteAllCardEffects(cardID int) error {
	query := `
		DELETE FROM card_effects 
		WHERE card_id = ?
	`

	_, err := DB.Exec(query, cardID)
	if err != nil {
		return fmt.Errorf("error deleting all card effects: %v", err)
	}

	return nil
}
