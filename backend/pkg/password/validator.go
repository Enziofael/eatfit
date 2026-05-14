// pkg/password/validator.go
package password

import (
	"fmt"
	"regexp"
)

// ValidationError описывает ошибку валидации пароля
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

// Validator - валидатор паролей
type Validator struct {
	minLength          int
	requireUppercase   bool
	requireLowercase   bool
	requireDigit       bool
	requireSpecialChar bool
	specialChars       string
}

// NewValidator создаёт валидатор с настройками по умолчанию
func NewValidator() *Validator {
	return &Validator{
		minLength:          5,
		requireUppercase:   true,
		requireLowercase:   true,
		requireDigit:       true,
		requireSpecialChar: true,
		specialChars:       `[!@#\$%\^&\*\(\),\.\?":\{\}\|<>]`,
	}
}

// ValidatePassword проверяет пароль на соответствие требованиям
func (v *Validator) ValidatePassword(password string) error {
	if len(password) < v.minLength {
		return &ValidationError{
			Field:   "password",
			Message: fmt.Sprintf("password must be at least %d characters long", v.minLength),
		}
	}

	if v.requireUppercase && !regexp.MustCompile(`[A-Z]`).MatchString(password) {
		return &ValidationError{
			Field:   "password",
			Message: "password must contain at least one uppercase letter",
		}
	}

	if v.requireLowercase && !regexp.MustCompile(`[a-z]`).MatchString(password) {
		return &ValidationError{
			Field:   "password",
			Message: "password must contain at least one lowercase letter",
		}
	}

	if v.requireDigit && !regexp.MustCompile(`[0-9]`).MatchString(password) {
		return &ValidationError{
			Field:   "password",
			Message: "password must contain at least one digit",
		}
	}

	if v.requireSpecialChar && !regexp.MustCompile(v.specialChars).MatchString(password) {
		return &ValidationError{
			Field:   "password",
			Message: "password must contain at least one special character",
		}
	}

	return nil
}
