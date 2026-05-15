package model

import (
	"time"

	"github.com/google/uuid"
)

// Profile - профиль пользователя
type Profile struct {
	UserID        uuid.UUID
	Login         string
	Email         string
	FirstName     *string
	LastName      *string
	AvatarURL     *string
	Height        *float64
	BirthDate     *time.Time
	Gender        *string
	Bio           *string
	CurrentWeight *float64
	Age           *int
	Norms         *NutritionNorms
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// DisplayName возвращает имя для отображения
func (p *Profile) DisplayName(login string) string {
	if p.FirstName != nil && *p.FirstName != "" {
		if p.LastName != nil && *p.LastName != "" {
			return *p.FirstName + " " + *p.LastName
		}
		return *p.FirstName
	}
	if p.LastName != nil && *p.LastName != "" {
		return *p.LastName
	}
	return login
}

// NutritionNorms - нормы КБЖУ и воды
type NutritionNorms struct {
	Calories  float64
	Proteins  float64
	Fats      float64
	Carbs     float64
	Water     float64
	UpdatedAt time.Time
}

// WeightEntry - запись об изменении веса
type WeightEntry struct {
	ID         uuid.UUID
	UserID     uuid.UUID
	Weight     float64
	Note       *string
	RecordedAt time.Time
}

// UpdateProfileRequest - запрос на обновление профиля
type UpdateProfileRequest struct {
	FirstName *string
	LastName  *string
	AvatarURL *string
	Height    *float64
	BirthDate *string // YYYY-MM-DD
	Gender    *string
	Bio       *string
}

// UpdateNormsRequest - запрос на обновление норм
type SetNormsRequest struct {
	Calories float64
	Proteins float64
	Fats     float64
	Carbs    float64
	Water    float64
}
