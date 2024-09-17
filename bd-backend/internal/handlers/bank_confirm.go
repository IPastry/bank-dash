package handlers

import (
	"bd-backend/internal/cache"
	"bd-backend/internal/db/banks"
	"bd-backend/internal/db/users"
	"bd-backend/internal/utils"
	"encoding/json"
	"log"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

// ConfirmBankHandler handles the bank confirmation and updates the user information.
func ConfirmBankHandler(dbPool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Decode request body
		var requestBody struct {
			BankID int `json:"bank_id"`
		}
		if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			log.Printf("Error decoding request body: %v", err)
			return
		}

		// Start a transaction for updating user info
		tx, err := dbPool.Begin(ctx)
		if err != nil {
			http.Error(w, "Error starting transaction", http.StatusInternalServerError)
			log.Printf("Error starting transaction: %v", err)
			return
		}
		defer func() {
			if err != nil {
				if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
					log.Printf("Error rolling back transaction: %v", rollbackErr)
				}
			}
		}()

		// Update user with selected bank information
		err = users.UpdateUserBankID(ctx, tx, requestBody.BankID)
		if err != nil {
			http.Error(w, "Error updating user bank information", http.StatusInternalServerError)
			log.Printf("Error updating user bank information: %v", err)
			return
		}

		// Commit the transaction
		if err = tx.Commit(ctx); err != nil {
			http.Error(w, "Error committing transaction", http.StatusInternalServerError)
			log.Printf("Error committing transaction: %v", err)
			return
		}

		// Retrieve peer group for the bank ID
		peerGroup, err := banks.GetPeerGroupFromBankID(ctx, dbPool, requestBody.BankID)
		if err != nil {
			http.Error(w, "Error retrieving peer group", http.StatusInternalServerError)
			log.Printf("Error retrieving peer group for bank ID %d: %v", requestBody.BankID, err)
			return
		}

		// Store metadata in Redis
		err = cache.SetMetadata(ctx, requestBody.BankID, peerGroup)
		if err != nil {
			http.Error(w, "Error storing metadata", http.StatusInternalServerError)
			log.Printf("Error storing metadata in cache: %v", err)
			return
		}

		// store report data in Redis 
		err = cache.SetReportData(ctx, dbPool, requestBody.BankID)
		if err != nil {
			http.Error(w, "Error caching report data", http.StatusInternalServerError)
			return
		}

		// Send success response
		utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Bank information updated successfully"})
	}
}
