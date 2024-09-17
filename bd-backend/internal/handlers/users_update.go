package handlers

import (
	"bd-backend/internal/db/users"
	"encoding/json"
	"net/http"


	"github.com/jackc/pgx/v5/pgxpool"
)

// UpdateProfileHandler handles profile update requests
func UpdateProfileHandler(dbPool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			FirstName   *string `json:"first_name"`
			LastName    *string `json:"last_name"`
			PhoneNumber *string `json:"phone_number"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		// Get context from the request
		ctx := r.Context()

		// Start a transaction
		tx, err := dbPool.Begin(ctx)
		if err != nil {
			http.Error(w, "Error starting transaction", http.StatusInternalServerError)
			return
		}
		defer tx.Rollback(ctx) // Ensure rollback on error

		// Update the user's profile
		err = users.UpdateProfile(ctx, tx, req.FirstName, req.LastName, req.PhoneNumber)
		if err != nil {
			http.Error(w, "Error updating profile", http.StatusInternalServerError)
			return
		}

		// Commit the transaction
		if err := tx.Commit(ctx); err != nil {
			http.Error(w, "Error committing transaction", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"message": "Profile updated successfully"}`))
	}
}

// ResetPasswordHandler handles password reset requests
func ResetPasswordHandler(dbPool *pgxpool.Pool) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var req struct {
            Token       string `json:"token"`
            NewPassword string `json:"new_password"`
        }
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
            http.Error(w, "Invalid request", http.StatusBadRequest)
            return
        }

        // Get context from the request
        ctx := r.Context()

        // Validate reset token and retrieve user ID
        userID, err := users.ValidateResetToken(ctx, dbPool, req.Token)
        if err != nil {
            http.Error(w, "Invalid or expired reset token", http.StatusBadRequest)
            return
        }

        // Begin transaction
        tx, err := dbPool.Begin(ctx)
        if err != nil {
            http.Error(w, "Failed to start transaction", http.StatusInternalServerError)
            return
        }

        defer tx.Rollback(ctx) // Rollback in case of error

        // Update the user's password
        if err := users.UpdatePassword(ctx, tx, req.NewPassword); err != nil {
            http.Error(w, "Failed to update password", http.StatusInternalServerError)
            return
        }

        // Commit the transaction
        if err := tx.Commit(ctx); err != nil {
            http.Error(w, "Failed to commit transaction", http.StatusInternalServerError)
            return
        }

        // Log the password reset
        users.LogPasswordReset(ctx, tx, userID)

        // Send success response
        w.WriteHeader(http.StatusOK)
        w.Header().Set("Content-Type", "application/json")
        w.Write([]byte(`{"message": "Password updated successfully"}`))
    }
}
