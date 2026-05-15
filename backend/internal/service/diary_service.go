package service

import (
	"context"

	"github.com/Enziofael/eatfit/backend/internal/model"
	"github.com/Enziofael/eatfit/backend/internal/repository/postgres"
	"github.com/google/uuid"
)

type DiaryService struct {
	consumptionRepo *postgres.ConsumptionRepo
}

func NewDiaryService(consumptionRepo *postgres.ConsumptionRepo) *DiaryService {
	return &DiaryService{consumptionRepo: consumptionRepo}
}

func (s *DiaryService) AddConsumption(ctx context.Context, req *model.AddConsumptionRequest) (*model.ConsumptionRecord, error) {
	return s.consumptionRepo.AddConsumption(ctx, req)
}

func (s *DiaryService) ListConsumptions(ctx context.Context, userID uuid.UUID, days int) ([]model.ConsumptionGroup, error) {
	return s.consumptionRepo.ListConsumptions(ctx, userID, days)
}

func (s *DiaryService) DeleteConsumption(ctx context.Context, recordID uuid.UUID) error {
	return s.consumptionRepo.DeleteConsumption(ctx, recordID)
}
