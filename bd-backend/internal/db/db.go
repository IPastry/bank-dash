package db

import (
    "context"
    "fmt"
    "os"
    "github.com/jackc/pgx/v5/pgxpool"
)

// ConnectDB initializes and returns a new database connection pool.
func ConnectDB(ctx context.Context) (*pgxpool.Pool, error) {
    // Load database URL from environment variable
    databaseURL := os.Getenv("DATABASE_URL")
    if databaseURL == "" {
        return nil, fmt.Errorf("DATABASE_URL environment variable is not set")
    }

    // Create a connection pool
    dbPool, err := pgxpool.New(ctx, databaseURL)
    if err != nil {
        return nil, fmt.Errorf("unable to connect to database: %v", err)
    }

    // Test the connection
    err = dbPool.Ping(ctx)
    if err != nil {
        dbPool.Close()
        return nil, fmt.Errorf("unable to ping database: %v", err)
    }

    fmt.Println("Successfully connected to the database")
    return dbPool, nil
}

// CloseDB closes the database connection pool.
func CloseDB(dbPool *pgxpool.Pool) {
    if dbPool != nil {
        dbPool.Close()
    }
}
