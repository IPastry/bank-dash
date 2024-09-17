package routes

import (
	"bd-backend/internal/handlers"
	"bd-backend/internal/middleware"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
)

// SetupRoutes initializes all the routes for your application
func SetupRoutes(dbPool *pgxpool.Pool) *mux.Router {
	router := mux.NewRouter()

	// Define and apply global middleware
	applyGlobalMiddleware(router)

	// Define route groups
	setupAuthRoutes(router, dbPool)
	setupUserRoutes(router, dbPool)
	// setupAdminRoutes(router, dbPool)

	return router
}

// applyGlobalMiddleware applies middleware that should be applied to all routes
func applyGlobalMiddleware(router *mux.Router) {
	loggingMiddleware := middleware.LoggingMiddleware
	corsMiddleware := middleware.CORSHeaderMiddleware
	recoveryMiddleware := middleware.RecoveryMiddleware
	timeoutMiddleware := middleware.TimeoutMiddleware(time.Second * 10) // 10 seconds timeout
	rateLimiterMiddleware := middleware.RateLimiterMiddleware

	router.Use(loggingMiddleware)     // Log every request
	router.Use(corsMiddleware)        // Set CORS headers
	router.Use(recoveryMiddleware)    // Handle panics safely
	router.Use(timeoutMiddleware)     // Set a timeout limit (e.g., 10 seconds)
	router.Use(rateLimiterMiddleware) // Apply rate limiting globally
}

// setupAuthRoutes sets up authentication-related routes
func setupAuthRoutes(router *mux.Router, dbPool *pgxpool.Pool) {
	router.HandleFunc("/login", handlers.LoginHandler(dbPool)).Methods(http.MethodPost)
	router.HandleFunc("/resend-verification-email", handlers.EmailResendHandler(dbPool)).Methods(http.MethodPost)
}

// setupUserRoutes sets up user-related routes
func setupUserRoutes(router *mux.Router, dbPool *pgxpool.Pool) {
	router.HandleFunc("/elevated-user", handlers.CreateElevatedUserHandler(dbPool)).Methods(http.MethodPost)
	router.HandleFunc("/update-profile", handlers.UpdateProfileHandler(dbPool)).Methods(http.MethodPost)
	router.HandleFunc("/bank-info/{query}", handlers.GetBankInfoHandler(dbPool)).Methods(http.MethodGet)
	router.HandleFunc("/confirm-bank", handlers.ConfirmBankHandler(dbPool)).Methods(http.MethodPost)
	router.HandleFunc("/shared-account", handlers.CreateSharedAccountHandler(dbPool)).Methods(http.MethodPost)

	// Protected routes (authentication required)
	protected := router.PathPrefix("/").Subrouter()
	protected.Use(middleware.AuthMiddleware)
}

// // setupAdminRoutes sets up admin-related routes
// func setupAdminRoutes(router *mux.Router, dbPool *pgxpool.Pool) {
// 	admin := router.PathPrefix("/admin").Subrouter()
// 	admin.Use(middleware.AuthMiddleware) // Ensure admin routes are protected
// 	admin.Use(func(next http.Handler) http.Handler {
// 		return middleware.RoleMiddleware("admin", next)
// 	})
// }
