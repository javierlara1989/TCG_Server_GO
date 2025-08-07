package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"tcg-server-go/database"

	"github.com/gorilla/mux"
)

// CreateTableRequest represents the request body for creating a table
type CreateTableRequest struct {
	Category string  `json:"category"`
	Privacy  string  `json:"privacy"`
	Password *string `json:"password,omitempty"`
	Prize    string  `json:"prize"`
	Amount   *int    `json:"amount,omitempty"`
}

// UpdateTableRequest represents the request body for updating a table
type UpdateTableRequest struct {
	Category string  `json:"category,omitempty"`
	Privacy  string  `json:"privacy,omitempty"`
	Password *string `json:"password,omitempty"`
	Prize    string  `json:"prize,omitempty"`
	Amount   *int    `json:"amount,omitempty"`
}

// TableResponse represents the response for table operations
type TableResponse struct {
	ID         uint    `json:"id"`
	Category   string  `json:"category"`
	Privacy    string  `json:"privacy"`
	Password   *string `json:"password,omitempty"`
	Prize      string  `json:"prize"`
	Amount     *int    `json:"amount,omitempty"`
	Winner     *bool   `json:"winner,omitempty"`
	CreatedAt  string  `json:"created_at"`
	UpdatedAt  string  `json:"updated_at"`
	FinishedAt *string `json:"finished_at,omitempty"`
}

// UserTableResponse represents the response for user table operations
type UserTableResponse struct {
	ID         uint          `json:"id"`
	UserID     uint          `json:"user_id"`
	RivalID    *uint         `json:"rival_id,omitempty"`
	TableID    uint          `json:"table_id"`
	Time       int           `json:"time"`
	Table      TableResponse `json:"table"`
	UserName   string        `json:"user_name"`
	UserEmail  string        `json:"user_email"`
	RivalName  *string       `json:"rival_name,omitempty"`
	RivalEmail *string       `json:"rival_email,omitempty"`
}

// CreateTable creates a new table and associates it with the logged-in user
func CreateTable(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context (set by auth middleware)
	userID, ok := r.Context().Value("user_id").(uint)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Parse request body
	var req CreateTableRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if req.Category == "" || req.Privacy == "" || req.Prize == "" {
		http.Error(w, "Category, privacy, and prize are required", http.StatusBadRequest)
		return
	}

	// Validate category
	validCategories := map[string]bool{"S": true, "A": true, "B": true, "C": true, "D": true}
	if !validCategories[req.Category] {
		http.Error(w, "Invalid category. Must be S, A, B, C, or D", http.StatusBadRequest)
		return
	}

	// Validate privacy
	if req.Privacy != "private" && req.Privacy != "public" {
		http.Error(w, "Invalid privacy. Must be 'private' or 'public'", http.StatusBadRequest)
		return
	}

	// Validate prize
	validPrizes := map[string]bool{"money": true, "card": true, "aura": true}
	if !validPrizes[req.Prize] {
		http.Error(w, "Invalid prize. Must be 'money', 'card', or 'aura'", http.StatusBadRequest)
		return
	}

	// Validate password if provided
	if req.Password != nil {
		if len(*req.Password) > 10 {
			http.Error(w, "Password must be 10 characters or less", http.StatusBadRequest)
			return
		}
		// Check if password contains only digits
		for _, char := range *req.Password {
			if char < '0' || char > '9' {
				http.Error(w, "Password must contain only numeric characters", http.StatusBadRequest)
				return
			}
		}
	}

	// Create table
	result, err := database.CreateTable(req.Category, req.Privacy, req.Prize, req.Password, req.Amount)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating table: %v", err), http.StatusInternalServerError)
		return
	}

	// Get the table ID
	tableID, err := (*result).LastInsertId()
	if err != nil {
		http.Error(w, "Error getting table ID", http.StatusInternalServerError)
		return
	}

	// Create user table association with rival_id as null
	err = database.CreateUserTable(uint(userID), uint(tableID), nil)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating user table association: %v", err), http.StatusInternalServerError)
		return
	}

	// Return success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":  "Table created successfully",
		"table_id": tableID,
	})
}

// UpdateTable updates table parameters (only if user is owner and table is waiting for rival)
func UpdateTable(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, ok := r.Context().Value("user_id").(uint)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get table ID from URL
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 3 {
		http.Error(w, "Table ID required", http.StatusBadRequest)
		return
	}

	tableIDStr := pathParts[len(pathParts)-1]
	tableID, err := strconv.ParseUint(tableIDStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid table ID", http.StatusBadRequest)
		return
	}

	// Check if user is the owner of the table
	isOwner, err := database.IsTableOwner(userID, uint(tableID))
	if err != nil {
		http.Error(w, fmt.Sprintf("Error checking table ownership: %v", err), http.StatusInternalServerError)
		return
	}

	if !isOwner {
		http.Error(w, "You can only update tables you own", http.StatusForbidden)
		return
	}

	// Check if table is waiting for rival
	isWaiting, err := database.IsTableWaitingForRival(uint(tableID))
	if err != nil {
		http.Error(w, fmt.Sprintf("Error checking table status: %v", err), http.StatusInternalServerError)
		return
	}

	if !isWaiting {
		http.Error(w, "Cannot update table that already has a rival", http.StatusForbidden)
		return
	}

	// Parse request body
	var req UpdateTableRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Get current table data to merge with updates
	row, err := database.GetTableByID(uint(tableID))
	if err != nil {
		http.Error(w, "Error retrieving table", http.StatusInternalServerError)
		return
	}

	var currentCategory, currentPrivacy, currentPrize string
	var currentPassword *string
	var currentAmount *int
	var currentWinner *bool
	var currentCreatedAt, currentUpdatedAt string
	var currentFinishedAt *string

	err = row.Scan(&tableID, &currentCategory, &currentPrivacy, &currentPassword, &currentPrize, &currentAmount,
		&currentWinner, &currentCreatedAt, &currentUpdatedAt, &currentFinishedAt)
	if err != nil {
		http.Error(w, "Error reading table data", http.StatusInternalServerError)
		return
	}

	// Merge updates with current values
	if req.Category != "" {
		validCategories := map[string]bool{"S": true, "A": true, "B": true, "C": true, "D": true}
		if !validCategories[req.Category] {
			http.Error(w, "Invalid category. Must be S, A, B, C, or D", http.StatusBadRequest)
			return
		}
		currentCategory = req.Category
	}

	if req.Privacy != "" {
		if req.Privacy != "private" && req.Privacy != "public" {
			http.Error(w, "Invalid privacy. Must be 'private' or 'public'", http.StatusBadRequest)
			return
		}
		currentPrivacy = req.Privacy
	}

	if req.Prize != "" {
		validPrizes := map[string]bool{"money": true, "card": true, "aura": true}
		if !validPrizes[req.Prize] {
			http.Error(w, "Invalid prize. Must be 'money', 'card', or 'aura'", http.StatusBadRequest)
			return
		}
		currentPrize = req.Prize
	}

	if req.Password != nil {
		if len(*req.Password) > 10 {
			http.Error(w, "Password must be 10 characters or less", http.StatusBadRequest)
			return
		}
		// Check if password contains only digits
		for _, char := range *req.Password {
			if char < '0' || char > '9' {
				http.Error(w, "Password must contain only numeric characters", http.StatusBadRequest)
				return
			}
		}
		currentPassword = req.Password
	}

	if req.Amount != nil {
		if *req.Amount < 0 {
			http.Error(w, "Amount must be a positive number", http.StatusBadRequest)
			return
		}
		currentAmount = req.Amount
	}

	// Update table
	err = database.UpdateTable(uint(tableID), currentCategory, currentPrivacy, currentPrize, currentPassword, currentAmount)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error updating table: %v", err), http.StatusInternalServerError)
		return
	}

	// Return success response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":  "Table updated successfully",
		"table_id": tableID,
	})
}

// GetUserTables retrieves all tables for the logged-in user
func GetUserTables(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, ok := r.Context().Value("user_id").(uint)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get user tables
	rows, err := database.GetUserTablesByUserID(userID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error retrieving user tables: %v", err), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var userTables []UserTableResponse
	for rows.Next() {
		var ut UserTableResponse
		var rivalName, rivalEmail sql.NullString

		err := rows.Scan(
			&ut.ID, &ut.UserID, &ut.RivalID, &ut.TableID, &ut.Time,
			&ut.UserName, &ut.UserEmail,
			&rivalName, &rivalEmail,
			&ut.Table.Category, &ut.Table.Privacy, &ut.Table.Prize, &ut.Table.Amount, &ut.Table.Winner,
			&ut.Table.CreatedAt, &ut.Table.UpdatedAt, &ut.Table.FinishedAt,
		)
		if err != nil {
			http.Error(w, "Error reading table data", http.StatusInternalServerError)
			return
		}

		// Handle nullable rival fields
		if rivalName.Valid {
			ut.RivalName = &rivalName.String
		}
		if rivalEmail.Valid {
			ut.RivalEmail = &rivalEmail.String
		}

		userTables = append(userTables, ut)
	}

	// Return response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"tables": userTables,
	})
}

// UpdateUserTableTimeRequest represents the request for updating user table time
type UpdateUserTableTimeRequest struct {
	Time int `json:"time"`
}

// UpdateUserTableTime updates the time field for a user table
func UpdateUserTableTime(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, ok := r.Context().Value("user_id").(uint)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get table ID from URL parameters
	vars := mux.Vars(r)
	tableIDStr, ok := vars["id"]
	if !ok {
		http.Error(w, "Table ID is required", http.StatusBadRequest)
		return
	}

	tableID, err := strconv.ParseUint(tableIDStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid table ID", http.StatusBadRequest)
		return
	}

	// Parse request body
	var req UpdateUserTableTimeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate time value
	if req.Time < 0 {
		http.Error(w, "Time must be a non-negative number", http.StatusBadRequest)
		return
	}

	// Check if user is associated with this table
	isOwner, err := database.IsTableOwner(userID, uint(tableID))
	if err != nil {
		http.Error(w, fmt.Sprintf("Error checking table ownership: %v", err), http.StatusInternalServerError)
		return
	}

	if !isOwner {
		http.Error(w, "You can only update time for your own tables", http.StatusForbidden)
		return
	}

	// Get the user table ID (we need to find the user table record)
	// For now, we'll use a simple approach - you might want to add a function to get user table by user and table IDs
	rows, err := database.GetUserTablesByUserID(userID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error retrieving user tables: %v", err), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var userTableID uint
	found := false
	for rows.Next() {
		var ut UserTableResponse
		var rivalName, rivalEmail sql.NullString

		err := rows.Scan(
			&ut.ID, &ut.UserID, &ut.RivalID, &ut.TableID, &ut.Time,
			&ut.UserName, &ut.UserEmail,
			&rivalName, &rivalEmail,
			&ut.Table.Category, &ut.Table.Privacy, &ut.Table.Prize, &ut.Table.Amount, &ut.Table.Winner,
			&ut.Table.CreatedAt, &ut.Table.UpdatedAt, &ut.Table.FinishedAt,
		)
		if err != nil {
			http.Error(w, "Error reading table data", http.StatusInternalServerError)
			return
		}

		if ut.TableID == uint(tableID) {
			userTableID = ut.ID
			found = true
			break
		}
	}

	if !found {
		http.Error(w, "Table not found or you don't have access to it", http.StatusNotFound)
		return
	}

	// Update the time
	err = database.UpdateUserTableTime(userTableID, req.Time)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error updating table time: %v", err), http.StatusInternalServerError)
		return
	}

	// Return success response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Table time updated successfully",
		"time":    req.Time,
	})
}
