package models

type Role string

const (
    Admin    Role = "admin"
    Elevated Role = "elevated"
    Regular  Role = "regular"
)