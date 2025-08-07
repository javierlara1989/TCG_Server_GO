package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"tcg-server-go/database"
	"tcg-server-go/models"

	"github.com/gorilla/mux"
)

// GetUserInfoHandler retrieves user info for the authenticated user
func GetUserInfoHandler(w http.ResponseWriter, r *http.Request) {
	// Get user ID from middleware
	userIDStr := r.Header.Get("X-User-ID")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	userInfo, err := database.GetUserInfoByUserID(userID)
	if err != nil {
		http.Error(w, "Error retrieving user info", http.StatusInternalServerError)
		return
	}

	if userInfo == nil {
		// Create default user info if it doesn't exist
		userInfo, err = database.CreateDefaultUserInfo(userID)
		if err != nil {
			http.Error(w, "Error creating user info", http.StatusInternalServerError)
			return
		}
	}

	response := models.UserInfoResponse{
		UserInfo: userInfo,
		Message:  "User info retrieved successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// User Cards Handlers

// GetUserCardsHandler retrieves all cards for the authenticated user
func GetUserCardsHandler(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context (set by auth middleware)
	userID := r.Context().Value("user_id").(int)

	userCards, err := database.GetUserCardsByUserID(userID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error retrieving user cards: %v", err), http.StatusInternalServerError)
		return
	}

	response := models.UserCardsResponse{
		UserCards: userCards,
		Message:   "User cards retrieved successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// GetUserCardHandler retrieves a specific user card by ID
func GetUserCardHandler(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context (set by auth middleware)
	userID := r.Context().Value("user_id").(int)

	// Get card ID from URL parameters
	vars := mux.Vars(r)
	cardIDStr := vars["id"]
	cardID, err := strconv.Atoi(cardIDStr)
	if err != nil {
		http.Error(w, "Invalid card ID", http.StatusBadRequest)
		return
	}

	userCard, err := database.GetUserCardByUserAndCard(userID, cardID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error retrieving user card: %v", err), http.StatusInternalServerError)
		return
	}

	if userCard == nil {
		http.Error(w, "User card not found", http.StatusNotFound)
		return
	}

	response := models.UserCardResponse{
		UserCard: userCard,
		Message:  "User card retrieved successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
