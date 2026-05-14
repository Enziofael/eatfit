// pkg/password/hasher.go
package password

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// Hasher - сервис хеширования паролей
type Hasher struct {
	cost int
}

// NewHasher создаёт новый экземпляр Hasher
// cost - стоимость bcrypt (рекомендуется 12 для production)
func NewHasher(cost int) *Hasher {
	if cost < bcrypt.MinCost || cost > bcrypt.MaxCost {
		cost = bcrypt.DefaultCost
	}
	return &Hasher{cost: cost}
}

// HashPassword хеширует пароль с использованием bcrypt
func (h *Hasher) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), h.cost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(bytes), nil
}

// VerifyPassword проверяет соответствие пароля хешу
func (h *Hasher) VerifyPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
