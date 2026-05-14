package domain

import (
    "time"
    "github.com/google/uuid"
)

type User struct {
    ID              uuid.UUID
    Email           string
    Username        string
    PasswordHash    string
    EmailVerified   bool
    CreatedAt       time.Time
    UpdatedAt       time.Time
}

type Session struct {
    ID           uuid.UUID
    UserID       uuid.UUID
    RefreshToken string
    UserAgent    string
    IPAddress    string
    ExpiresAt    time.Time
    CreatedAt    time.Time
    Revoked      bool
}