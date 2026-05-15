// internal/service/auth_service.go
package service

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"strings"
	"time"

	emailValidator "github.com/Enziofael/eatfit/backend/pkg/email"
	loginValidator "github.com/Enziofael/eatfit/backend/pkg/login"
	"github.com/Enziofael/eatfit/backend/pkg/password"

	"github.com/Enziofael/eatfit/backend/internal/model"
	"github.com/Enziofael/eatfit/backend/internal/repository/postgres"
	redisRepo "github.com/Enziofael/eatfit/backend/internal/repository/redis"
	"github.com/google/uuid"
)

// AuthService - сервис аутентификации и регистрации
type AuthService struct {
	accountRepo       *postgres.AccountRepo
	tokenRepo         *redisRepo.TokenRepo
	jwtService        *JWTService
	emailService      *EmailService
	passwordHasher    *password.Hasher
	passwordValidator *password.Validator
	emailValidator    *emailValidator.Validator
	loginValidator    *loginValidator.Validator
}

// NewAuthService создаёт новый экземпляр AuthService
func NewAuthService(
	accountRepo *postgres.AccountRepo,
	tokenRepo *redisRepo.TokenRepo,
	jwtService *JWTService,
	emailService *EmailService,
	passwordHasher *password.Hasher,
) *AuthService {
	return &AuthService{
		accountRepo:       accountRepo,
		tokenRepo:         tokenRepo,
		jwtService:        jwtService,
		emailService:      emailService,
		passwordHasher:    passwordHasher,
		passwordValidator: password.NewValidator(),
		emailValidator:    emailValidator.NewValidator(),
		loginValidator:    loginValidator.NewValidator(),
	}
}

// Register регистрирует нового пользователя
func (s *AuthService) Register(ctx context.Context, req *model.RegisterAccountRequest) (*model.Account, string, error) {
	// Нормализуем email
	req.Email = emailValidator.NormalizeEmail(req.Email)

	// Валидация email
	if err := s.emailValidator.Validate(req.Email); err != nil {
		return nil, "", fmt.Errorf("invalid email: %w", err)
	}

	// Валидация логина
	if err := s.loginValidator.Validate(req.Login); err != nil {
		return nil, "", fmt.Errorf("invalid login: %w", err)
	}

	// Валидация пароля
	if err := s.passwordValidator.ValidatePassword(req.Password); err != nil {
		return nil, "", fmt.Errorf("invalid password: %w", err)
	}

	// Проверка совпадения паролей
	if req.Password != req.PasswordConfirmation {
		return nil, "", fmt.Errorf("password and confirmation do not match")
	}

	// Проверка уникальности email
	emailExists, err := s.accountRepo.IsEmailExists(ctx, req.Email)
	if err != nil {
		return nil, "", fmt.Errorf("failed to check email: %w", err)
	}
	if emailExists {
		return nil, "", fmt.Errorf("email already exists")
	}

	// Проверка уникальности логина
	loginExists, err := s.accountRepo.IsLoginExists(ctx, req.Login)
	if err != nil {
		return nil, "", fmt.Errorf("failed to check login: %w", err)
	}
	if loginExists {
		return nil, "", fmt.Errorf("login already exists")
	}

	// Хеширование пароля
	passwordHash, err := s.passwordHasher.HashPassword(req.Password)
	if err != nil {
		return nil, "", fmt.Errorf("failed to hash password: %w", err)
	}

	// Создание аккаунта
	account, err := s.accountRepo.CreateAccount(ctx, req, passwordHash)
	if err != nil {
		return nil, "", fmt.Errorf("failed to create account: %w", err)
	}

	// Генерация кода верификации
	verificationCode, err := s.generateVerificationCode()
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate verification code: %w", err)
	}

	// Сохранение кода верификации (в БД и Redis для быстрого доступа)
	expiresAt := time.Now().Add(15 * time.Minute)
	if err := s.accountRepo.UpdateVerificationCode(ctx, account.ID, verificationCode, expiresAt); err != nil {
		return nil, "", fmt.Errorf("failed to save verification code: %w", err)
	}

	if err := s.tokenRepo.SaveVerificationCode(ctx, account.ID, verificationCode, expiresAt); err != nil {
		return nil, "", fmt.Errorf("failed to cache verification code: %w", err)
	}

	// Отправка кода на email
	if err := s.emailService.SendVerificationEmail(account.Email, verificationCode); err != nil {
		return nil, "", fmt.Errorf("failed to send verification email: %w", err)
	}

	return account, verificationCode, nil
}

// VerifyEmail подтверждает email пользователя
func (s *AuthService) VerifyEmail(ctx context.Context, accountID uuid.UUID, code string) error {
	// Получаем аккаунт
	account, err := s.accountRepo.GetAccountByID(ctx, accountID)
	if err != nil {
		return fmt.Errorf("account not found: %w", err)
	}

	// Проверяем, не подтверждён ли уже email
	if account.EmailVerified {
		return fmt.Errorf("email already verified")
	}

	// Проверяем количество попыток
	if account.VerificationAttempts >= 5 {
		return fmt.Errorf("too many verification attempts, request a new code")
	}

	// Проверяем срок действия кода
	if account.VerificationCodeExpiresAt != nil && time.Now().After(*account.VerificationCodeExpiresAt) {
		return fmt.Errorf("verification code expired, request a new one")
	}

	// Проверяем код (сначала из Redis, потом из БД)
	storedCode, err := s.tokenRepo.GetVerificationCode(ctx, accountID)
	if err != nil {
		// Если в Redis нет, проверяем в БД
		if account.VerificationCode != code {
			_ = s.accountRepo.IncrementVerificationAttempts(ctx, accountID)
			return fmt.Errorf("invalid verification code")
		}
	} else if storedCode != code {
		_ = s.accountRepo.IncrementVerificationAttempts(ctx, accountID)
		return fmt.Errorf("invalid verification code")
	}

	// Подтверждаем email
	if err := s.accountRepo.VerifyEmail(ctx, accountID); err != nil {
		return fmt.Errorf("failed to verify email: %w", err)
	}

	// Удаляем код из Redis
	_ = s.tokenRepo.DeleteVerificationCode(ctx, accountID)

	// Отправляем приветственное письмо
	_ = s.emailService.SendWelcomeEmail(account.Email, account.Login)

	return nil
}

// Login выполняет вход пользователя
func (s *AuthService) Login(ctx context.Context, req *model.LoginRequest) (*model.LoginResult, error) {
	// Ищем пользователя по email или логину
	var account *model.Account
	var err error

	// Определяем, что ввёл пользователь: email или логин
	if emailValidator.NewValidator().Validate(req.LoginIdentifier) == nil {
		account, err = s.accountRepo.GetAccountByEmail(ctx, req.LoginIdentifier)
	} else {
		account, err = s.accountRepo.GetAccountByLogin(ctx, req.LoginIdentifier)
	}

	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Проверяем, подтверждён ли email
	if !account.EmailVerified {
		return nil, fmt.Errorf("email not verified, please verify your email first")
	}

	// Проверяем пароль
	if !s.passwordHasher.VerifyPassword(req.Password, account.PasswordHash) {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Генерируем access токен
	accessToken, expiresAt, err := s.jwtService.GenerateAccessToken(account)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	// Генерируем refresh токен
	refreshToken, refreshExpiresAt, err := s.jwtService.GenerateRefreshToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// Сохраняем refresh токен в Redis
	tokenHash := redisRepo.HashToken(refreshToken)
	if err := s.tokenRepo.SaveRefreshToken(ctx, account.ID, tokenHash, refreshExpiresAt); err != nil {
		return nil, fmt.Errorf("failed to save refresh token: %w", err)
	}

	// Обновляем время последнего входа
	_ = s.accountRepo.UpdateLastLogin(ctx, account.ID)

	return &model.LoginResult{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(time.Until(expiresAt).Seconds()),
		TokenType:    "Bearer",
		Account:      account,
	}, nil
}

// RefreshToken обновляет токены доступа
func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (*model.LoginResult, error) {
	// Хешируем токен для поиска в Redis
	tokenHash := redisRepo.HashToken(refreshToken)

	// Получаем ID аккаунта по токену
	accountID, err := s.tokenRepo.GetRefreshToken(ctx, tokenHash)
	if err != nil {
		return nil, fmt.Errorf("invalid or expired refresh token")
	}

	// Отзываем старый токен
	_ = s.tokenRepo.RevokeRefreshToken(ctx, tokenHash)

	// Получаем аккаунт
	account, err := s.accountRepo.GetAccountByID(ctx, accountID)
	if err != nil {
		return nil, fmt.Errorf("account not found")
	}

	// Генерируем новые токены
	accessToken, expiresAt, err := s.jwtService.GenerateAccessToken(account)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	newRefreshToken, refreshExpiresAt, err := s.jwtService.GenerateRefreshToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// Сохраняем новый refresh токен
	newTokenHash := redisRepo.HashToken(newRefreshToken)
	if err := s.tokenRepo.SaveRefreshToken(ctx, account.ID, newTokenHash, refreshExpiresAt); err != nil {
		return nil, fmt.Errorf("failed to save refresh token: %w", err)
	}

	return &model.LoginResult{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		ExpiresIn:    int64(time.Until(expiresAt).Seconds()),
		TokenType:    "Bearer",
		Account:      account,
	}, nil
}

// Logout выполняет выход из системы
func (s *AuthService) Logout(ctx context.Context, refreshToken string) error {
	tokenHash := redisRepo.HashToken(refreshToken)
	return s.tokenRepo.RevokeRefreshToken(ctx, tokenHash)
}

// ResendVerification отправляет новый код верификации
func (s *AuthService) ResendVerification(ctx context.Context, accountID uuid.UUID) error {
	// Получаем аккаунт
	account, err := s.accountRepo.GetAccountByID(ctx, accountID)
	if err != nil {
		return fmt.Errorf("account not found: %w", err)
	}

	// Проверяем, не подтверждён ли уже email
	if account.EmailVerified {
		return fmt.Errorf("email already verified")
	}

	// Генерируем новый код
	verificationCode, err := s.generateVerificationCode()
	if err != nil {
		return fmt.Errorf("failed to generate verification code: %w", err)
	}

	// Сохраняем код
	expiresAt := time.Now().Add(15 * time.Minute)
	if err := s.accountRepo.UpdateVerificationCode(ctx, accountID, verificationCode, expiresAt); err != nil {
		return fmt.Errorf("failed to save verification code: %w", err)
	}

	if err := s.tokenRepo.SaveVerificationCode(ctx, accountID, verificationCode, expiresAt); err != nil {
		return fmt.Errorf("failed to cache verification code: %w", err)
	}

	// Отправляем код
	if err := s.emailService.SendVerificationEmail(account.Email, verificationCode); err != nil {
		return fmt.Errorf("failed to send verification email: %w", err)
	}

	return nil
}

// generateVerificationCode генерирует 6-значный код подтверждения
func (s *AuthService) generateVerificationCode() (string, error) {
	code := ""
	for i := 0; i < 6; i++ {
		n, err := rand.Int(rand.Reader, big.NewInt(10))
		if err != nil {
			return "", fmt.Errorf("failed to generate random number: %w", err)
		}
		code += n.String()
	}
	return code, nil
}

// ForgotPassword отправляет код сброса пароля
func (s *AuthService) ForgotPassword(ctx context.Context, loginIdentifier string) (string, error) {
	// Ищем пользователя
	var account *model.Account
	var err error

	if emailValidator.NewValidator().Validate(loginIdentifier) == nil {
		account, err = s.accountRepo.GetAccountByEmail(ctx, loginIdentifier)
	} else {
		account, err = s.accountRepo.GetAccountByLogin(ctx, loginIdentifier)
	}

	if err != nil {
		// Не раскрываем, существует ли аккаунт
		return "", nil
	}

	// Генерируем код
	verificationCode, err := s.generateVerificationCode()
	if err != nil {
		return "", fmt.Errorf("failed to generate code: %w", err)
	}

	// Генерируем reset токен
	resetToken := uuid.New().String()

	// Сохраняем код в Redis с привязкой к reset токену
	expiresAt := time.Now().Add(15 * time.Minute)
	key := fmt.Sprintf("password_reset:%s", resetToken)
	// Сохраняем: account_id, code, новый пароль пока пустой
	data := fmt.Sprintf("%s|%s", account.ID.String(), verificationCode)
	if err := s.tokenRepo.SavePasswordResetToken(ctx, key, data, expiresAt); err != nil {
		return "", fmt.Errorf("failed to save reset token: %w", err)
	}

	// Отправляем код на email
	if err := s.emailService.SendPasswordResetEmail(account.Email, verificationCode); err != nil {
		return "", fmt.Errorf("failed to send email: %w", err)
	}

	return resetToken, nil
}

// ResetPassword сбрасывает пароль
func (s *AuthService) ResetPassword(ctx context.Context, resetToken, code, newPassword, confirmation string) error {
	// Валидация пароля
	if err := s.passwordValidator.ValidatePassword(newPassword); err != nil {
		return err
	}

	if newPassword != confirmation {
		return fmt.Errorf("password and confirmation do not match")
	}

	// Получаем данные из Redis
	key := fmt.Sprintf("password_reset:%s", resetToken)
	data, err := s.tokenRepo.GetPasswordResetToken(ctx, key)
	if err != nil {
		return fmt.Errorf("invalid or expired reset token")
	}

	// Парсим данные: account_id|code
	parts := strings.Split(data, "|")
	if len(parts) != 2 {
		return fmt.Errorf("invalid reset token data")
	}

	accountID, err := uuid.Parse(parts[0])
	if err != nil {
		return fmt.Errorf("invalid account id in token")
	}

	storedCode := parts[1]
	if storedCode != code {
		return fmt.Errorf("invalid verification code")
	}

	// Хешируем новый пароль
	passwordHash, err := s.passwordHasher.HashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// Обновляем пароль
	if err := s.accountRepo.UpdatePassword(ctx, accountID, passwordHash); err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	// Удаляем reset токен
	_ = s.tokenRepo.DeletePasswordResetToken(ctx, key)

	return nil
}
