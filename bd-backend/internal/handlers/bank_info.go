package handlers

import (
	"bd-backend/internal/db/banks"
	"bd-backend/internal/models"
	"bd-backend/internal/utils"
	"context"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
)

// GetBankInfoHandler retrieves bank information based on various criteria from the query parameter.
func GetBankInfoHandler(dbPool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		vars := mux.Vars(r)
		query := vars["query"] // The query parameter can be name, routing number, bank ID, or cert

		// Normalize the query input
		normalizedQuery := strings.TrimSpace(query)
		if len(normalizedQuery) == 0 {
			http.Error(w, "Query parameter cannot be empty", http.StatusBadRequest)
			return
		}

		// Attempt to find the bank using the query
		bank, err := findBankByCriteria(ctx, dbPool, normalizedQuery)
		if err != nil {
			http.Error(w, "Failed to retrieve bank information", http.StatusInternalServerError)
			return
		}

		if bank == nil {
			http.Error(w, "Bank not found", http.StatusNotFound)
			return
		}

		utils.RespondWithJSON(w, http.StatusOK, bank)
	}
}

// findBankByCriteria tries to find a bank based on various criteria.
func findBankByCriteria(ctx context.Context, dbPool *pgxpool.Pool, query string) (*models.UserBankInfo, error) {
	// Attempt to find the bank by different criteria
	// 1. Search by name
	bank, err := banks.FindBank(ctx, dbPool, "name", query)
	if err == nil && bank != nil {
		return bank, nil
	}

	// 2. Search by routing number
	bank, err = banks.FindBank(ctx, dbPool, "routing", query)
	if err == nil && bank != nil {
		return bank, nil
	}

	// 3. Search by ID (try both string and integer representations)
	bank, err = banks.FindBank(ctx, dbPool, "bank_id", query)
	if err == nil && bank != nil {
		return bank, nil
	}

	// 4. Search by certification number
	bank, err = banks.FindBank(ctx, dbPool, "cert", query)
	if err == nil && bank != nil {
		return bank, nil
	}

	return nil, nil
}
