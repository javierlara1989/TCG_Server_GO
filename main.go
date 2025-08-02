package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"tcg-server-go/handlers"
)

func main() {
	router := handlers.SetupRoutes()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Server started on port %s\n", port)
	fmt.Println("Available endpoints:")
	fmt.Println("  POST /login - User authentication")
	fmt.Println("  GET  /health - Server health check")
	fmt.Println("  GET  /api/validate - Token validation (requires authentication)")

	log.Fatal(http.ListenAndServe(":"+port, router))
}
