package utils

import (
	"bd-backend/internal/models"
	"fmt"
	"regexp"
)

// ValidateEmail checks if an email is valid
func ValidateEmail(email string) error {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !re.MatchString(email) {
		return fmt.Errorf("invalid email address")
	}
	return nil
}

// ValidatePhoneNumber checks if a phone number is valid based on a regex pattern
func ValidatePhoneNumber(phoneNumber string) error {
	// Basic phone number regex (adjust pattern as needed)
	phoneRegex := `^\+?[1-9]\d{1,14}$`
	matched, err := regexp.MatchString(phoneRegex, phoneNumber)
	if err != nil {
		return fmt.Errorf("error validating phone number: %w", err)
	}
	if !matched {
		return fmt.Errorf("invalid phone number format")
	}
	return nil
}

// ValidatePassword checks if the password meets security requirements
func ValidatePassword(password string) error {
	// Define password requirements
	const (
		minLength = 8
		maxLength = 20
	)

	// Check password length
	if len(password) < minLength || len(password) > maxLength {
		return fmt.Errorf("password must be between %d and %d characters", minLength, maxLength)
	}

	// Check for at least one uppercase letter
	if !regexp.MustCompile(`[A-Z]`).MatchString(password) {
		return fmt.Errorf("password must contain at least one uppercase letter")
	}

	// Check for at least one lowercase letter
	if !regexp.MustCompile(`[a-z]`).MatchString(password) {
		return fmt.Errorf("password must contain at least one lowercase letter")
	}

	// Check for at least one digit
	if !regexp.MustCompile(`[0-9]`).MatchString(password) {
		return fmt.Errorf("password must contain at least one digit")
	}

	// Check for at least one special character
	if !regexp.MustCompile(`[!@#$%^&*()_+{}\[\]:;"'<>,.?/~\\|-]`).MatchString(password) {
		return fmt.Errorf("password must contain at least one special character")
	}

	return nil
}

// Helper function to validate if the role is valid
func IsValidRole(role models.Role) bool {
	validRoles := []models.Role{
		models.Admin,
		models.Elevated,
		models.Regular,
	}

	for _, validRole := range validRoles {
		if role == validRole {
			return true
		}
	}
	return false
}
