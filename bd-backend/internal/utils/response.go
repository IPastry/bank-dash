// utils/response.go
package utils

import (
    "encoding/json"
    "net/http"
)
// ResponseWriterWithStatusCode wraps the http.ResponseWriter to capture the status code
type ResponseWriterWithStatusCode struct {
    http.ResponseWriter
    StatusCode int
}

// WriteHeader overrides the default WriteHeader method to capture the status code
func (rw *ResponseWriterWithStatusCode) WriteHeader(statusCode int) {
    rw.StatusCode = statusCode
    rw.ResponseWriter.WriteHeader(statusCode)
}

// RespondWithJSON sends a JSON response with the given status code and data.
func RespondWithJSON(w http.ResponseWriter, statusCode int, data interface{}) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(statusCode)
    json.NewEncoder(w).Encode(data)
}