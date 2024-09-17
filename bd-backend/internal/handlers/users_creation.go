package handlers

import (
	"bd-backend/internal/db/users"
	"bd-backend/internal/mail"
	"bd-backend/internal/utils"
	"encoding/json"
	"log"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

// CreateElevatedUserHandler handles the creation of an elevated user
func CreateElevatedUserHandler(dbPool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Define a structure for request data
		var requestData struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		// Decode JSON body
		if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
			log.Printf("Error decoding request body: %v", err)
			http.Error(w, `{"error": "Invalid request body"}`, http.StatusBadRequest)
			return
		}

		// Validate email
		if err := utils.ValidateEmail(requestData.Email); err != nil {
			log.Printf("Invalid email address %s: %v", requestData.Email, err)
			http.Error(w, `{"error": "Invalid email address"}`, http.StatusBadRequest)
			return
		}

		// Validate password
		if err := utils.ValidatePassword(requestData.Password); err != nil {
			log.Printf("Invalid password for email %s: %v", requestData.Email, err)
			http.Error(w, `{"error": "Invalid password"}`, http.StatusBadRequest)
			return
		}

		// Start a new transaction
		tx, err := dbPool.Begin(ctx)
		if err != nil {
			log.Printf("Error starting transaction: %v", err)
			http.Error(w, `{"error": "Error starting transaction"}`, http.StatusInternalServerError)
			return
		}
		defer tx.Rollback(ctx) // Ensure rollback in case of error

		// Create elevated user
		if err := users.CreateElevatedUser(ctx, tx, requestData.Email, requestData.Password); err != nil {
			log.Printf("Error creating elevated user with email %s: %v", requestData.Email, err)
			http.Error(w, `{"error": "Error creating user"}`, http.StatusInternalServerError)
			return
		}

		// Commit the transaction
		if err := tx.Commit(ctx); err != nil {
			log.Printf("Error committing transaction: %v", err)
			http.Error(w, `{"error": "Error committing transaction"}`, http.StatusInternalServerError)
			return
		}

		// Initiate email verification
		if err := mail.InitiateEmailVerification(ctx, dbPool, requestData.Email); err != nil {
			log.Printf("Error sending verification email to %s: %v", requestData.Email, err)
			http.Error(w, `{"error": "Error sending verification email"}`, http.StatusInternalServerError)
			return
		}

		// Send success response
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"message": "User created successfully"}`))
	}
}

// CreateSharedAccountHandler handles the creation of a shared account
func CreateSharedAccountHandler(dbPool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Define a structure for request data
		var requestData struct {
			Email       string  `json:"email"`
			Password    string  `json:"password"`
			PhoneNumber *string `json:"phone_number"`
		}

		// Decode JSON body
		if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
			http.Error(w, `{"error": "Invalid request body"}`, http.StatusBadRequest)
			return
		}

		// Validate email
		if err := utils.ValidateEmail(requestData.Email); err != nil {
			http.Error(w, `{"error": "Invalid email address"}`, http.StatusBadRequest)
			return
		}

		// Validate password
		if err := utils.ValidatePassword(requestData.Password); err != nil {
			http.Error(w, `{"error": "Invalid password"}`, http.StatusBadRequest)
			return
		}

		// Start transaction
		tx, err := dbPool.Begin(ctx)
		if err != nil {
			http.Error(w, `{"error": "Error starting transaction"}`, http.StatusInternalServerError)
			return
		}
		defer tx.Rollback(ctx) // Rollback if something goes wrong

		// Create shared account
		if err := users.CreateSharedAccount(ctx, tx, requestData.Email, requestData.Password, requestData.PhoneNumber); err != nil {
			http.Error(w, `{"error": "Error creating shared account"}`, http.StatusInternalServerError)
			return
		}

		// Commit transaction
		if err := tx.Commit(ctx); err != nil {
			http.Error(w, `{"error": "Error committing transaction"}`, http.StatusInternalServerError)
			return
		}

		// Send success response
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"message": "Shared account created successfully"}`))
	}
}