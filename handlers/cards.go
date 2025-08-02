package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"tcg-server-go/database"
	"tcg-server-go/models"

	"github.com/gorilla/mux"
)

// GetAllCardsHandler retrieves all cards
func GetAllCardsHandler(w http.ResponseWriter, r *http.Request) {
	cards, err := database.GetAllCards()
	if err != nil {
		http.Error(w, "Error retrieving cards", http.StatusInternalServerError)
		return
	}

	response := models.CardsResponse{
		Cards:   cards,
		Message: "Cards retrieved successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// GetCardByIDHandler retrieves a card by ID
func GetCardByIDHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid card ID", http.StatusBadRequest)
		return
	}

	card, err := database.GetCardByID(id)
	if err != nil {
		http.Error(w, "Error retrieving card", http.StatusInternalServerError)
		return
	}

	if card == nil {
		http.Error(w, "Card not found", http.StatusNotFound)
		return
	}

	response := models.CardResponse{
		Card:    card,
		Message: "Card retrieved successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// GetCardsByTypeHandler retrieves all cards of a specific type
func GetCardsByTypeHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	cardType := models.CardType(vars["type"])

	// Validate card type
	if cardType != models.CardTypeMonster && cardType != models.CardTypeSpell && cardType != models.CardTypeEnergy {
		http.Error(w, "Invalid card type", http.StatusBadRequest)
		return
	}

	cards, err := database.GetCardsByType(cardType)
	if err != nil {
		http.Error(w, "Error retrieving cards", http.StatusInternalServerError)
		return
	}

	response := models.CardsResponse{
		Cards:   cards,
		Message: "Cards retrieved successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// GetCardsByElementHandler retrieves all cards of a specific element
func GetCardsByElementHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	element := models.CardElement(vars["element"])

	// Validate element
	validElements := []models.CardElement{
		models.CardElementFire, models.CardElementWater, models.CardElementWind,
		models.CardElementEarth, models.CardElementNeutral, models.CardElementHoly, models.CardElementDark,
	}

	isValid := false
	for _, validElement := range validElements {
		if element == validElement {
			isValid = true
			break
		}
	}

	if !isValid {
		http.Error(w, "Invalid element", http.StatusBadRequest)
		return
	}

	cards, err := database.GetCardsByElement(element)
	if err != nil {
		http.Error(w, "Error retrieving cards", http.StatusInternalServerError)
		return
	}

	response := models.CardsResponse{
		Cards:   cards,
		Message: "Cards retrieved successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// SearchCardsHandler searches for cards by name
func SearchCardsHandler(w http.ResponseWriter, r *http.Request) {
	searchTerm := r.URL.Query().Get("q")
	if searchTerm == "" {
		http.Error(w, "Search term is required", http.StatusBadRequest)
		return
	}

	cards, err := database.SearchCards(searchTerm)
	if err != nil {
		http.Error(w, "Error searching cards", http.StatusInternalServerError)
		return
	}

	response := models.CardsResponse{
		Cards:   cards,
		Message: "Cards found successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
