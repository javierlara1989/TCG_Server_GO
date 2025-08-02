package handlers

import (
	"encoding/json"
	"net/http"

	"tcg-server-go/models"
)

func ValidateTokenHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.ValidateResponse{Message: "OK"})
} 