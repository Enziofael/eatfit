// internal/model/account.go
package model

import (
	"time"

	"github.com/google/uuid"
)

// VerificationStatus - статус верификации email
type VerificationStatus string

const (
	VerificationPending  VerificationStatus = "pending"
	VerificationVerified VerificationStatus = "verified"
	VerificationExpired  VerificationStatus = "expired"
)

// Account - учётная запись пользователя (ядро аутентификации)
// Содержит только данные, необходимые для входа и безопасности.
// Профильные данные (имя, аватар, био) будут в отдельной сущности Profile.
type Account struct {
	ID    uuid.UUID
	Email string
	Login string

	// Безопасность
	PasswordHash       string
	EmailVerified      bool
	VerificationStatus VerificationStatus

	// Код верификации (хранится только до подтверждения email)
	VerificationCode          string
	VerificationCodeExpiresAt *time.Time
	VerificationAttempts      int

	// Временные метки
	CreatedAt   time.Time
	UpdatedAt   time.Time
	LastLoginAt *time.Time
}

// RegisterAccountRequest - данные для регистрации нового аккаунта
type RegisterAccountRequest struct {
	Email                string
	Login                string
	Password             string
	PasswordConfirmation string
}

// LoginRequest - данные для входа
type LoginRequest struct {
	LoginIdentifier string // email или login
	Password        string
	DeviceInfo      string
}

// LoginResult - результат успешного входа
type LoginResult struct {
	AccessToken  string
	RefreshToken string
	ExpiresIn    int64
	TokenType    string
	Account      *Account
}

// RefreshToken - сессионный токен (refresh token)
type RefreshToken struct {
	ID         uuid.UUID
	AccountID  uuid.UUID
	TokenHash  string
	DeviceInfo string
	IPAddress  string
	UserAgent  string
	CreatedAt  time.Time
	ExpiresAt  time.Time
	RevokedAt  *time.Time
	IsRevoked  bool
}

// VerificationCode - код подтверждения email
type VerificationCode struct {
	AccountID uuid.UUID
	Code      string
	ExpiresAt time.Time
	Attempts  int
}
