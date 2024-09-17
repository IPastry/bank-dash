package users

import (
	"bd-backend/internal/utils"
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
)

// UpdateProfile updates the user's profile with the provided details.
// It handles both initial updates (where all fields must be provided) and partial updates (where some fields can be nil).
func UpdateProfile(ctx context.Context, tx pgx.Tx, firstName, lastName, phoneNumber *string) error {
	// Retrieve the user ID from the context
	userID, err := utils.GetUserIDFromContext(ctx)
	if err != nil {
		return fmt.Errorf("error retrieving user ID from context: %w", err)
	}

	// Prepare the SQL update query
	query := `UPDATE users SET`
	var args []interface{}
	var argCount int

	// Dynamically build the query based on which fields are not nil
	if firstName != nil {
		query += ` first_name = $1,`
		args = append(args, *firstName)
		argCount++
	}
	if lastName != nil {
		query += ` last_name = $2,`
		args = append(args, *lastName)
		argCount++
	}
	if phoneNumber != nil {
		if err := utils.ValidatePhoneNumber(*phoneNumber); err != nil {
			return fmt.Errorf("error validating phone number: %w", err)
		}
		query += ` phone_number = $3,`
		args = append(args, *phoneNumber)
		argCount++
	}

	// Remove the trailing comma if any fields are updated
	if argCount > 0 {
		query = query[:len(query)-1] // Remove last comma
		query += ` WHERE user_id = $4`
		args = append(args, userID)
	} else {
		return fmt.Errorf("no fields to update")
	}

	// Perform the database operation with the provided context and transaction
	_, err = tx.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("error updating user profile: %w", err)
	}

	return nil
}

// UpdateUserBankInfo updates the bank ID for a user in the database.
func UpdateUserBankID(ctx context.Context, tx pgx.Tx, bankID int) error {
	// Retrieve the user ID from context
	userID, err := utils.GetUserIDFromContext(ctx)
	if err != nil {
		return fmt.Errorf("error retrieving user ID from context: %w", err)
	}

	// Update the user's bank ID
	_, err = tx.Exec(ctx, `UPDATE users SET bank_id = $1 WHERE user_id = $2`, bankID, userID)
	if err != nil {
		return fmt.Errorf("error updating bank ID for user %d: %w", userID, err)
	}

	return nil
}

// UpdateEmail updates the user's email address and marks it as unverified.
// It returns an error if the update fails or if the email format is invalid.
func UpdateEmail(ctx context.Context, tx pgx.Tx, newEmail string) error {
	// Retrieve the user ID from the context
	userID, err := utils.GetUserIDFromContext(ctx)
	if err != nil {
		return fmt.Errorf("error retrieving user ID from context: %w", err)
	}

	// Validate the new email format
	if err := utils.ValidateEmail(newEmail); err != nil {
		return fmt.Errorf("invalid new email address: %w", err)
	}

	// Execute the update query with the provided context and transaction
	result, err := tx.Exec(ctx, `
		UPDATE users
		SET email = $1, is_verified = false
		WHERE user_id = $2`,
		newEmail, userID)

	if err != nil {
		return fmt.Errorf("error executing update query: %w", err)
	}

	// Check if any rows were affected by the update
	if result.RowsAffected() == 0 {
		return fmt.Errorf("no user found with ID %d", userID)
	}

	return nil
}

// UpdatePassword updates the user's password.
// It returns an error if hashing the password fails or if the update query fails.
func UpdatePassword(ctx context.Context, tx pgx.Tx, newPassword string) error {
	// Retrieve the user ID from the context
	userID, err := utils.GetUserIDFromContext(ctx)
	if err != nil {
		return fmt.Errorf("error retrieving user ID from context: %w", err)
	}

	// Validate the new password
	if err := utils.ValidatePassword(newPassword); err != nil {
		return fmt.Errorf("error validating password: %w", err)
	}

	// Hash the new password
	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("error hashing password: %w", err)
	}

	// Execute the update query with the provided context and transaction
	result, err := tx.Exec(ctx, `
		UPDATE users
		SET password_hash = $1
		WHERE user_id = $2`,
		hashedPassword, userID)

	if err != nil {
		return fmt.Errorf("error executing update query: %w", err)
	}

	// Check if any rows were affected by the update
	if result.RowsAffected() == 0 {
		return fmt.Errorf("no user found with ID %d", userID)
	}

	return nil
}

func logFailedLoginAttempt(ctx context.Context, tx pgx.Tx, userID int, success bool) {
    _, err := tx.Exec(ctx, `INSERT INTO login_attempts (user_id, success) VALUES ($1, $2)`, userID, success)
    if err != nil {
        log.Printf("Error logging failed login attempt for user %d: %v", userID, err)
    }
}

func logSuccessfulLoginAttempt(ctx context.Context, tx pgx.Tx, userID int) {
    logFailedLoginAttempt(ctx, tx, userID, true) // Reuse the same function
}

// LogPasswordReset logs the password reset event for a user.
func LogPasswordReset(ctx context.Context, tx pgx.Tx, userID int) error {
    _, err := tx.Exec(ctx, `INSERT INTO password_resets (user_id, reset_time) VALUES ($1, NOW())`, userID)
    if err != nil {
        log.Printf("Error logging password reset for user %d: %v", userID, err)
        return fmt.Errorf("error logging password reset: %w", err)
    }
    return nil
}