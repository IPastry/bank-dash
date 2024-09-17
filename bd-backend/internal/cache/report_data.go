package cache

import (
	"bd-backend/internal/db/data"
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v5/pgxpool"
)

// SetReportData retrieves recent report data from the database and stores it in Redis.
func SetReportData(ctx context.Context, db *pgxpool.Pool, userID int) error {
    client := GetClient() // Get Redis client

    // Fetch the bank ID associated with the user from the cache
    bankID, err := GetBankID(ctx, userID)
    if err != nil {
        return fmt.Errorf("error retrieving bank ID from cache: %w", err)
    }

    // Retrieve recent report data for the bank ID
    reportData, err := data.FetchRecentReportData(ctx, db, bankID)
    if err != nil {
        return fmt.Errorf("error fetching recent report data: %w", err)
    }

    // Create a Redis key for the user's report data
    cacheKey := fmt.Sprintf("user:%d:report_data", userID)

    // Start a Redis pipeline transaction
    pipe := client.Pipeline()
    for _, data := range reportData {
        // Use section, date, name, and metric as keys in the nested hash
        sectionKey := string(data.Section) // Directly use sectionKey as string

        // Convert CustomDate to string
        dateStr, err := data.Date.Value()
        if err != nil {
            return fmt.Errorf("error formatting date %v: %w", data.Date, err)
        }

        // Ensure dateStr is a string
        dateValue, ok := dateStr.(string)
        if !ok {
            return fmt.Errorf("error converting date value to string: %v", dateStr)
        }

        nameKey := data.Name

        // Ensure field is a string
        field := fmt.Sprintf("%s:%s", nameKey, data.Metric)

        // Handle nullable value
        var valueStr string
        if data.Value != nil {
            valueStr = *data.Value
        } else {
            valueStr = "N/A" // or provide a different representation for nil
        }

        // Store each piece of report data in Redis
        pipe.HSet(ctx, fmt.Sprintf("%s:%s", cacheKey, sectionKey), dateValue, fmt.Sprintf("%s:%s", field, valueStr))
    }

    // Execute the pipeline
    _, err = pipe.Exec(ctx)
    if err != nil {
        return fmt.Errorf("error executing Redis pipeline: %w", err)
    }

    return nil
}

// GetCachedReportData retrieves cached report data for a specific user and date.
func GetCachedReportData(ctx context.Context, userID int, date string, metric string) (string, error) {
	client := GetClient()
	cacheKey := fmt.Sprintf("user:%d:report_data", userID)

	field := fmt.Sprintf("%s:%s", date, metric)
	result, err := client.HGet(ctx, cacheKey, field).Result()
	if err != nil {
		if err == redis.Nil {
			return "", nil // Data not found in cache
		}
		return "", fmt.Errorf("error retrieving cached data: %w", err)
	}

	return result, nil
}
