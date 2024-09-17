package auth

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"bd-backend/internal/models"
	"bd-backend/internal/utils"

	"github.com/golang-jwt/jwt/v4"
	"github.com/joho/godotenv"
)

var (
	jwtSecret []byte
	// Configurable expiration times
	accessTokenExpiry  time.Duration
	refreshTokenExpiry time.Duration
)

func Init() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		panic("Error loading .env file")
	}

	jwtSecret = []byte(os.Getenv("JWT_SECRET_KEY"))

	// Load expiration times from environment variables or default values
	accessTokenExpiry = time.Hour * time.Duration(getEnvInt("ACCESS_TOKEN_EXPIRY_HOURS", 72))
	refreshTokenExpiry = time.Hour * 24 * time.Duration(getEnvInt("REFRESH_TOKEN_EXPIRY_DAYS", 30))
}

func getEnvInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}

// GenerateToken generates a JWT token for a given user ID and role
func GenerateToken(userID int, role models.Role) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"exp":     time.Now().Add(accessTokenExpiry).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return signedToken, nil
}

// ParseToken parses and validates the JWT token
func ParseToken(tokenStr string) (*jwt.Token, map[string]interface{}, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse token: %w", err)
	}

	claims, ok := token.Claims.(*jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, nil, fmt.Errorf("invalid token")
	}

	roleStr, ok := (*claims)["role"].(string)
	if !ok {
		return nil, nil, fmt.Errorf("role claim missing or invalid")
	}

	role := models.Role(roleStr)
	if !utils.IsValidRole(role) {
		return nil, nil, fmt.Errorf("invalid role in token: %s", role)
	}

	return token, *claims, nil
}

// GenerateRefreshToken generates a refresh token
func GenerateRefreshToken(userID int) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(refreshTokenExpiry).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", fmt.Errorf("failed to sign refresh token: %w", err)
	}

	return signedToken, nil
}

// ParseRefreshToken parses and validates the refresh token
func ParseRefreshToken(tokenStr string) (*jwt.Token, map[string]interface{}, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse refresh token: %w", err)
	}

	claims, ok := token.Claims.(*jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, nil, fmt.Errorf("invalid refresh token")
	}

	return token, *claims, nil
}
