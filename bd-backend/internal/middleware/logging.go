package middleware

import (
	"bd-backend/internal/utils"
	"log"
	"net/http"
	"time"
)

// LoggingMiddleware logs details of each request
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		rw := &utils.ResponseWriterWithStatusCode{
    ResponseWriter: w,
    StatusCode:     http.StatusOK, // default status code
}

		// Log request details
		log.Printf("Request received: Method: %s, URL: %s, Headers: %v",
			r.Method,
			r.URL.Path,
			r.Header,
		)

		next.ServeHTTP(rw, r)

		duration := time.Since(start)
		log.Printf(
			"Request completed: Method: %s, URL: %s, Status: %d, Duration: %s",
			r.Method,
			r.URL.Path,
			rw.StatusCode,
			duration,
		)
	})
}
