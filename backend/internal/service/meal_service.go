package service

import (
	"context"
	"fmt"

	"github.com/Enziofael/eatfit/backend/internal/model"
	"github.com/Enziofael/eatfit/backend/internal/repository/postgres"
	"github.com/google/uuid"
)

type MealService struct {
	mealRepo *postgres.MealRepo
}

func NewMealService(mealRepo *postgres.MealRepo) *MealService {
	return &MealService{mealRepo: mealRepo}
}

func (s *MealService) CreateMeal(ctx context.Context, req *model.CreateMealRequest) (*model.Meal, error) {
	if req.Name == "" {
		return nil, fmt.Errorf("name is required")
	}

	// Если КБЖУ не указаны и есть компоненты — вычисляем из компонентов
	if req.Calories == 0 && req.Proteins == 0 && req.Fats == 0 && req.Carbs == 0 && req.Water == 0 && len(req.Components) > 0 {
		var calories, proteins, fats, carbs, water float64
		for _, c := range req.Components {
			component, err := s.mealRepo.GetMeal(ctx, c.ComponentMealID)
			if err != nil {
				return nil, fmt.Errorf("component meal not found: %s", c.ComponentMealID)
			}
			ratio := c.Amount / 100.0
			calories += component.Calories * ratio
			proteins += component.Proteins * ratio
			fats += component.Fats * ratio
			carbs += component.Carbs * ratio
			water += component.Water * ratio
		}
		req.Calories = calories
		req.Proteins = proteins
		req.Fats = fats
		req.Carbs = carbs
		req.Water = water
	}

	return s.mealRepo.CreateMeal(ctx, req)
}

func (s *MealService) GetMeal(ctx context.Context, mealID uuid.UUID) (*model.Meal, error) {
	return s.mealRepo.GetMeal(ctx, mealID)
}

func (s *MealService) UpdateMeal(ctx context.Context, req *model.UpdateMealRequest) (*model.Meal, error) {
	// Если переданы компоненты и КБЖУ не указаны явно — пересчитываем
	if req.Components != nil && req.Calories == nil {
		var calories, proteins, fats, carbs, water float64
		for _, c := range req.Components {
			component, err := s.mealRepo.GetMeal(ctx, c.ComponentMealID)
			if err != nil {
				return nil, fmt.Errorf("component meal not found: %s", c.ComponentMealID)
			}
			ratio := c.Amount / 100.0
			calories += component.Calories * ratio
			proteins += component.Proteins * ratio
			fats += component.Fats * ratio
			carbs += component.Carbs * ratio
			water += component.Water * ratio
		}
		req.Calories = &calories
		req.Proteins = &proteins
		req.Fats = &fats
		req.Carbs = &carbs
		req.Water = &water
	}

	return s.mealRepo.UpdateMeal(ctx, req)
}

func (s *MealService) DeleteMeal(ctx context.Context, mealID uuid.UUID) error {
	return s.mealRepo.DeleteMeal(ctx, mealID)
}

func (s *MealService) ListMeals(ctx context.Context, req *model.ListMealsRequest) ([]model.Meal, int, error) {
	return s.mealRepo.ListMeals(ctx, req)
}

func (s *MealService) SearchMeals(ctx context.Context, req *model.SearchMealsRequest) ([]model.Meal, int, error) {
	return s.mealRepo.SearchMeals(ctx, req)
}
