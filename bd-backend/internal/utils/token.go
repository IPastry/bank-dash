package utils

import (
    "crypto/rand"
    "encoding/hex"
)

// GenerateEmailVerificationToken generates a unique token for email verification
func GenerateEmailVerificationToken() (string, error) {
    token := make([]byte, 16)
    _, err := rand.Read(token)
    if err != nil {
        return "", err
    }
    return hex.EncodeToString(token), nil
}
