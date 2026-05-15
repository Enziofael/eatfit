package service

import (
	"context"
	"fmt"

	"github.com/Enziofael/eatfit/backend/internal/model"
	"github.com/Enziofael/eatfit/backend/internal/repository/postgres"
	"github.com/google/uuid"
)

type ProfileService struct {
	profileRepo *postgres.ProfileRepo
}

func NewProfileService(profileRepo *postgres.ProfileRepo) *ProfileService {
	return &ProfileService{profileRepo: profileRepo}
}

// GetProfile получает профиль
func (s *ProfileService) GetProfile(ctx context.Context, userID uuid.UUID) (*model.Profile, error) {
	return s.profileRepo.GetProfile(ctx, userID)
}

// CreateProfile создаёт профиль при регистрации
func (s *ProfileService) CreateProfile(ctx context.Context, userID uuid.UUID) error {
	return s.profileRepo.CreateProfile(ctx, userID)
}

// UpdateProfile обновляет профиль
func (s *ProfileService) UpdateProfile(ctx context.Context, userID uuid.UUID, req *model.UpdateProfileRequest) (*model.Profile, error) {
	// Валидация роста
	if req.Height != nil && *req.Height <= 0 {
		return nil, fmt.Errorf("height must be positive")
	}

	// Валидация пола
	if req.Gender != nil {
		valid := map[string]bool{"male": true, "female": true, "other": true}
		if !valid[*req.Gender] {
			return nil, fmt.Errorf("invalid gender")
		}
	}

	return s.profileRepo.UpdateProfile(ctx, userID, req)
}

// UpdateWeight обновляет вес
func (s *ProfileService) UpdateWeight(ctx context.Context, userID uuid.UUID, weight float64, note *string) (float64, error) {
	if weight <= 0 || weight > 500 {
		return 0, fmt.Errorf("weight must be between 0 and 500 kg")
	}

	return s.profileRepo.AddWeightEntry(ctx, userID, weight, note)
}

// UpdateNorms обновляет нормы
func (s *ProfileService) SetNorms(ctx context.Context, userID uuid.UUID, req *model.SetNormsRequest) (*model.NutritionNorms, error) {
	if req.Calories < 0 || req.Proteins < 0 || req.Fats < 0 || req.Carbs < 0 || req.Water < 0 {
		return nil, fmt.Errorf("nutrition values must be non-negative")
	}

	return s.profileRepo.UpdateNorms(ctx, userID, req)
}

// GetWeightHistory получает историю веса
func (s *ProfileService) GetWeightHistory(ctx context.Context, userID uuid.UUID, limit int) ([]model.WeightEntry, error) {
	return s.profileRepo.GetWeightHistory(ctx, userID, limit)
}