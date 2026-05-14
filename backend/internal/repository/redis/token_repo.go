// internal/repository/redis/token_repo.go
package redis

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

// TokenRepo - репозиторий для работы с токенами в Redis
type TokenRepo struct {
	client *redis.Client
}

// NewTokenRepo создаёт новый экземпляр TokenRepo
func NewTokenRepo(client *redis.Client) *TokenRepo {
	return &TokenRepo{client: client}
}

// SaveRefreshToken сохраняет refresh токен в Redis
func (r *TokenRepo) SaveRefreshToken(ctx context.Context, accountID uuid.UUID, tokenHash string, expiresAt time.Time) error {
	key := fmt.Sprintf("refresh_token:%s", tokenHash)
	ttl := time.Until(expiresAt)

	err := r.client.Set(ctx, key, accountID.String(), ttl).Err()
	if err != nil {
		return fmt.Errorf("failed to save refresh token: %w", err)
	}

	return nil
}

// GetRefreshToken получает ID аккаунта по refresh токену
func (r *TokenRepo) GetRefreshToken(ctx context.Context, tokenHash string) (uuid.UUID, error) {
	key := fmt.Sprintf("refresh_token:%s", tokenHash)

	accountIDStr, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return uuid.Nil, fmt.Errorf("refresh token not found or expired")
	}
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to get refresh token: %w", err)
	}

	accountID, err := uuid.Parse(accountIDStr)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid account ID in token: %w", err)
	}

	return accountID, nil
}

// RevokeRefreshToken отзывает refresh токен
func (r *TokenRepo) RevokeRefreshToken(ctx context.Context, tokenHash string) error {
	key := fmt.Sprintf("refresh_token:%s", tokenHash)

	err := r.client.Del(ctx, key).Err()
	if err != nil {
		return fmt.Errorf("failed to revoke refresh token: %w", err)
	}

	return nil
}

// SaveVerificationCode сохраняет код верификации в Redis
func (r *TokenRepo) SaveVerificationCode(ctx context.Context, accountID uuid.UUID, code string, expiresAt time.Time) error {
	key := fmt.Sprintf("verification:%s", accountID.String())
	ttl := time.Until(expiresAt)

	err := r.client.Set(ctx, key, code, ttl).Err()
	if err != nil {
		return fmt.Errorf("failed to save verification code: %w", err)
	}

	return nil
}

// GetVerificationCode получает код верификации из Redis
func (r *TokenRepo) GetVerificationCode(ctx context.Context, accountID uuid.UUID) (string, error) {
	key := fmt.Sprintf("verification:%s", accountID.String())

	code, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", fmt.Errorf("verification code not found or expired")
	}
	if err != nil {
		return "", fmt.Errorf("failed to get verification code: %w", err)
	}

	return code, nil
}

// DeleteVerificationCode удаляет код верификации
func (r *TokenRepo) DeleteVerificationCode(ctx context.Context, accountID uuid.UUID) error {
	key := fmt.Sprintf("verification:%s", accountID.String())

	err := r.client.Del(ctx, key).Err()
	if err != nil {
		return fmt.Errorf("failed to delete verification code: %w", err)
	}

	return nil
}

// HashToken хеширует токен для хранения
func HashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}
