package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"tcg-server-go/database"
	"tcg-server-go/models"
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
