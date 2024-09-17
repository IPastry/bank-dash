package mail

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// SaveTokenInDB stores the token and its expiration time in the database.
func SaveTokenInDB(ctx context.Context, tx pgx.Tx, email, token string, expirationTime time.Time, lastSentTime time.Time) error {
	_, err := tx.Exec(ctx, `
        INSERT INTO email_verifications (email, token, expiration_time, last_verification_email_sent)
        VALUES ($1, $2, $3, $4)
        ON CONFLICT (email)
        DO UPDATE SET token = EXCLUDED.token, expiration_time = EXCLUDED.expiration_time, last_verification_email_sent = EXCLUDED.last_verification_email_sent
    `, email, token, expirationTime, lastSentTime)
	return err
}

// VerifyToken checks if the token is valid and has not expired, and retrieves the associated email.
func VerifyToken(ctx context.Context, dbPool *pgxpool.Pool, token string) (string, bool, error) {
	var email string
	var expirationTime time.Time

	err := dbPool.QueryRow(ctx, `
        SELECT email, expiration_time
        FROM email_verifications
        WHERE token = $1`, token).Scan(&email, &expirationTime)

	if err != nil {
		if err == pgx.ErrNoRows {
			return "", false, nil // Token not found
		}
		return "", false, err
	}

	if time.Now().After(expirationTime) {
		return "", false, nil // Token has expired
	}

	return email, true, nil // Token is valid
}

// CleanupExpiredTokens removes tokens that have expired.
func CleanupExpiredTokens(ctx context.Context, dbPool *pgxpool.Pool) error {
	_, err := dbPool.Exec(ctx, `
        DELETE FROM email_verifications
        WHERE expiration_time < $1`, time.Now())
	return err
}

// CheckLastSentTime checks if the last email was sent recently and enforces a cooldown period.
func CheckLastSentTime(ctx context.Context, dbPool *pgxpool.Pool, email string) (bool, error) {
	var lastSentTime time.Time
	err := dbPool.QueryRow(ctx, `SELECT last_verification_email_sent FROM email_verifications WHERE email = $1`, email).Scan(&lastSentTime)
	if err != nil {
		return false, err
	}

	// Check if the last sent time is within the cooldown period (e.g., 5 minutes)
	if time.Since(lastSentTime) < time.Minute*5 {
		return false, nil
	}

	return true, nil
}
