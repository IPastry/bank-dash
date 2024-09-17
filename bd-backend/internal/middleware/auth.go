package middleware

import (
	"bd-backend/internal/auth"
	"bd-backend/internal/utils"
	"bd-backend/internal/models"
	"context"
	"net/http"
	"strings"
)

// AuthMiddleware is a middleware handler that validates JWT tokens and sets user context
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Extract the token from the "Bearer " prefix
		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		// Parse and validate the token
		token, claims, err := auth.ParseToken(tokenStr)
		if err != nil || token == nil || !token.Valid {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Assert user_id to integer
		userIDFloat, ok := claims["user_id"].(float64)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		userID := int(userIDFloat)

		// Assert role to string and validate
		roleStr, ok := claims["role"].(string)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		role := models.Role(roleStr)

		// Validate role using the helper function
		if !utils.IsValidRole(role) {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Set values in context
		ctx := context.WithValue(r.Context(), utils.UserIDKey, userID)
		ctx = context.WithValue(ctx, utils.RoleKey, role)
		r = r.WithContext(ctx)

		// Proceed to the next handler
		next.ServeHTTP(w, r)
	})
}
