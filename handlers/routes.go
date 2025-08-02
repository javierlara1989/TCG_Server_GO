package handlers

import (
	"github.com/gorilla/mux"
	"tcg-server-go/middleware"
)

func SetupRoutes() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/login", LoginHandler).Methods("POST")
	r.HandleFunc("/health", HealthHandler).Methods("GET")

	protected := r.PathPrefix("/api").Subrouter()
	protected.Use(middleware.AuthMiddleware)
	protected.HandleFunc("/validate", ValidateTokenHandler).Methods("GET")

	return r
} 