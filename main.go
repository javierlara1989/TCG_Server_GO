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

	log.Fatal(http.ListenAndServe(":"+port, router))
}
