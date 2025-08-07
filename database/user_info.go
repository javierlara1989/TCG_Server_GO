package database

import (
	"database/sql"
	"fmt"
	"tcg-server-go/models"
	"time"
)

// CreateUserInfo creates a new user info record in the database
func CreateUserInfo(userInfo *models.UserInfo) error {
	query := `
		INSERT INTO user_info (user_id, level, experience, money, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`

	now := time.Now()
	userInfo.CreatedAt = now
	userInfo.UpdatedAt = now

	result, err := DB.Exec(query, userInfo.UserID, userInfo.Level, userInfo.Experience, userInfo.Money, userInfo.CreatedAt, userInfo.UpdatedAt)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	userInfo.ID = int(id)
	return nil
}

// GetUserInfoByUserID retrieves user info by user ID
func GetUserInfoByUserID(userID int) (*models.UserInfo, error) {
	query := `
		SELECT id, user_id, level, experience, money, created_at, updated_at
		FROM user_info
		WHERE user_id = ?
	`

	userInfo := &models.UserInfo{}
	err := DB.QueryRow(query, userID).Scan(
		&userInfo.ID,
		&userInfo.UserID,
		&userInfo.Level,
		&userInfo.Experience,
		&userInfo.Money,
		&userInfo.CreatedAt,
		&userInfo.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // User info not found
		}
		return nil, err
	}

	return userInfo, nil
}

// GetUserInfoByID retrieves user info by its own ID
func GetUserInfoByID(id int) (*models.UserInfo, error) {
	query := `
		SELECT id, user_id, level, experience, money, created_at, updated_at
		FROM user_info
		WHERE id = ?
	`

	userInfo := &models.UserInfo{}
	err := DB.QueryRow(query, id).Scan(
		&userInfo.ID,
		&userInfo.UserID,
		&userInfo.Level,
		&userInfo.Experience,
		&userInfo.Money,
		&userInfo.CreatedAt,
		&userInfo.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // User info not found
		}
		return nil, err
	}

	return userInfo, nil
}

// UpdateUserInfo updates user info
func UpdateUserInfo(userInfo *models.UserInfo) error {
	query := `
		UPDATE user_info
		SET level = ?, experience = ?, money = ?, updated_at = ?
		WHERE id = ?
	`

	userInfo.UpdatedAt = time.Now()

	result, err := DB.Exec(query, userInfo.Level, userInfo.Experience, userInfo.Money, userInfo.UpdatedAt, userInfo.ID)
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

// UpdateUserInfoPartial updates specific fields of user info
func UpdateUserInfoPartial(userID int, req *models.UpdateUserInfoRequest) error {
	// Build dynamic query based on provided fields
	query := "UPDATE user_info SET updated_at = ?"
	args := []interface{}{time.Now()}

	if req.Level != nil {
		query += ", level = ?"
		args = append(args, *req.Level)
	}

	if req.Experience != nil {
		query += ", experience = ?"
		args = append(args, *req.Experience)
	}

	if req.Money != nil {
		query += ", money = ?"
		args = append(args, *req.Money)
	}

	query += " WHERE user_id = ?"
	args = append(args, userID)

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

// AddExperience adds experience points to a user and handles level up
func AddExperience(userID int, experienceToAdd int) (*models.UserInfo, error) {
	// Start a transaction
	tx, err := DB.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Get current user info
	query := `
		SELECT id, user_id, level, experience, money, created_at, updated_at
		FROM user_info
		WHERE user_id = ? FOR UPDATE
	`

	userInfo := &models.UserInfo{}
	err = tx.QueryRow(query, userID).Scan(
		&userInfo.ID,
		&userInfo.UserID,
		&userInfo.Level,
		&userInfo.Experience,
		&userInfo.Money,
		&userInfo.CreatedAt,
		&userInfo.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user info not found")
		}
		return nil, err
	}

	// Add experience
	userInfo.Experience += experienceToAdd
	userInfo.UpdatedAt = time.Now()

	// Check for level up (simple formula: 1000 exp per level)
	experienceForNextLevel := userInfo.Level * 1000
	if userInfo.Experience >= experienceForNextLevel {
		userInfo.Level++
		// Add bonus money for leveling up
		userInfo.Money += userInfo.Level * 100
	}

	// Update in database
	updateQuery := `
		UPDATE user_info
		SET level = ?, experience = ?, money = ?, updated_at = ?
		WHERE user_id = ?
	`

	_, err = tx.Exec(updateQuery, userInfo.Level, userInfo.Experience, userInfo.Money, userInfo.UpdatedAt, userID)
	if err != nil {
		return nil, err
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return userInfo, nil
}

// AddMoney adds money to a user's account
func AddMoney(userID int, amount int) (*models.UserInfo, error) {
	query := `
		UPDATE user_info
		SET money = money + ?, updated_at = ?
		WHERE user_id = ?
	`

	result, err := DB.Exec(query, amount, time.Now(), userID)
	if err != nil {
		return nil, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}

	if rowsAffected == 0 {
		return nil, sql.ErrNoRows
	}

	// Return updated user info
	return GetUserInfoByUserID(userID)
}

// SpendMoney spends money from a user's account (with validation)
func SpendMoney(userID int, amount int) (*models.UserInfo, error) {
	// Start a transaction
	tx, err := DB.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Get current user info with lock
	query := `
		SELECT id, user_id, level, experience, money, created_at, updated_at
		FROM user_info
		WHERE user_id = ? FOR UPDATE
	`

	userInfo := &models.UserInfo{}
	err = tx.QueryRow(query, userID).Scan(
		&userInfo.ID,
		&userInfo.UserID,
		&userInfo.Level,
		&userInfo.Experience,
		&userInfo.Money,
		&userInfo.CreatedAt,
		&userInfo.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user info not found")
		}
		return nil, err
	}

	// Check if user has enough money
	if userInfo.Money < amount {
		return nil, fmt.Errorf("insufficient funds: required %d, available %d", amount, userInfo.Money)
	}

	// Spend money
	userInfo.Money -= amount
	userInfo.UpdatedAt = time.Now()

	// Update in database
	updateQuery := `
		UPDATE user_info
		SET money = ?, updated_at = ?
		WHERE user_id = ?
	`

	_, err = tx.Exec(updateQuery, userInfo.Money, userInfo.UpdatedAt, userID)
	if err != nil {
		return nil, err
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return userInfo, nil
}

// DeleteUserInfo deletes user info by user ID
func DeleteUserInfo(userID int) error {
	query := `DELETE FROM user_info WHERE user_id = ?`

	result, err := DB.Exec(query, userID)
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

// UserInfoExists checks if user info exists for a user
func UserInfoExists(userID int) (bool, error) {
	query := `SELECT COUNT(*) FROM user_info WHERE user_id = ?`

	var count int
	err := DB.QueryRow(query, userID).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// CreateDefaultUserInfo creates default user info for a new user
func CreateDefaultUserInfo(userID int) (*models.UserInfo, error) {
	defaultUserInfo := &models.UserInfo{
		UserID:     userID,
		Level:      1,
		Experience: 0,
		Money:      100, // Starting money
	}

	return defaultUserInfo, CreateUserInfo(defaultUserInfo)
}

// User Cards Database Operations

// CreateUserCard creates a new user card record in the database
func CreateUserCard(userCard *models.UserCard) error {
	query := `
		INSERT INTO user_cards (user_id, card_id, amount, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?)
	`

	now := time.Now()
	userCard.CreatedAt = now
	userCard.UpdatedAt = now

	result, err := DB.Exec(query, userCard.UserID, userCard.CardID, userCard.Amount, userCard.CreatedAt, userCard.UpdatedAt)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	userCard.ID = int(id)
	return nil
}

// GetUserCardByID retrieves a user card by its ID
func GetUserCardByID(id int) (*models.UserCard, error) {
	query := `
		SELECT uc.id, uc.user_id, uc.card_id, uc.amount, uc.created_at, uc.updated_at,
		       c.id, c.name, c.type, c.legend, c.element, c.created_at, c.updated_at
		FROM user_cards uc
		JOIN cards c ON uc.card_id = c.id
		WHERE uc.id = ?
	`

	userCard := &models.UserCard{}
	card := &models.Card{}
	err := DB.QueryRow(query, id).Scan(
		&userCard.ID,
		&userCard.UserID,
		&userCard.CardID,
		&userCard.Amount,
		&userCard.CreatedAt,
		&userCard.UpdatedAt,
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
			return nil, nil // User card not found
		}
		return nil, err
	}

	userCard.Card = card
	return userCard, nil
}

// GetUserCardByUserAndCard retrieves a user card by user ID and card ID
func GetUserCardByUserAndCard(userID, cardID int) (*models.UserCard, error) {
	query := `
		SELECT uc.id, uc.user_id, uc.card_id, uc.amount, uc.created_at, uc.updated_at,
		       c.id, c.name, c.type, c.legend, c.element, c.created_at, c.updated_at
		FROM user_cards uc
		JOIN cards c ON uc.card_id = c.id
		WHERE uc.user_id = ? AND uc.card_id = ?
	`

	userCard := &models.UserCard{}
	card := &models.Card{}
	err := DB.QueryRow(query, userID, cardID).Scan(
		&userCard.ID,
		&userCard.UserID,
		&userCard.CardID,
		&userCard.Amount,
		&userCard.CreatedAt,
		&userCard.UpdatedAt,
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
			return nil, nil // User card not found
		}
		return nil, err
	}

	userCard.Card = card
	return userCard, nil
}

// GetUserCardsByUserID retrieves all cards for a specific user
func GetUserCardsByUserID(userID int) ([]models.UserCard, error) {
	query := `
		SELECT uc.id, uc.user_id, uc.card_id, uc.amount, uc.created_at, uc.updated_at,
		       c.id, c.name, c.type, c.legend, c.element, c.created_at, c.updated_at
		FROM user_cards uc
		JOIN cards c ON uc.card_id = c.id
		WHERE uc.user_id = ?
		ORDER BY c.name
	`

	rows, err := DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var userCards []models.UserCard
	for rows.Next() {
		userCard := models.UserCard{}
		card := models.Card{}
		err := rows.Scan(
			&userCard.ID,
			&userCard.UserID,
			&userCard.CardID,
			&userCard.Amount,
			&userCard.CreatedAt,
			&userCard.UpdatedAt,
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
		userCard.Card = &card
		userCards = append(userCards, userCard)
	}

	return userCards, nil
}

// UpdateUserCard updates a user card's amount
func UpdateUserCard(userCard *models.UserCard) error {
	query := `
		UPDATE user_cards
		SET amount = ?, updated_at = ?
		WHERE id = ?
	`

	userCard.UpdatedAt = time.Now()
	_, err := DB.Exec(query, userCard.Amount, userCard.UpdatedAt, userCard.ID)
	return err
}

// AddOrUpdateUserCard adds a new user card or updates the amount if it already exists
func AddOrUpdateUserCard(userID, cardID, amount int) error {
	// First, try to get existing user card
	existingCard, err := GetUserCardByUserAndCard(userID, cardID)
	if err != nil {
		return err
	}

	if existingCard != nil {
		// Update existing card
		existingCard.Amount += amount
		return UpdateUserCard(existingCard)
	} else {
		// Create new user card
		userCard := &models.UserCard{
			UserID: userID,
			CardID: cardID,
			Amount: amount,
		}
		return CreateUserCard(userCard)
	}
}

// DeleteUserCard deletes a user card by ID
func DeleteUserCard(id int) error {
	query := `DELETE FROM user_cards WHERE id = ?`
	_, err := DB.Exec(query, id)
	return err
}

// DeleteUserCardByUserAndCard deletes a user card by user ID and card ID
func DeleteUserCardByUserAndCard(userID, cardID int) error {
	query := `DELETE FROM user_cards WHERE user_id = ? AND card_id = ?`
	_, err := DB.Exec(query, userID, cardID)
	return err
}
