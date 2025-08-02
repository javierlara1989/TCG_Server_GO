package database

import (
	"database/sql"
	"tcg-server-go/models"
	"time"
)

// CreateUser creates a new user in the database
func CreateUser(user *models.User) error {
	query := `
		INSERT INTO users (nombre, email, password, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?)
	`

	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	result, err := DB.Exec(query, user.Nombre, user.Email, user.Password, user.CreatedAt, user.UpdatedAt)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	user.ID = int(id)
	return nil
}

// GetUserByEmail retrieves a user by email
func GetUserByEmail(email string) (*models.User, error) {
	query := `
		SELECT id, nombre, email, password, created_at, updated_at, deleted_at
		FROM users
		WHERE email = ? AND deleted_at IS NULL
	`

	user := &models.User{}
	err := DB.QueryRow(query, email).Scan(
		&user.ID,
		&user.Nombre,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.DeletedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // User not found
		}
		return nil, err
	}

	return user, nil
}

// GetUserByID retrieves a user by ID
func GetUserByID(id int) (*models.User, error) {
	query := `
		SELECT id, nombre, email, password, created_at, updated_at, deleted_at
		FROM users
		WHERE id = ? AND deleted_at IS NULL
	`

	user := &models.User{}
	err := DB.QueryRow(query, id).Scan(
		&user.ID,
		&user.Nombre,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.DeletedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // User not found
		}
		return nil, err
	}

	return user, nil
}

// UpdateUser updates user information
func UpdateUser(user *models.User) error {
	query := `
		UPDATE users
		SET nombre = ?, email = ?, updated_at = ?
		WHERE id = ? AND deleted_at IS NULL
	`

	user.UpdatedAt = time.Now()

	result, err := DB.Exec(query, user.Nombre, user.Email, user.UpdatedAt, user.ID)
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

// UpdatePassword updates user password
func UpdatePassword(userID int, hashedPassword string) error {
	query := `
		UPDATE users
		SET password = ?, updated_at = ?
		WHERE id = ? AND deleted_at IS NULL
	`

	result, err := DB.Exec(query, hashedPassword, time.Now(), userID)
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

// SoftDeleteUser marks a user as deleted (soft delete)
func SoftDeleteUser(userID int) error {
	query := `
		UPDATE users
		SET deleted_at = ?, updated_at = ?
		WHERE id = ? AND deleted_at IS NULL
	`

	now := time.Now()
	result, err := DB.Exec(query, now, now, userID)
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

// HardDeleteUser permanently deletes a user
func HardDeleteUser(userID int) error {
	query := `DELETE FROM users WHERE id = ?`

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

// EmailExists checks if an email already exists in the database
func EmailExists(email string) (bool, error) {
	query := `SELECT COUNT(*) FROM users WHERE email = ? AND deleted_at IS NULL`

	var count int
	err := DB.QueryRow(query, email).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}
