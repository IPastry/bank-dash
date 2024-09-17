package users

import (
	"bd-backend/internal/models"
	"bd-backend/internal/utils"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// VerifyUserPasswordAndRole verifies the user's password and retrieves their role.
func VerifyUserPasswordAndRole(ctx context.Context, tx pgx.Tx, email, password string) (bool, int, models.Role, error) {
    var hashedPassword string
    var userID int
    var role models.Role

    // Retrieve user ID, hashed password, and role
    err := tx.QueryRow(ctx, `SELECT user_id, password_hash, role FROM users WHERE email = $1`, email).Scan(&userID, &hashedPassword, &role)
    if err != nil {
        if err == pgx.ErrNoRows {
            log.Printf("User with email %s not found", email)
            logFailedLoginAttempt(ctx, tx, 0, false) // Log attempt with userID=0 for unknown user
            return false, 0, "", nil // User not found
        }
        log.Printf("Error retrieving user with email %s: %v", email, err)
        return false, 0, "", err
    }

    // Compare password
    err = utils.CompareHashAndPassword(hashedPassword, password)
    if err != nil {
        log.Printf("Invalid password for user %d: %v", userID, err)
        logFailedLoginAttempt(ctx, tx, userID, false) // Log failed attempt
        return false, 0, "", nil // Incorrect password
    }

    logSuccessfulLoginAttempt(ctx, tx, userID) // Log successful attempt
    log.Printf("Successful login for user %d", userID)
    return true, userID, role, nil // Password correct
}

// RetrieveUserRole retrieves the user's role based on their user ID.
func RetrieveUserRole(ctx context.Context, dbPool *pgxpool.Pool, userID int) (models.Role, error) {
    var role models.Role
    err := dbPool.QueryRow(ctx, "SELECT role FROM users WHERE user_id = $1", userID).Scan(&role)
    if err != nil {
        return "", err
    }
    return role, nil
}

// IsProfileComplete checks if the user's profile information is complete
func IsProfileComplete(ctx context.Context, tx pgx.Tx, userID int) (bool, error) {
	var firstName, lastName, phoneNumber sql.NullString
	err := tx.QueryRow(ctx, `SELECT first_name, last_name, phone_number FROM users WHERE user_id = $1`, userID).Scan(&firstName, &lastName, &phoneNumber)
	if err != nil {
		log.Printf("Error checking profile completeness for user %d: %v", userID, err)
		return false, err
	}

	isComplete := firstName.Valid && lastName.Valid && phoneNumber.Valid
	if isComplete {
		log.Printf("User %d profile is complete", userID)
	} else {
		log.Printf("User %d profile is incomplete", userID)
	}

	return isComplete, nil
}

// GetUserByEmail retrieves a user by email, used to verify the email address.
func GetUserByEmail(ctx context.Context, dbPool *pgxpool.Pool, email string) (map[string]interface{}, error) {
	var emailStr string
	var isVerified bool
	row := dbPool.QueryRow(ctx, `SELECT email, is_verified FROM users WHERE email = $1`, email)
	err := row.Scan(&emailStr, &isVerified)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("error querying user by email: %w", err)
	}
	user := map[string]interface{}{
		"email":       emailStr,
		"is_verified": isVerified,
	}
	return user, nil
}

// GetUserInfo retrieves the full profile of a user by email.
func GetUserInfo(ctx context.Context, dbPool *pgxpool.Pool, email string) (map[string]interface{}, error) {
	var emailStr string
	var firstName, lastName, phoneNumber *string
	var bankID *int
	row := dbPool.QueryRow(ctx, `SELECT email, first_name, last_name, phone_number, bank_id FROM users WHERE email = $1`, email)
	err := row.Scan(&emailStr, &firstName, &lastName, &phoneNumber, &bankID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("error querying user info: %w", err)
	}
	userInfo := map[string]interface{}{
		"email":        emailStr,
		"first_name":   defaultString(firstName),
		"last_name":    defaultString(lastName),
		"phone_number": defaultString(phoneNumber),
		"bank_id":      defaultInt(bankID),
	}
	return userInfo, nil
}

// VerifyUser updates the is_verified field in the database to true.
func VerifyUser(ctx context.Context, dbPool *pgxpool.Pool, email string) error {
	result, err := dbPool.Exec(ctx, `UPDATE users SET is_verified = true WHERE email = $1`, email)
	if err != nil {
		return fmt.Errorf("error updating user verification status: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("no user found with email %s", email)
	}
	return nil
}

// validateResetToken validates the password reset token and retrieves the user ID
func ValidateResetToken(ctx context.Context, dbPool *pgxpool.Pool, token string) (int, error) {
    var userID int

    // Example query to validate token and get user ID
    err := dbPool.QueryRow(ctx, `SELECT user_id FROM password_resets WHERE token = $1 AND reset_time > NOW() - INTERVAL '1 hour'`, token).Scan(&userID)
    if err != nil {
        if err == pgx.ErrNoRows {
            return 0, fmt.Errorf("invalid or expired token")
        }
        return 0, fmt.Errorf("error validating token: %w", err)
    }

    return userID, nil
}