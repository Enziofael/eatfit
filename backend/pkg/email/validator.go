// pkg/email/validator.go
package email

import (
	"fmt"
	"regexp"
	"strings"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

// Validator - валидатор email адресов
type Validator struct {
	maxLength int
}

// NewValidator создаёт новый валидатор email
func NewValidator() *Validator {
	return &Validator{
		maxLength: 255,
	}
}

// Validate проверяет корректность email адреса
// Возвращает nil если email валиден, иначе ошибку с описанием
func (v *Validator) Validate(email string) error {
	if email == "" {
		return fmt.Errorf("email is required")
	}

	if len(email) > v.maxLength {
		return fmt.Errorf("email is too long (max %d characters)", v.maxLength)
	}

	if !emailRegex.MatchString(email) {
		return fmt.Errorf("invalid email format")
	}

	return nil
}

// NormalizeEmail приводит email к нижнему регистру и удаляет пробелы
func NormalizeEmail(email string) string {
	return strings.TrimSpace(strings.ToLower(email))
}
