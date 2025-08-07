package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"tcg-server-go/database"
	"tcg-server-go/handlers"
)

func main() {
	// Initialize database connection
	if err := database.Connect(); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer database.Close()

	// Create database tables
	if err := database.CreateTables(); err != nil {
		log.Fatal("Failed to create database tables:", err)
	}

	router := handlers.SetupRoutes()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Server started on port %s\n", port)
	fmt.Println("Available endpoints:")
	fmt.Println("  POST /register - User registration")
	fmt.Println("  POST /verify-email - Email verification")
	fmt.Println("  POST /resend-code - Resend validation code")
	fmt.Println("  POST /login - User authentication")
	fmt.Println("  GET  /health - Server health check")
	fmt.Println("  GET  /api/validate - Token validation (requires authentication)")
	fmt.Println("  GET  /api/user-info - Get user game info (requires authentication)")
	fmt.Println("  GET  /api/user-cards - Get user's card inventory (requires authentication)")
	fmt.Println("  GET  /api/user-cards/{id} - Get specific user card (requires authentication)")
	fmt.Println("  GET  /api/decks - Get user's decks (requires authentication)")
	fmt.Println("  POST /api/decks - Create new deck (requires authentication)")
	fmt.Println("  GET  /api/decks/limit - Get deck limit information (requires authentication)")
	fmt.Println("  GET  /api/decks/{id} - Get specific deck (requires authentication)")
	fmt.Println("  GET  /api/decks/{id}/cards - Get deck with cards (requires authentication)")
	fmt.Println("  DELETE /api/decks/{id} - Delete deck (requires authentication)")
	fmt.Println("")
	fmt.Println("Card Management (Read-only):")
	fmt.Println("  GET  /cards - Get all cards")
	fmt.Println("  GET  /cards/search?q=<term> - Search cards by name")
	fmt.Println("  GET  /cards/type/{type} - Get cards by type (Monster/Spell/Energy)")
	fmt.Println("  GET  /cards/element/{element} - Get cards by element")
	fmt.Println("  GET  /cards/{id} - Get card by ID")

	log.Fatal(http.ListenAndServe(":"+port, router))
}
