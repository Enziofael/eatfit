// internal/service/jwt_service.go
package service

import (
	"fmt"
	"time"

	"github.com/Enziofael/eatfit/backend/internal/model"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// JWTService - сервис для работы с JWT токенами
type JWTService struct {
	secret          []byte
	accessDuration  time.Duration
	refreshDuration time.Duration
}

// NewJWTService создаёт новый экземпляр JWTService
func NewJWTService(secret string, accessDuration, refreshDuration time.Duration) *JWTService {
	return &JWTService{
		secret:          []byte(secret),
		accessDuration:  accessDuration,
		refreshDuration: refreshDuration,
	}
}

// Claims - JWT claims (данные внутри токена)
type Claims struct {
	AccountID string `json:"account_id"`
	Email     string `json:"email"`
	Login     string `json:"login"`
	jwt.RegisteredClaims
}

// GenerateAccessToken создаёт access токен
func (s *JWTService) GenerateAccessToken(account *model.Account) (string, time.Time, error) {
	expiresAt := time.Now().Add(s.accessDuration)

	claims := &Claims{
		AccountID: account.ID.String(),
		Email:     account.Email,
		Login:     account.Login,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "eatfit",
			Subject:   account.ID.String(),
			ID:        uuid.New().String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(s.secret)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to sign access token: %w", err)
	}

	return tokenString, expiresAt, nil
}

// GenerateRefreshToken создаёт refresh токен (opaque, не JWT)
func (s *JWTService) GenerateRefreshToken() (string, time.Time, error) {
	token := uuid.New().String() + "-" + uuid.New().String()
	expiresAt := time.Now().Add(s.refreshDuration)
	return token, expiresAt, nil
}

// ValidateAccessToken проверяет и парсит access токен
func (s *JWTService) ValidateAccessToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Проверяем метод подписи
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.secret, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}
