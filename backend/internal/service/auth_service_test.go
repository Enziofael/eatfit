// backend/internal/service/auth_service_test.go
package service_test

import (
	"context"
	"testing"

	"github.com/Enziofael/eatfit/backend/internal/service"
	"github.com/stretchr/testify/assert"
)

// TestRegisterUser_Success - успешная регистрация нового пользователя
func TestRegisterUser_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(MockUserRepository)
	mockEmail := new(MockEmailSender)
	mockHasher := new(MockPasswordHasher)

	service := service.NewAuthService(mockRepo, mockEmail, mockHasher)

	req := &RegisterRequest{
		Email:    "test@example.com",
		Login:    "testuser",
		Password: "Test123!",
	}

	// Assert
	assert.NoError(t, err)
	assert.NotEmpty(t, response.UserId)
}

// TestRegisterUser_DuplicateEmail - ошибка при дублировании email
func TestRegisterUser_DuplicateEmail(t *testing.T) {
	// ... тест на дубликат email
}

// TestRegisterUser_WeakPassword - ошибка при слабом пароле
func TestRegisterUser_WeakPassword(t *testing.T) {
	// ... тест на валидацию пароля
}

// TestLogin_Success - успешный вход в систему
func TestLogin_Success(t *testing.T) {
	// ... тест на успешный вход
}

// TestLogin_InvalidCredentials - ошибка при неверных учётных данных
func TestLogin_InvalidCredentials(t *testing.T) {
	// ... тест на неверный пароль
}
