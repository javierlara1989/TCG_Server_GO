package handlers

import (
	"encoding/json"
	"net/http"

	"tcg-server-go/auth"
	"tcg-server-go/models"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var loginReq models.LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&loginReq); err != nil {
		http.Error(w, "Error decoding request", http.StatusBadRequest)
		return
	}

	// Validate the request
	validationErrors := ValidateLoginRequest(&loginReq)
	if len(validationErrors) > 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ValidationResponse{Errors: validationErrors})
		return
	}

	if !auth.ValidateCredentials(loginReq.Email, loginReq.Password) {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	token, err := auth.GenerateToken(loginReq.Email)
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.LoginResponse{Token: token})
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var createReq models.CreateUserRequest

	if err := json.NewDecoder(r.Body).Decode(&createReq); err != nil {
		http.Error(w, "Error decoding request", http.StatusBadRequest)
		return
	}

	// Validate the request
	validationErrors := ValidateCreateUserRequest(&createReq)
	if len(validationErrors) > 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ValidationResponse{Errors: validationErrors})
		return
	}

	// Check if user already exists
	if auth.UserExists(createReq.Email) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(map[string]string{"error": "User already exists"})
		return
	}

	// Create the user
	user, err := auth.CreateUser(&createReq)
	if err != nil {
		http.Error(w, "Error creating user", http.StatusInternalServerError)
		return
	}

	// Generate token for the new user
	token, err := auth.GenerateToken(user.Email)
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(models.LoginResponse{Token: token})
}
