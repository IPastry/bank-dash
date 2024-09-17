package models

import "time"

type Users struct {
    UserID       int        `json:"user_id"`
    Email        string     `json:"email"`
    PasswordHash string     `json:"password_hash"`
    Role         Role       `json:"role"` 
    ManagerID    *int       `json:"manager_id"`
    BankID       *int       `json:"bank_id"`
    FirstName    *string    `json:"first_name"`
    LastName     *string    `json:"last_name"`
    Preferences  *string    `json:"preferences"` // JSON as string for simplicity
    IsActive     bool       `json:"is_active"`
    CreatedAt    time.Time  `json:"created_at"`
    UpdatedAt    time.Time  `json:"updated_at"`
    LastLogin    *time.Time `json:"last_login"`
		IsVerified   bool       `json:"is_verified"`
		PhoneNumber  *string    `json:"phone_number"`
}
