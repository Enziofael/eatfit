package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/Enziofael/eatfit/backend/internal/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ConsumptionRepo struct {
	pool *pgxpool.Pool
}

func NewConsumptionRepo(pool *pgxpool.Pool) *ConsumptionRepo {
	return &ConsumptionRepo{pool: pool}
}

func (r *ConsumptionRepo) AddConsumption(ctx context.Context, req *model.AddConsumptionRequest) (*model.ConsumptionRecord, error) {
	var mealID *uuid.UUID
	if req.MealID != "" {
		id, err := uuid.Parse(req.MealID)
		if err == nil {
			mealID = &id
		}
	}

	query := `
		INSERT INTO consumption_records (user_id, meal_id, meal_name, amount, calories, proteins, fats, carbs, water)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, consumed_at
	`

	record := &model.ConsumptionRecord{
		UserID:   req.UserID,
		MealID:   mealID,
		MealName: req.MealName,
		Amount:   req.Amount,
		Calories: req.Calories,
		Proteins: req.Proteins,
		Fats:     req.Fats,
		Carbs:    req.Carbs,
		Water:    req.Water,
	}

	err := r.pool.QueryRow(ctx, query,
		record.UserID, mealID, record.MealName, record.Amount,
		record.Calories, record.Proteins, record.Fats, record.Carbs, record.Water,
	).Scan(&record.ID, &record.ConsumedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to add consumption: %w", err)
	}

	return record, nil
}

func (r *ConsumptionRepo) ListConsumptions(ctx context.Context, userID uuid.UUID, days int) ([]model.ConsumptionGroup, error) {
	if days <= 0 {
		days = 30
	}
	since := time.Now().AddDate(0, 0, -days)

	query := `
		SELECT id, user_id, meal_id, meal_name, amount, calories, proteins, fats, carbs, water, consumed_at
		FROM consumption_records
		WHERE user_id = $1 AND consumed_at >= $2
		ORDER BY consumed_at DESC
	`

	rows, err := r.pool.Query(ctx, query, userID, since)
	if err != nil {
		return nil, fmt.Errorf("failed to list consumptions: %w", err)
	}
	defer rows.Close()

	groups := make(map[string][]model.ConsumptionRecord)
	var orderedDates []string

	for rows.Next() {
		var r model.ConsumptionRecord
		err := rows.Scan(&r.ID, &r.UserID, &r.MealID, &r.MealName, &r.Amount,
			&r.Calories, &r.Proteins, &r.Fats, &r.Carbs, &r.Water, &r.ConsumedAt)
		if err != nil {
			continue
		}
		date := r.ConsumedAt.Format("2006-01-02")
		if _, exists := groups[date]; !exists {
			orderedDates = append(orderedDates, date)
		}
		groups[date] = append(groups[date], r)
	}

	var result []model.ConsumptionGroup
	for _, date := range orderedDates {
		result = append(result, model.ConsumptionGroup{
			Date:    date,
			Records: groups[date],
		})
	}

	return result, nil
}

func (r *ConsumptionRepo) DeleteConsumption(ctx context.Context, recordID uuid.UUID) error {
	_, err := r.pool.Exec(ctx, "DELETE FROM consumption_records WHERE id = $1", recordID)
	return err
}
