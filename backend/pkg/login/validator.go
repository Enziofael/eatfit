// pkg/login/validator.go
package login

import (
	"fmt"
	"regexp"
)

var loginRegex = regexp.MustCompile(`^[a-zA-Z0-9_]{3,30}$`)

// Validator - валидатор логинов
type Validator struct {
	minLength int
	maxLength int
}

// NewValidator создаёт новый валидатор логинов
func NewValidator() *Validator {
	return &Validator{
		minLength: 3,
		maxLength: 30,
	}
}

// Validate проверяет корректность логина
// Возвращает nil если логин валиден, иначе ошибку с описанием
func (v *Validator) Validate(login string) error {
	if login == "" {
		return fmt.Errorf("login is required")
	}

	if len(login) < v.minLength {
		return fmt.Errorf("login must be at least %d characters long", v.minLength)
	}

	if len(login) > v.maxLength {
		return fmt.Errorf("login must be at most %d characters long", v.maxLength)
	}

	if !loginRegex.MatchString(login) {
		return fmt.Errorf("login must contain only latin letters, digits and underscore")
	}

	return nil
}
