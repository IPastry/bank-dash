package middleware

import (
    "net/http"
    "os"
)

// CORSHeaderMiddleware adds CORS headers to the response
func CORSHeaderMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Restrict allowed origins from environment variable
        allowedOrigin := os.Getenv("ALLOWED_ORIGIN")
        if allowedOrigin == "" {
            allowedOrigin = "*"
        }

        w.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

        // Handle preflight requests
        if r.Method == http.MethodOptions {
            w.WriteHeader(http.StatusOK)
            return
        }

        next.ServeHTTP(w, r)
    })
}
