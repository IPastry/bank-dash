package middleware

import (
	"bd-backend/internal/utils"
	"context"
	"net/http"
	"time"

	"log"

	"github.com/gorilla/mux"
)

// TimeoutMiddleware sets a timeout for handling requests
func TimeoutMiddleware(timeout time.Duration) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx, cancel := context.WithTimeout(r.Context(), timeout)
			defer cancel()

			// Wrap the ResponseWriter to log timeout events
			rw := &utils.ResponseWriterWithStatusCode{
				ResponseWriter: w,
				StatusCode:     http.StatusOK, // default status code
			}

			r = r.WithContext(ctx)
			done := make(chan bool, 1)

			go func() {
				next.ServeHTTP(rw, r)
				done <- true
			}()

			select {
			case <-ctx.Done():
				if ctx.Err() == context.DeadlineExceeded {
					log.Printf("Request to %s timed out", r.URL.Path)
					http.Error(w, "Request timed out", http.StatusRequestTimeout)
				}
			case <-done:
				// Request completed before the timeout
			}
		})
	}
}
