package utils

import (
	"context"
	"fmt"
	"bd-backend/internal/models"
)

type contextKey string

const (
	UserIDKey contextKey = "user_id"
	RoleKey   contextKey = "role"
)

// GetRoleFromContext retrieves the user's role from the context
func GetRoleFromContext(ctx context.Context) (models.Role, error) {
	role, ok := ctx.Value(RoleKey).(models.Role)
	if !ok {
		return "", fmt.Errorf("role not found in context or invalid type")
	}
	return role, nil
}

// GetUserIDFromContext retrieves the user's ID from the context
func GetUserIDFromContext(ctx context.Context) (int, error) {
	userID, ok := ctx.Value(UserIDKey).(int)
	if !ok {
		return 0, fmt.Errorf("user ID not found in context or invalid type")
	}
	return userID, nil
}