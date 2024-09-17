package handlers

import (
	"bd-backend/internal/auth"
	"bd-backend/internal/db/users"
	"bd-backend/internal/models"
	"bd-backend/internal/utils"
	"encoding/json"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

// LoginHandler handles user login requests
func LoginHandler(dbPool *pgxpool.Pool) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var req struct {
            Email    string `json:"email"`
            Password string `json:"password"`
        }
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
            utils.RespondWithJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request"})
            return
        }

        // Get context from the request
        ctx := r.Context()

        // Start a transaction
        tx, err := dbPool.Begin(ctx)
        if err != nil {
            utils.RespondWithJSON(w, http.StatusInternalServerError, map[string]string{"error": "Error starting transaction"})
            return
        }
        defer tx.Rollback(ctx) // Rollback if there's any error

        // Verify user password and retrieve user ID and role
        valid, userID, role, err := users.VerifyUserPasswordAndRole(ctx, tx, req.Email, req.Password)
        if err != nil || !valid {
            utils.RespondWithJSON(w, http.StatusUnauthorized, map[string]string{"error": "Invalid credentials"})
            return
        }

        // Check profile completeness for Elevated role
        if role == models.Elevated {
            isComplete, err := users.IsProfileComplete(ctx, tx, userID)
            if err != nil {
                utils.RespondWithJSON(w, http.StatusInternalServerError, map[string]string{"error": "Error checking profile completeness"})
                return
            }
            if !isComplete {
                utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Profile incomplete. Please update your profile details."})
                return
            }
        }

        // Commit the transaction
        if err := tx.Commit(ctx); err != nil {
            utils.RespondWithJSON(w, http.StatusInternalServerError, map[string]string{"error": "Error committing transaction"})
            return
        }

        // Generate token for the user
        tokenString, err := auth.GenerateToken(userID, role)
        if err != nil {
            utils.RespondWithJSON(w, http.StatusInternalServerError, map[string]string{"error": "Error generating token"})
            return
        }

        // Send the token in response
        utils.RespondWithJSON(w, http.StatusOK, map[string]string{"token": tokenString})
    }
}


// RefreshTokenHandler handles token refresh requests
func RefreshTokenHandler(dbPool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			RefreshToken string `json:"refresh_token"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.RespondWithJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request"})
			return
		}

		_, claims, err := auth.ParseRefreshToken(req.RefreshToken)
		if err != nil {
			utils.RespondWithJSON(w, http.StatusUnauthorized, map[string]string{"error": "Invalid refresh token"})
			return
		}

		userID, ok := claims["user_id"].(float64)
		if !ok {
			utils.RespondWithJSON(w, http.StatusUnauthorized, map[string]string{"error": "Invalid token claims"})
			return
		}

		role, err := users.RetrieveUserRole(r.Context(), dbPool, int(userID))
		if err != nil {
			utils.RespondWithJSON(w, http.StatusInternalServerError, map[string]string{"error": "Error retrieving user role"})
			return
		}

		accessToken, err := auth.GenerateToken(int(userID), role)
		if err != nil {
			utils.RespondWithJSON(w, http.StatusInternalServerError, map[string]string{"error": "Error generating access token"})
			return
		}

		newRefreshToken, err := auth.GenerateRefreshToken(int(userID))
		if err != nil {
			utils.RespondWithJSON(w, http.StatusInternalServerError, map[string]string{"error": "Error generating refresh token"})
			return
		}

		utils.RespondWithJSON(w, http.StatusOK, map[string]string{"access_token": accessToken, "refresh_token": newRefreshToken})
	}
}
