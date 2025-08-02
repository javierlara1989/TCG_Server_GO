package database

import (
	"database/sql"
	"tcg-server-go/models"
	"time"
)

// CreateCard creates a new card record in the database
func CreateCard(card *models.Card) error {
	query := `
		INSERT INTO cards (name, type, legend, element, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`

	now := time.Now()
	card.CreatedAt = now
	card.UpdatedAt = now

	result, err := DB.Exec(query, card.Name, card.Type, card.Legend, card.Element, card.CreatedAt, card.UpdatedAt)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	card.ID = int(id)
	return nil
}

// GetCardByID retrieves a card by its ID
func GetCardByID(id int) (*models.Card, error) {
	query := `
		SELECT id, name, type, legend, element, created_at, updated_at
		FROM cards
		WHERE id = ?
	`

	card := &models.Card{}
	err := DB.QueryRow(query, id).Scan(
		&card.ID,
		&card.Name,
		&card.Type,
		&card.Legend,
		&card.Element,
		&card.CreatedAt,
		&card.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Card not found
		}
		return nil, err
	}

	return card, nil
}

// GetCardByName retrieves a card by its name
func GetCardByName(name string) (*models.Card, error) {
	query := `
		SELECT id, name, type, legend, element, created_at, updated_at
		FROM cards
		WHERE name = ?
	`

	card := &models.Card{}
	err := DB.QueryRow(query, name).Scan(
		&card.ID,
		&card.Name,
		&card.Type,
		&card.Legend,
		&card.Element,
		&card.CreatedAt,
		&card.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Card not found
		}
		return nil, err
	}

	return card, nil
}

// GetAllCards retrieves all cards from the database
func GetAllCards() ([]*models.Card, error) {
	query := `
		SELECT id, name, type, legend, element, created_at, updated_at
		FROM cards
		ORDER BY id
	`

	rows, err := DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cards []*models.Card
	for rows.Next() {
		card := &models.Card{}
		err := rows.Scan(
			&card.ID,
			&card.Name,
			&card.Type,
			&card.Legend,
			&card.Element,
			&card.CreatedAt,
			&card.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		cards = append(cards, card)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return cards, nil
}

// GetCardsByType retrieves all cards of a specific type
func GetCardsByType(cardType models.CardType) ([]*models.Card, error) {
	query := `
		SELECT id, name, type, legend, element, created_at, updated_at
		FROM cards
		WHERE type = ?
		ORDER BY id
	`

	rows, err := DB.Query(query, cardType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cards []*models.Card
	for rows.Next() {
		card := &models.Card{}
		err := rows.Scan(
			&card.ID,
			&card.Name,
			&card.Type,
			&card.Legend,
			&card.Element,
			&card.CreatedAt,
			&card.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		cards = append(cards, card)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return cards, nil
}

// GetCardsByElement retrieves all cards of a specific element
func GetCardsByElement(element models.CardElement) ([]*models.Card, error) {
	query := `
		SELECT id, name, type, legend, element, created_at, updated_at
		FROM cards
		WHERE element = ?
		ORDER BY id
	`

	rows, err := DB.Query(query, element)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cards []*models.Card
	for rows.Next() {
		card := &models.Card{}
		err := rows.Scan(
			&card.ID,
			&card.Name,
			&card.Type,
			&card.Legend,
			&card.Element,
			&card.CreatedAt,
			&card.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		cards = append(cards, card)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return cards, nil
}

// UpdateCard updates a card
func UpdateCard(card *models.Card) error {
	query := `
		UPDATE cards
		SET name = ?, type = ?, legend = ?, element = ?, updated_at = ?
		WHERE id = ?
	`

	card.UpdatedAt = time.Now()

	result, err := DB.Exec(query, card.Name, card.Type, card.Legend, card.Element, card.UpdatedAt, card.ID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// UpdateCardPartial updates specific fields of a card
func UpdateCardPartial(id int, req *models.UpdateCardRequest) error {
	// Build dynamic query based on provided fields
	query := "UPDATE cards SET updated_at = ?"
	args := []interface{}{time.Now()}

	if req.Name != nil {
		query += ", name = ?"
		args = append(args, *req.Name)
	}

	if req.Type != nil {
		query += ", type = ?"
		args = append(args, *req.Type)
	}

	if req.Legend != nil {
		query += ", legend = ?"
		args = append(args, *req.Legend)
	}

	if req.Element != nil {
		query += ", element = ?"
		args = append(args, *req.Element)
	}

	query += " WHERE id = ?"
	args = append(args, id)

	result, err := DB.Exec(query, args...)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// DeleteCard deletes a card by ID
func DeleteCard(id int) error {
	query := `DELETE FROM cards WHERE id = ?`

	result, err := DB.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// CardExists checks if a card exists by ID
func CardExists(id int) (bool, error) {
	query := `SELECT COUNT(*) FROM cards WHERE id = ?`

	var count int
	err := DB.QueryRow(query, id).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// CardNameExists checks if a card exists by name
func CardNameExists(name string) (bool, error) {
	query := `SELECT COUNT(*) FROM cards WHERE name = ?`

	var count int
	err := DB.QueryRow(query, name).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// SearchCards searches for cards by name (partial match)
func SearchCards(searchTerm string) ([]*models.Card, error) {
	query := `
		SELECT id, name, type, legend, element, created_at, updated_at
		FROM cards
		WHERE name LIKE ?
		ORDER BY name
	`

	searchPattern := "%" + searchTerm + "%"
	rows, err := DB.Query(query, searchPattern)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cards []*models.Card
	for rows.Next() {
		card := &models.Card{}
		err := rows.Scan(
			&card.ID,
			&card.Name,
			&card.Type,
			&card.Legend,
			&card.Element,
			&card.CreatedAt,
			&card.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		cards = append(cards, card)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return cards, nil
}
