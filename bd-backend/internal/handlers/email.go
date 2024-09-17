package handlers

import (
	"bd-backend/internal/mail"
	"bd-backend/internal/utils"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// EmailResendHandler handles the request to resend the verification email.
func EmailResendHandler(dbPool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		var requestData struct {
			Email string `json:"email"`
		}

		// Decode the request body
		if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
			log.Printf("Error decoding request body for email %s: %v", requestData.Email, err)
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Start a database transaction
		tx, err := dbPool.Begin(ctx)
		if err != nil {
			log.Printf("Error starting transaction for email %s: %v", requestData.Email, err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		defer tx.Rollback(ctx) // Rollback on error

		// Check if email can be resent
		isAllowed, err := mail.CheckLastSentTime(ctx, dbPool, requestData.Email)
		if err != nil {
			log.Printf("Error checking last sent time for email %s: %v", requestData.Email, err)
			http.Error(w, "Error checking last sent time", http.StatusInternalServerError)
			return
		}
		if !isAllowed {
			http.Error(w, "Verification email already sent recently", http.StatusTooManyRequests)
			return
		}

		// Generate a new verification token
		token, err := utils.GenerateEmailVerificationToken()
		if err != nil {
			log.Printf("Error generating verification token for email %s: %v", requestData.Email, err)
			http.Error(w, "Error generating verification token", http.StatusInternalServerError)
			return
		}
		expirationTime := time.Now().Add(24 * time.Hour)

		// Save the token in the database
		err = mail.SaveTokenInDB(ctx, tx, requestData.Email, token, expirationTime, time.Now())
		if err != nil {
			log.Printf("Error saving new token in database for email %s: %v", requestData.Email, err)
			http.Error(w, "Error saving new token in database", http.StatusInternalServerError)
			return
		}

		// Send the verification email
		err = mail.SendVerificationEmail(ctx, requestData.Email, token)
		if err != nil {
			log.Printf("Error sending verification email for email %s: %v", requestData.Email, err)
			http.Error(w, "Error sending verification email", http.StatusInternalServerError)
			return
		}

		// Commit the transaction
		if err := tx.Commit(ctx); err != nil {
			log.Printf("Error committing transaction for email %s: %v", requestData.Email, err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		// Respond with success
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Verification email resent successfully"))
	}
}
