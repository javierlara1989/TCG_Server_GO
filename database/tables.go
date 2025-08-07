package database

import (
	"database/sql"
	"fmt"
)

// CreateTable creates a new table
func CreateTable(category, privacy, prize string, password *string, amount *int) (*sql.Result, error) {
	query := `
		INSERT INTO tables (category, privacy, password, prize, amount, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, NOW(), NOW())
	`

	var result sql.Result
	var err error

	if password != nil && amount != nil {
		result, err = DB.Exec(query, category, privacy, *password, prize, *amount)
	} else if password != nil {
		result, err = DB.Exec(query, category, privacy, *password, prize, nil)
	} else if amount != nil {
		result, err = DB.Exec(query, category, privacy, nil, prize, *amount)
	} else {
		result, err = DB.Exec(query, category, privacy, nil, prize, nil)
	}

	if err != nil {
		return nil, fmt.Errorf("error creating table: %v", err)
	}

	return &result, nil
}

// CreateUserTable creates a new user table association
func CreateUserTable(userID, tableID uint, rivalID *uint) error {
	query := `
		INSERT INTO user_tables (user_id, rival_id, table_id, time)
		VALUES (?, ?, ?, 0)
	`

	_, err := DB.Exec(query, userID, rivalID, tableID)
	if err != nil {
		return fmt.Errorf("error creating user table: %v", err)
	}

	return nil
}

// GetTableByID retrieves a table by its ID
func GetTableByID(id uint) (*sql.Row, error) {
	query := `
		SELECT id, category, privacy, password, prize, amount, winner, created_at, updated_at, finished_at
		FROM tables WHERE id = ?
	`

	row := DB.QueryRow(query, id)
	return row, nil
}

// GetUserTableByTableID retrieves user table by table ID
func GetUserTableByTableID(tableID uint) (*sql.Row, error) {
	query := `
		SELECT ut.id, ut.user_id, ut.rival_id, ut.table_id, ut.time,
		       u.name as user_name, u.email as user_email,
		       r.name as rival_name, r.email as rival_email,
		       t.category, t.privacy, t.prize, t.amount, t.winner, t.created_at, t.updated_at, t.finished_at
		FROM user_tables ut
		JOIN users u ON ut.user_id = u.id
		LEFT JOIN users r ON ut.rival_id = r.id
		JOIN tables t ON ut.table_id = t.id
		WHERE ut.table_id = ?
	`

	row := DB.QueryRow(query, tableID)
	return row, nil
}

// GetUserTablesByUserID retrieves all user tables for a specific user
func GetUserTablesByUserID(userID uint) (*sql.Rows, error) {
	query := `
		SELECT ut.id, ut.user_id, ut.rival_id, ut.table_id, ut.time,
		       u.name as user_name, u.email as user_email,
		       r.name as rival_name, r.email as rival_email,
		       t.category, t.privacy, t.prize, t.amount, t.winner, t.created_at, t.updated_at, t.finished_at
		FROM user_tables ut
		JOIN users u ON ut.user_id = u.id
		LEFT JOIN users r ON ut.rival_id = r.id
		JOIN tables t ON ut.table_id = t.id
		WHERE ut.user_id = ? OR ut.rival_id = ?
	`

	rows, err := DB.Query(query, userID, userID)
	if err != nil {
		return nil, fmt.Errorf("error querying user tables: %v", err)
	}

	return rows, nil
}

// UpdateTable updates table fields
func UpdateTable(id uint, category, privacy, prize string, password *string, amount *int) error {
	query := `
		UPDATE tables 
		SET category = ?, privacy = ?, password = ?, prize = ?, amount = ?, updated_at = NOW()
		WHERE id = ?
	`

	var err error
	if password != nil && amount != nil {
		_, err = DB.Exec(query, category, privacy, *password, prize, *amount, id)
	} else if password != nil {
		_, err = DB.Exec(query, category, privacy, *password, prize, nil, id)
	} else if amount != nil {
		_, err = DB.Exec(query, category, privacy, nil, prize, *amount, id)
	} else {
		_, err = DB.Exec(query, category, privacy, nil, prize, nil, id)
	}

	if err != nil {
		return fmt.Errorf("error updating table: %v", err)
	}

	return nil
}

// IsTableOwner checks if a user is the owner of a table (user_id matches and rival_id is null)
func IsTableOwner(userID, tableID uint) (bool, error) {
	query := `
		SELECT COUNT(*) FROM user_tables 
		WHERE user_id = ? AND table_id = ? AND rival_id IS NULL
	`

	var count int
	err := DB.QueryRow(query, userID, tableID).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("error checking table ownership: %v", err)
	}

	return count > 0, nil
}

// IsTableWaitingForRival checks if a table is waiting for a rival (rival_id is null)
func IsTableWaitingForRival(tableID uint) (bool, error) {
	query := `
		SELECT COUNT(*) FROM user_tables 
		WHERE table_id = ? AND rival_id IS NULL
	`

	var count int
	err := DB.QueryRow(query, tableID).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("error checking if table is waiting for rival: %v", err)
	}

	return count > 0, nil
}

// UpdateUserTableTime updates the time field for a user table
func UpdateUserTableTime(userTableID uint, time int) error {
	query := `
		UPDATE user_tables 
		SET time = ?
		WHERE id = ?
	`

	_, err := DB.Exec(query, time, userTableID)
	if err != nil {
		return fmt.Errorf("error updating user table time: %v", err)
	}

	return nil
}

// DeleteTable deletes a table and its associated user table
func DeleteTable(tableID uint) error {
	// Start a transaction
	tx, err := DB.Begin()
	if err != nil {
		return fmt.Errorf("error starting transaction: %v", err)
	}

	// Delete associated user table first
	_, err = tx.Exec("DELETE FROM user_tables WHERE table_id = ?", tableID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("error deleting user table: %v", err)
	}

	// Delete the table
	_, err = tx.Exec("DELETE FROM tables WHERE id = ?", tableID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("error deleting table: %v", err)
	}

	return tx.Commit()
}
