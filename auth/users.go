package auth

import (
	"golang.org/x/crypto/bcrypt"
)

// Sample users (use database in production)
var users = map[string]string{
	"admin": "$2a$10$N.zmdr9k7uOCQb376NoUnuTJ8iAt6Z5EHsM8lE9lBOsl7iKTVEFDa", // password: "admin123"
	"user":  "$2a$10$8K1p/a0dL1LXMIgoEDFrwOfgqwAGcwZQh3UPHz3UaCgHpVqKqKqKq", // password: "user123"
}

func ValidateCredentials(username, password string) bool {
	hashedPassword, exists := users[username]
	if !exists {
		return false
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); err != nil {
		return false
	}

	return true
}

func UserExists(username string) bool {
	_, exists := users[username]
	return exists
}

// AddUser adds a new user (useful for testing)
func AddUser(username, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	users[username] = string(hashedPassword)
	return nil
} 