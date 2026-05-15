package model

import (
	"time"

	"github.com/google/uuid"
)

type ConsumptionRecord struct {
	ID         uuid.UUID
	UserID     uuid.UUID
	MealID     *uuid.UUID // nullable
	MealName   string
	Amount     float64
	Calories   float64
	Proteins   float64
	Fats       float64
	Carbs      float64
	Water      float64
	ConsumedAt time.Time
}

type ConsumptionGroup struct {
	Date    string
	Records []ConsumptionRecord
}

type AddConsumptionRequest struct {
	UserID    uuid.UUID
	MealID    string
	MealName  string
	Amount    float64
	Calories  float64
	Proteins  float64
	Fats      float64
	Carbs     float64
	Water     float64
}