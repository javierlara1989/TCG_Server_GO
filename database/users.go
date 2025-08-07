package database

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"fmt"
	"strings"
	"tcg-server-go/models"
	"time"
)

// CreateUser creates a new user in the database
func CreateUser(user *models.User) error {
	// Generate validation code
	validationCode := generateValidationCode()
	expiresAt := time.Now().Add(24 * time.Hour) // Code expires in 24 hours

	user.ValidationCode = &validationCode
	user.ValidationCodeExpiresAt = &expiresAt

	query := `
		INSERT INTO users (name, email, password, validation_code, validation_code_expires_at, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	result, err := DB.Exec(query, user.Name, user.Email, user.Password, user.ValidationCode, user.ValidationCodeExpiresAt, user.CreatedAt, user.UpdatedAt)
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
		SELECT id, name, email, password, validation_code, validation_code_expires_at, validated_at, created_at, updated_at, deleted_at
		FROM users
		WHERE email = ? AND deleted_at IS NULL
	`

	user := &models.User{}
	err := DB.QueryRow(query, email).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.ValidationCode,
		&user.ValidationCodeExpiresAt,
		&user.ValidatedAt,
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
		SELECT id, name, email, password, validation_code, validation_code_expires_at, validated_at, created_at, updated_at, deleted_at
		FROM users
		WHERE id = ? AND deleted_at IS NULL
	`

	user := &models.User{}
	err := DB.QueryRow(query, id).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.ValidationCode,
		&user.ValidationCodeExpiresAt,
		&user.ValidatedAt,
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

// VerifyEmail verifies a user's email with the provided validation code
func VerifyEmail(email, validationCode string) (*models.User, error) {
	user, err := GetUserByEmail(email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, fmt.Errorf("user not found")
	}

	// Check if already validated
	if user.ValidatedAt != nil {
		return nil, fmt.Errorf("email already verified")
	}

	// Check if validation code matches
	if user.ValidationCode == nil || *user.ValidationCode != validationCode {
		return nil, fmt.Errorf("invalid validation code")
	}

	// Check if validation code has expired
	if user.ValidationCodeExpiresAt != nil && time.Now().After(*user.ValidationCodeExpiresAt) {
		return nil, fmt.Errorf("validation code has expired")
	}

	// Mark email as verified
	now := time.Now()
	query := `
		UPDATE users
		SET validated_at = ?, updated_at = ?, validation_code = NULL, validation_code_expires_at = NULL
		WHERE id = ? AND deleted_at IS NULL
	`

	result, err := DB.Exec(query, now, now, user.ID)
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

	user.ValidatedAt = &now
	user.UpdatedAt = now
	user.ValidationCode = nil
	user.ValidationCodeExpiresAt = nil

	return user, nil
}

// ResendValidationCode generates a new validation code for a user
func ResendValidationCode(email string) error {
	user, err := GetUserByEmail(email)
	if err != nil {
		return err
	}
	if user == nil {
		return fmt.Errorf("user not found")
	}

	// Check if already validated
	if user.ValidatedAt != nil {
		return fmt.Errorf("email already verified")
	}

	// Generate new validation code
	validationCode := generateValidationCode()
	expiresAt := time.Now().Add(24 * time.Hour)

	query := `
		UPDATE users
		SET validation_code = ?, validation_code_expires_at = ?, updated_at = ?
		WHERE id = ? AND deleted_at IS NULL
	`

	result, err := DB.Exec(query, validationCode, expiresAt, time.Now(), user.ID)
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

// generateValidationCode generates a random 6-character validation code
func generateValidationCode() string {
	bytes := make([]byte, 3)
	rand.Read(bytes)
	return strings.ToUpper(hex.EncodeToString(bytes)[:6])
}

// UpdateUser updates user information
func UpdateUser(user *models.User) error {
	query := `
		UPDATE users
		SET name = ?, email = ?, updated_at = ?
		WHERE id = ? AND deleted_at IS NULL
	`

	user.UpdatedAt = time.Now()

	result, err := DB.Exec(query, user.Name, user.Email, user.UpdatedAt, user.ID)
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
