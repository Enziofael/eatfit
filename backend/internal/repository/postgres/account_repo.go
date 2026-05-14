// internal/repository/postgres/account_repo.go
package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/Enziofael/eatfit/backend/internal/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// AccountRepo - репозиторий для работы с аккаунтами в PostgreSQL
type AccountRepo struct {
	pool *pgxpool.Pool
}

// NewAccountRepo создаёт новый экземпляр AccountRepo
func NewAccountRepo(pool *pgxpool.Pool) *AccountRepo {
	return &AccountRepo{pool: pool}
}

// CreateAccount создаёт новую учётную запись
func (r *AccountRepo) CreateAccount(ctx context.Context, req *model.RegisterAccountRequest, passwordHash string) (*model.Account, error) {
	query := `
		INSERT INTO users (email, login, password_hash)
		VALUES ($1, $2, $3)
		RETURNING id, email, login, password_hash, email_verified, 
		          verification_status, verification_code, 
		          verification_code_expires_at, verification_attempts,
		          created_at, updated_at, last_login_at
	`

	account := &model.Account{}
	var verificationCode *string
	var verificationCodeExpiresAt *time.Time
	var lastLoginAt *time.Time

	err := r.pool.QueryRow(ctx, query, req.Email, req.Login, passwordHash).Scan(
		&account.ID,
		&account.Email,
		&account.Login,
		&account.PasswordHash,
		&account.EmailVerified,
		&account.VerificationStatus,
		&verificationCode,
		&verificationCodeExpiresAt,
		&account.VerificationAttempts,
		&account.CreatedAt,
		&account.UpdatedAt,
		&lastLoginAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create account: %w", err)
	}

	// Преобразуем *string в string
	if verificationCode != nil {
		account.VerificationCode = *verificationCode
	}
	account.VerificationCodeExpiresAt = verificationCodeExpiresAt
	account.LastLoginAt = lastLoginAt

	return account, nil
}

// GetAccountByEmail находит аккаунт по email
func (r *AccountRepo) GetAccountByEmail(ctx context.Context, email string) (*model.Account, error) {
	query := `
		SELECT id, email, login, password_hash, email_verified, 
		       verification_status, verification_code,
		       verification_code_expires_at, verification_attempts,
		       created_at, updated_at, last_login_at
		FROM users
		WHERE email = $1
	`

	account := &model.Account{}
	var verificationCode *string
	var verificationCodeExpiresAt *time.Time
	var lastLoginAt *time.Time

	err := r.pool.QueryRow(ctx, query, email).Scan(
		&account.ID,
		&account.Email,
		&account.Login,
		&account.PasswordHash,
		&account.EmailVerified,
		&account.VerificationStatus,
		&verificationCode,
		&verificationCodeExpiresAt,
		&account.VerificationAttempts,
		&account.CreatedAt,
		&account.UpdatedAt,
		&lastLoginAt,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("account not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get account by email: %w", err)
	}

	if verificationCode != nil {
		account.VerificationCode = *verificationCode
	}
	account.VerificationCodeExpiresAt = verificationCodeExpiresAt
	account.LastLoginAt = lastLoginAt

	return account, nil
}

// GetAccountByLogin находит аккаунт по логину
func (r *AccountRepo) GetAccountByLogin(ctx context.Context, login string) (*model.Account, error) {
	query := `
		SELECT id, email, login, password_hash, email_verified,
		       verification_status, verification_code,
		       verification_code_expires_at, verification_attempts,
		       created_at, updated_at, last_login_at
		FROM users
		WHERE login = $1
	`

	account := &model.Account{}
	var verificationCode *string
	var verificationCodeExpiresAt *time.Time
	var lastLoginAt *time.Time

	err := r.pool.QueryRow(ctx, query, login).Scan(
		&account.ID,
		&account.Email,
		&account.Login,
		&account.PasswordHash,
		&account.EmailVerified,
		&account.VerificationStatus,
		&verificationCode,
		&verificationCodeExpiresAt,
		&account.VerificationAttempts,
		&account.CreatedAt,
		&account.UpdatedAt,
		&lastLoginAt,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("account not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get account by login: %w", err)
	}

	if verificationCode != nil {
		account.VerificationCode = *verificationCode
	}
	account.VerificationCodeExpiresAt = verificationCodeExpiresAt
	account.LastLoginAt = lastLoginAt

	return account, nil
}

// GetAccountByID находит аккаунт по ID
func (r *AccountRepo) GetAccountByID(ctx context.Context, id uuid.UUID) (*model.Account, error) {
	query := `
		SELECT id, email, login, password_hash, email_verified,
		       verification_status, verification_code,
		       verification_code_expires_at, verification_attempts,
		       created_at, updated_at, last_login_at
		FROM users
		WHERE id = $1
	`

	account := &model.Account{}
	var verificationCode *string
	var verificationCodeExpiresAt *time.Time
	var lastLoginAt *time.Time

	err := r.pool.QueryRow(ctx, query, id).Scan(
		&account.ID,
		&account.Email,
		&account.Login,
		&account.PasswordHash,
		&account.EmailVerified,
		&account.VerificationStatus,
		&verificationCode,
		&verificationCodeExpiresAt,
		&account.VerificationAttempts,
		&account.CreatedAt,
		&account.UpdatedAt,
		&lastLoginAt,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("account not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get account by id: %w", err)
	}

	if verificationCode != nil {
		account.VerificationCode = *verificationCode
	}
	account.VerificationCodeExpiresAt = verificationCodeExpiresAt
	account.LastLoginAt = lastLoginAt

	return account, nil
}

// UpdateVerificationCode обновляет код верификации
func (r *AccountRepo) UpdateVerificationCode(ctx context.Context, id uuid.UUID, code string, expiresAt time.Time) error {
	query := `
		UPDATE users
		SET verification_code = $2,
		    verification_code_expires_at = $3,
		    verification_attempts = 0,
		    updated_at = NOW()
		WHERE id = $1
	`

	_, err := r.pool.Exec(ctx, query, id, code, expiresAt)
	if err != nil {
		return fmt.Errorf("failed to update verification code: %w", err)
	}

	return nil
}

// VerifyEmail подтверждает email аккаунта
func (r *AccountRepo) VerifyEmail(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE users
		SET email_verified = TRUE,
		    verification_status = $2,
		    verification_code = NULL,
		    verification_code_expires_at = NULL,
		    updated_at = NOW()
		WHERE id = $1
	`

	_, err := r.pool.Exec(ctx, query, id, model.VerificationVerified)
	if err != nil {
		return fmt.Errorf("failed to verify email: %w", err)
	}

	return nil
}

// IncrementVerificationAttempts увеличивает счётчик попыток верификации
func (r *AccountRepo) IncrementVerificationAttempts(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE users
		SET verification_attempts = verification_attempts + 1,
		    updated_at = NOW()
		WHERE id = $1
	`

	_, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to increment verification attempts: %w", err)
	}

	return nil
}

// UpdateLastLogin обновляет время последнего входа
func (r *AccountRepo) UpdateLastLogin(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE users
		SET last_login_at = NOW(),
		    updated_at = NOW()
		WHERE id = $1
	`

	_, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to update last login: %w", err)
	}

	return nil
}

// IsEmailExists проверяет существование email
func (r *AccountRepo) IsEmailExists(ctx context.Context, email string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`

	var exists bool
	err := r.pool.QueryRow(ctx, query, email).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check email existence: %w", err)
	}

	return exists, nil
}

// IsLoginExists проверяет существование логина
func (r *AccountRepo) IsLoginExists(ctx context.Context, login string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE login = $1)`

	var exists bool
	err := r.pool.QueryRow(ctx, query, login).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check login existence: %w", err)
	}

	return exists, nil
}
