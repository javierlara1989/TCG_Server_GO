package auth

import (
	"golang.org/x/crypto/bcrypt"
	"tcg-server-go/database"
	"tcg-server-go/models"
)

// Sample users (fallback for testing when database is not available)
var users = map[string]string{
	"admin@example.com": "$2a$10$N.zmdr9k7uOCQb376NoUnuTJ8iAt6Z5EHsM8lE9lBOsl7iKTVEFDa", // password: "admin123"
	"user@example.com":  "$2a$10$8K1p/a0dL1LXMIgoEDFrwOfgqwAGcwZQh3UPHz3UaCgHpVqKqKqKq", // password: "user123"
}

func ValidateCredentials(email, password string) bool {
	// Try database first
	if database.DB != nil {
		user, err := database.GetUserByEmail(email)
		if err != nil || user == nil {
			return false
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
			return false
		}

		return true
	}

	// Fallback to in-memory users for testing
	hashedPassword, exists := users[email]
	if !exists {
		return false
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); err != nil {
		return false
	}

	return true
}

func UserExists(email string) bool {
	// Try database first
	if database.DB != nil {
		exists, err := database.EmailExists(email)
		if err != nil {
			return false
		}
		return exists
	}

	// Fallback to in-memory users for testing
	_, exists := users[email]
	return exists
}

// AddUser adds a new user to database or fallback to in-memory
func AddUser(user *models.User) error {
	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)

	// Try database first
	if database.DB != nil {
		return database.CreateUser(user)
	}

	// Fallback to in-memory users for testing
	users[user.Email] = user.Password
	return nil
}

// CreateUser creates a new user with validation
func CreateUser(req *models.CreateUserRequest) (*models.User, error) {
	user := &models.User{
		Nombre:   req.Nombre,
		Email:    req.Email,
		Password: req.Password,
	}

	return user, AddUser(user)
} 