package middleware

import (
    "fmt"
    "net/http"
    "os"
    "strconv"
    "golang.org/x/time/rate"
    "github.com/joho/godotenv"
)

var limiter *rate.Limiter

func init() {
    // Load environment variables from .env file
    if err := godotenv.Load(); err != nil {
        panic("Error loading .env file")
    }

    // Get rate limit settings from environment variables
    rateLimitPerSecond, err := strconv.Atoi(os.Getenv("RATE_LIMIT_PER_SECOND"))
    if err != nil {
        panic(fmt.Sprintf("Invalid RATE_LIMIT_PER_SECOND value: %v", err))
    }

    burstSize, err := strconv.Atoi(os.Getenv("BURST_SIZE"))
    if err != nil {
        panic(fmt.Sprintf("Invalid BURST_SIZE value: %v", err))
    }

    limiter = rate.NewLimiter(rate.Limit(rateLimitPerSecond), burstSize)
}

// RateLimiterMiddleware enforces the global rate limit for all endpoints
func RateLimiterMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        ctx := r.Context()

        // Apply the rate limit using the request's context
        if err := limiter.Wait(ctx); err != nil {
            http.Error(w, fmt.Sprintf("Rate limit exceeded: %v", err), http.StatusTooManyRequests)
            return
        }

        // Proceed to the next handler if the rate limit is not exceeded
        next.ServeHTTP(w, r)
    })
}
