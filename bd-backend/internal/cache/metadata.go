package cache

import (
	"bd-backend/internal/utils"
	"context"
	"fmt"
	"strconv"

	"github.com/go-redis/redis/v8"
)

func SetMetadata(ctx context.Context, bankID int, peerGroup string) error {
	client := GetClient()

	// Retrieve userID from context
	userID, err := utils.GetUserIDFromContext(ctx)
	if err != nil {
		return fmt.Errorf("error retrieving user ID from context: %w", err)
	}

	// Convert userID to string for Redis key
	userIDStr := strconv.Itoa(userID)

	// Store bank ID and peer group in Redis
	return client.HSet(ctx, "user:"+userIDStr+":meta", "bank_id", bankID, "peergroup", peerGroup).Err()
}

func GetMetadata(ctx context.Context, userID int) (int, string, error) {
	client := GetClient()
	userIDStr := strconv.Itoa(userID)
	result := client.HGetAll(ctx, "user:"+userIDStr+":meta")

	// Wait for all results
	metadata, err := result.Result()
	if err != nil {
		return 0, "", err
	}

	// Extract bank_id and peergroup from the metadata
	bankIDStr, bankIDOk := metadata["bank_id"]
	peergroup, peergroupOk := metadata["peergroup"]

	if !bankIDOk || !peergroupOk {
		return 0, "", redis.Nil
	}

	bankID, err := strconv.Atoi(bankIDStr)
	if err != nil {
		return 0, "", err
	}

	return bankID, peergroup, nil
}

// GetBankID retrieves the bank ID from metadata.
func GetBankID(ctx context.Context, userID int) (int, error) {
	bankID, _, err := GetMetadata(ctx, userID)
	if err != nil {
		return 0, err
	}
	return bankID, nil
}

// GetPeerGroup retrieves the peer group from metadata.
func GetPeerGroup(ctx context.Context, userID int) (string, error) {
	_, peergroup, err := GetMetadata(ctx, userID)
	if err != nil {
		return "", err
	}
	return peergroup, nil
}
