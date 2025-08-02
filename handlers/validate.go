package handlers

import (
	"encoding/json"
	"net/http"
	"regexp"
	"strings"
	"unicode"

	"github.com/go-playground/validator/v10"
	"tcg-server-go/models"
)

// ValidationError represents a validation error
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ValidationResponse represents the response for validation errors
type ValidationResponse struct {
	Errors []ValidationError `json:"errors"`
}

// Custom validator instance
var validate = validator.New()

// init registers custom validation functions
func init() {
	// Register custom validation for nombre (only letters, minimum 6 characters)
	validate.RegisterValidation("nombre", validateNombre)
	
	// Register custom validation for password (letters and numbers, minimum 6 characters)
	validate.RegisterValidation("password", validatePassword)
}

// validateNombre validates that nombre contains only letters and is at least 6 characters
func validateNombre(fl validator.FieldLevel) bool {
	nombre := fl.Field().String()
	
	// Check minimum length
	if len(nombre) < 6 {
		return false
	}
	
	// Check if contains only letters and spaces
	letterRegex := regexp.MustCompile(`^[a-zA-ZáéíóúÁÉÍÓÚñÑ\s]+$`)
	return letterRegex.MatchString(nombre)
}

// validatePassword validates that password contains both letters and numbers, minimum 6 characters
func validatePassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	
	// Check minimum length
	if len(password) < 6 {
		return false
	}
	
	// Check if contains at least one letter and one number
	hasLetter := false
	hasNumber := false
	
	for _, char := range password {
		if unicode.IsLetter(char) {
			hasLetter = true
		} else if unicode.IsNumber(char) {
			hasNumber = true
		}
	}
	
	return hasLetter && hasNumber
}

// ValidateStruct validates a struct and returns validation errors
func ValidateStruct(s interface{}) []ValidationError {
	err := validate.Struct(s)
	if err == nil {
		return nil
	}

	var errors []ValidationError
	
	for _, err := range err.(validator.ValidationErrors) {
		field := strings.ToLower(err.Field())
		var message string
		
		switch err.Tag() {
		case "required":
			message = field + " is required"
		case "email":
			message = "Invalid email format"
		case "min":
			message = field + " must be at least " + err.Param() + " characters"
		case "nombre":
			message = "Nombre must be at least 6 characters and contain only letters"
		case "password":
			message = "Password must be at least 6 characters and contain both letters and numbers"
		case "alphanum":
			message = field + " must contain only letters and numbers"
		case "alpha":
			message = field + " must contain only letters"
		default:
			message = field + " is invalid"
		}
		
		errors = append(errors, ValidationError{
			Field:   field,
			Message: message,
		})
	}
	
	return errors
}

// ValidateLoginRequest validates login request
func ValidateLoginRequest(req *models.LoginRequest) []ValidationError {
	return ValidateStruct(req)
}

// ValidateCreateUserRequest validates create user request
func ValidateCreateUserRequest(req *models.CreateUserRequest) []ValidationError {
	return ValidateStruct(req)
}

// ValidateTokenHandler handles token validation
func ValidateTokenHandler(w http.ResponseWriter, r *http.Request) {
	response := models.ValidateResponse{
		Message: "OK",
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
} 