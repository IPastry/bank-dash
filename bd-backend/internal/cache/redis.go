package cache

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"

	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
)

var rdb *redis.Client
var ctx = context.Background()

func InitRedis() {
    // Load environment variables from .env file
    if err := godotenv.Load(); err != nil {
        log.Fatalf("Error loading .env file: %v", err)
    }

    // Get Redis configuration from environment variables
    addr := getEnv("REDIS_ADDR", "localhost:6379")
    password := getEnv("REDIS_PASSWORD", "")
    db := getEnvAsInt("REDIS_DB", 0)

    rdb = redis.NewClient(&redis.Options{
        Addr:     addr,
        Password: password,
        DB:       db,
    })

    _, err := rdb.Ping(ctx).Result()
    if err != nil {
        log.Fatalf("Could not connect to Redis: %v", err)
    }

		// Start the graceful shutdown routine
    go setupGracefulShutdown()
}

func GetClient() *redis.Client {
    return rdb
}

func TestConnection() {
    result, err := rdb.Ping(ctx).Result()
    if err != nil {
        log.Fatalf("Redis connection failed: %v", err)
    }
    fmt.Println("Redis connection test result:", result)
}

func setupGracefulShutdown() {
    // Create a channel to listen for interrupt signals
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, os.Interrupt)

    // Block until an interrupt signal is received
    <-quit

    // Initiate graceful shutdown
    fmt.Println("Shutting down Redis connection...")
    if err := rdb.Close(); err != nil {
        log.Fatalf("Error closing Redis connection: %v", err)
    }
    fmt.Println("Redis connection closed.")
}

func getEnv(key, defaultValue string) string {
    if value, exists := os.LookupEnv(key); exists {
        return value
    }
    return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
    valueStr := getEnv(key, fmt.Sprintf("%d", defaultValue))
    value, err := strconv.Atoi(valueStr)
    if err != nil {
        return defaultValue
    }
    return value
}
