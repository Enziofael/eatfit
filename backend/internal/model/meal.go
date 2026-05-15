package model

import (
	"time"

	"github.com/google/uuid"
)

type Meal struct {
	ID          uuid.UUID
	UserID      uuid.UUID
	Name        string
	Description *string
	Recipe      *string
	ImageURL    *string
	Calories    float64
	Proteins    float64
	Fats        float64
	Carbs       float64
	Water       float64
	Components  []MealComponent
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type MealComponent struct {
	ComponentMealID uuid.UUID
	ComponentName   string
	Amount          float64
	Calories        float64
	Proteins        float64
	Fats            float64
	Carbs           float64
	Water           float64
	SortOrder       int
}

type CreateMealRequest struct {
	UserID      uuid.UUID
	Name        string
	Description *string
	Recipe      *string
	ImageURL    *string
	Calories    float64
	Proteins    float64
	Fats        float64
	Carbs       float64
	Water       float64
	Components  []MealComponentInput
}

type MealComponentInput struct {
	ComponentMealID uuid.UUID
	Amount          float64
}

type UpdateMealRequest struct {
	ID          uuid.UUID
	Name        *string
	Description *string
	Recipe      *string
	ImageURL    *string
	Calories    *float64
	Proteins    *float64
	Fats        *float64
	Carbs       *float64
	Water       *float64
	Components  []MealComponentInput
}

type ListMealsRequest struct {
	UserID    uuid.UUID
	SortBy    string
	SortOrder string
	Page      int
	PageSize  int
}

type SearchMealsRequest struct {
	UserID    uuid.UUID
	Query     string
	SortBy    string
	SortOrder string
	Page      int
	PageSize  int
}
