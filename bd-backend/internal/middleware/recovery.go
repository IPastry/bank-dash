package middleware

import (
    "net/http"
    "log"
)

// RecoveryMiddleware recovers from panics and logs them with request details
func RecoveryMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        defer func() {
            if err := recover(); err != nil {
                log.Printf("Recovered from panic: %v | Method: %s, URL: %s", err, r.Method, r.URL.Path)
                http.Error(w, "Internal Server Error", http.StatusInternalServerError)
            }
        }()
        next.ServeHTTP(w, r)
    })
}
