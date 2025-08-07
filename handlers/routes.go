package handlers

import (
	"tcg-server-go/middleware"

	"github.com/gorilla/mux"
)

func SetupRoutes() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/login", LoginHandler).Methods("POST")
	r.HandleFunc("/register", RegisterHandler).Methods("POST")
	r.HandleFunc("/verify-email", VerifyEmailHandler).Methods("POST")
	r.HandleFunc("/resend-code", ResendCodeHandler).Methods("POST")
	r.HandleFunc("/health", HealthHandler).Methods("GET")

	protected := r.PathPrefix("/api").Subrouter()
	protected.Use(middleware.AuthMiddleware)
	protected.HandleFunc("/validate", ValidateTokenHandler).Methods("GET")

	// User Info endpoint (read-only, requires authentication)
	protected.HandleFunc("/user-info", GetUserInfoHandler).Methods("GET")

	// User Cards endpoints (requires authentication)
	protected.HandleFunc("/user-cards", GetUserCardsHandler).Methods("GET")
	protected.HandleFunc("/user-cards/{id}", GetUserCardHandler).Methods("GET")

	// Card endpoints (public access for reading only)
	r.HandleFunc("/cards", GetAllCardsHandler).Methods("GET")
	r.HandleFunc("/cards/search", SearchCardsHandler).Methods("GET")
	r.HandleFunc("/cards/type/{type}", GetCardsByTypeHandler).Methods("GET")
	r.HandleFunc("/cards/element/{element}", GetCardsByElementHandler).Methods("GET")
	r.HandleFunc("/cards/{id}", GetCardByIDHandler).Methods("GET")

	// Table endpoints (protected, requires authentication)
	protected.HandleFunc("/tables", CreateTable).Methods("POST")
	protected.HandleFunc("/tables", GetUserTables).Methods("GET")
	protected.HandleFunc("/tables/{id}", UpdateTable).Methods("PUT")
	protected.HandleFunc("/tables/{id}/time", UpdateUserTableTime).Methods("PUT")

	return r
}
