package users

import (
	// "bd-backend/internal/cache"
	"bd-backend/internal/models"
	"bd-backend/internal/utils"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

// CreateElevatedUser inserts a new elevated user into the database using a transaction
func CreateElevatedUser(ctx context.Context, tx pgx.Tx, email, password string) error {
	if err := utils.ValidateEmail(email); err != nil {
		return fmt.Errorf("invalid email address: %w", err)
	}

	if err := utils.ValidatePassword(password); err != nil {
		return fmt.Errorf("invalid password: %w", err)
	}

	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return fmt.Errorf("error hashing password: %w", err)
	}

	// Use models.Elevated for the role
	_, err = tx.Exec(ctx, `INSERT INTO users (email, password_hash, role) VALUES ($1, $2, $3)`,
		email, hashedPassword, models.Elevated) // Ensure models.Elevated is the correct value for your enum
	if err != nil {
		return fmt.Errorf("error inserting elevated user into database: %w", err)
	}

	return nil
}

// CreateSharedAccount creates a new shared account for the bank and sets its managerID
func CreateSharedAccount(ctx context.Context, tx pgx.Tx, email, password string, phoneNumber *string) error {
	// Retrieve user ID from context
	userID, err := utils.GetUserIDFromContext(ctx)
	if err != nil {
		return fmt.Errorf("error retrieving user ID from context: %w", err)
	}

	// // Retrieve bank ID from metadata using the optimized function
	// bankID, err := cache.GetBankID(ctx, userID)
	// if err != nil {
	// 	return fmt.Errorf("error retrieving bank ID from metadata: %w", err)
	// }

	// Validate the email format
	if err := utils.ValidateEmail(email); err != nil {
		return fmt.Errorf("invalid email address: %w", err)
	}

	// Validate the password
	if err := utils.ValidatePassword(password); err != nil {
		return fmt.Errorf("invalid password: %w", err)
	}

	// Check if phoneNumber is nil and return an error if it is
	if phoneNumber == nil {
		return fmt.Errorf("phone number must be provided")
	}

	// Hash the password
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return fmt.Errorf("error hashing password: %w", err)
	}


	//TEST VERSION DELETE AFTER
	_, err = tx.Exec(ctx, `
        INSERT INTO users (email, password_hash, role, manager_id, phone_number, is_verified)
        VALUES ($1, $2, $3, $4, $5, true)`,
		email, hashedPassword, models.Regular, userID, *phoneNumber)
	if err != nil {
		return fmt.Errorf("error inserting shared account into database: %w", err)
	}
	// // Perform the database operation within the transaction
	// _, err = tx.Exec(ctx, `
  //       INSERT INTO users (email, password_hash, role, manager_id, bank_id, phone_number, is_verified)
  //       VALUES ($1, $2, $3, $4, $5, $6, true)`,
	// 	email, hashedPassword, models.Regular, userID, bankID, *phoneNumber)
	// if err != nil {
	// 	return fmt.Errorf("error inserting shared account into database: %w", err)
	// }

	return nil
}