package mail

import (
	"bd-backend/internal/utils"
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// InitiateEmailVerification sends the verification email and stores the token
func InitiateEmailVerification(ctx context.Context, dbPool *pgxpool.Pool, email string) error {
	// Start a transaction
	tx, err := dbPool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			// Rollback the transaction in case of an error
			tx.Rollback(ctx)
		} else {
			// Commit the transaction if everything went well
			tx.Commit(ctx)
		}
	}()

	// Generate the email verification token
	token, err := utils.GenerateEmailVerificationToken()
	if err != nil {
		return err
	}

	// Set token expiration time (24 hours)
	expirationTime := time.Now().Add(24 * time.Hour)

	// Save the token and expiration in the database
	err = SaveTokenInDB(ctx, tx, email, token, expirationTime, time.Now()) // Pass current time for last sent time
	if err != nil {
		return err
	}

	// Send the verification email with the generated token
	err = SendVerificationEmail(ctx, email, token) // Use context when sending the email
	if err != nil {
		return err
	}

	return nil
}
