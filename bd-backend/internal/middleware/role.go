package middleware

import (
    "bd-backend/internal/models"
    "bd-backend/internal/utils"
    "net/http"
    "log"
)

func RoleMiddleware(requiredRole models.Role, next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        role, ok := r.Context().Value(utils.RoleKey).(models.Role)
        if !ok || role != requiredRole {
            log.Printf("Access denied. User role: %s, Required role: %s", role, requiredRole)
            http.Error(w, "Forbidden", http.StatusForbidden)
            return
        }

        next.ServeHTTP(w, r)
    })
}
