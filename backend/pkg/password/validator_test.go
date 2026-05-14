// backend/internal/pkg/password/validator_test.go
package password_test

import (
	"testing"

	"github.com/Enziofael/eatfit/backend/internal/pkg/password"
	"github.com/stretchr/testify/assert"
)

func TestPasswordValidator(t *testing.T) {
	validator := password.NewPasswordValidator()

	tests := []struct {
		name     string
		password string
		wantErr  bool
		errMsg   string
	}{
		{
			name:     "Валидный пароль",
			password: "Test123!",
			wantErr:  false,
		},
		{
			name:     "Слишком короткий пароль",
			password: "T1!a",
			wantErr:  true,
			errMsg:   "password must be at least 5 characters",
		},
		{
			name:     "Без заглавной буквы",
			password: "test123!",
			wantErr:  true,
			errMsg:   "password must contain at least one uppercase letter",
		},
		{
			name:     "Без строчной буквы",
			password: "TEST123!",
			wantErr:  true,
			errMsg:   "password must contain at least one lowercase letter",
		},
		{
			name:     "Без цифры",
			password: "TestTest!",
			wantErr:  true,
			errMsg:   "password must contain at least one digit",
		},
		{
			name:     "Без спецсимвола",
			password: "Test1234",
			wantErr:  true,
			errMsg:   "password must contain at least one special character",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.Validate(tt.password)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
