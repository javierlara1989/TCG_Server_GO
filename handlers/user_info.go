package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

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

// Deck Handlers

// GetDecksHandler retrieves all decks for the authenticated user
func GetDecksHandler(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context (set by auth middleware)
	userID := r.Context().Value("user_id").(int)

	decks, err := database.GetDecksByUserID(userID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error retrieving decks: %v", err), http.StatusInternalServerError)
		return
	}

	response := models.DecksResponse{
		Decks:   decks,
		Message: "Decks retrieved successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// GetDeckHandler retrieves a specific deck by ID
func GetDeckHandler(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context (set by auth middleware)
	userID := r.Context().Value("user_id").(int)

	// Get deck ID from URL parameters
	vars := mux.Vars(r)
	deckIDStr := vars["id"]
	deckID, err := strconv.Atoi(deckIDStr)
	if err != nil {
		http.Error(w, "Invalid deck ID", http.StatusBadRequest)
		return
	}

	deck, err := database.GetDeckByID(deckID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error retrieving deck: %v", err), http.StatusInternalServerError)
		return
	}

	if deck == nil {
		http.Error(w, "Deck not found", http.StatusNotFound)
		return
	}

	// Check if the deck belongs to the authenticated user
	if deck.UserID != userID {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	response := models.DeckResponse{
		Deck:    deck,
		Message: "Deck retrieved successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// GetDeckWithCardsHandler retrieves a deck with all its cards
func GetDeckWithCardsHandler(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context (set by auth middleware)
	userID := r.Context().Value("user_id").(int)

	// Get deck ID from URL parameters
	vars := mux.Vars(r)
	deckIDStr := vars["id"]
	deckID, err := strconv.Atoi(deckIDStr)
	if err != nil {
		http.Error(w, "Invalid deck ID", http.StatusBadRequest)
		return
	}

	deck, err := database.GetDeckByID(deckID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error retrieving deck: %v", err), http.StatusInternalServerError)
		return
	}

	if deck == nil {
		http.Error(w, "Deck not found", http.StatusNotFound)
		return
	}

	// Check if the deck belongs to the authenticated user
	if deck.UserID != userID {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	// Get deck cards
	deckCards, err := database.GetDeckCards(deckID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error retrieving deck cards: %v", err), http.StatusInternalServerError)
		return
	}

	deckWithCards := &models.DeckWithCards{
		Deck:  deck,
		Cards: deckCards,
	}

	response := models.DeckWithCardsResponse{
		DeckWithCards: deckWithCards,
		Message:       "Deck with cards retrieved successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// CreateDeckHandler creates a new deck with validation
func CreateDeckHandler(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context (set by auth middleware)
	userID := r.Context().Value("user_id").(int)

	var req models.CreateDeckRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate request
	if err := validate.Struct(req); err != nil {
		http.Error(w, fmt.Sprintf("Validation error: %v", err), http.StatusBadRequest)
		return
	}

	// Validate that card_ids and card_count arrays have the same length
	if len(req.CardIDs) != len(req.CardCount) {
		http.Error(w, "card_ids and card_count arrays must have the same length", http.StatusBadRequest)
		return
	}

	// Create deck with validation
	deck, err := database.CreateDeckWithValidation(userID, req.Name, req.CardIDs, req.CardCount)
	if err != nil {
		if err.Error() == "user does not have all required cards" {
			http.Error(w, "Cannot create deck: you do not have all the required cards", http.StatusBadRequest)
			return
		}
		if strings.Contains(err.Error(), "deck must have at least 40 cards") {
			http.Error(w, "Cannot create deck: deck must have at least 40 cards", http.StatusBadRequest)
			return
		}
		if strings.Contains(err.Error(), "deck limit reached") {
			http.Error(w, fmt.Sprintf("Cannot create deck: %v", err.Error()), http.StatusBadRequest)
			return
		}
		http.Error(w, fmt.Sprintf("Error creating deck: %v", err), http.StatusInternalServerError)
		return
	}

	response := models.DeckResponse{
		Deck:    deck,
		Message: "Deck created successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// DeleteDeckHandler deletes a deck
func DeleteDeckHandler(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context (set by auth middleware)
	userID := r.Context().Value("user_id").(int)

	// Get deck ID from URL parameters
	vars := mux.Vars(r)
	deckIDStr := vars["id"]
	deckID, err := strconv.Atoi(deckIDStr)
	if err != nil {
		http.Error(w, "Invalid deck ID", http.StatusBadRequest)
		return
	}

	// Check if deck exists and belongs to user
	deck, err := database.GetDeckByID(deckID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error checking deck: %v", err), http.StatusInternalServerError)
		return
	}
	if deck == nil {
		http.Error(w, "Deck not found", http.StatusNotFound)
		return
	}
	if deck.UserID != userID {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	// Delete deck
	err = database.DeleteDeck(deckID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error deleting deck: %v", err), http.StatusInternalServerError)
		return
	}

	response := map[string]string{
		"message": "Deck deleted successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// GetDeckLimitHandler retrieves deck limit information for the authenticated user
func GetDeckLimitHandler(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context (set by auth middleware)
	userID := r.Context().Value("user_id").(int)

	// Get current decks
	currentDecks, err := database.GetDecksByUserID(userID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error retrieving decks: %v", err), http.StatusInternalServerError)
		return
	}

	// Get deck limit
	deckLimit, err := database.GetUserDeckLimit(userID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error calculating deck limit: %v", err), http.StatusInternalServerError)
		return
	}

	// Get user info for level
	userInfo, err := database.GetUserInfoByUserID(userID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error retrieving user info: %v", err), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"current_decks": len(currentDecks),
		"deck_limit":    deckLimit,
		"user_level":    userInfo.Level,
		"can_create":    len(currentDecks) < deckLimit,
		"message":       "Deck limit information retrieved successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
