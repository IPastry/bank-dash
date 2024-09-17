package main

import (
	"bd-backend/internal/db"
	// "bd-backend/internal/redis"
	"bd-backend/internal/routes"
	"context"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// // Initialize Redis
	// err = redis.InitRedis()
	// if err != nil {
	//     log.Fatalf("Error initializing Redis: %v", err)
	// }

	// // Test Redis connection
	// err = redis.TestConnection()
	// if err != nil {
	//     log.Fatalf("Redis connection test failed: %v", err)
	// }

	// Create a context
	ctx := context.Background()

	// Initialize the database connection
	dbPool, err := db.ConnectDB(ctx)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	defer db.CloseDB(dbPool)

	// Initialize the router with middlewares
	router := routes.SetupRoutes(dbPool)

	// Start the HTTP server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port if not set in .env
	}
	log.Printf("Server is running on port %s", port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatalf("Could not start server: %v", err)
	}
}
