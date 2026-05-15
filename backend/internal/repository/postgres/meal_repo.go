package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/Enziofael/eatfit/backend/internal/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type MealRepo struct {
	pool *pgxpool.Pool
}

func NewMealRepo(pool *pgxpool.Pool) *MealRepo {
	return &MealRepo{pool: pool}
}

func (r *MealRepo) CreateMeal(ctx context.Context, req *model.CreateMealRequest) (*model.Meal, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	query := `
		INSERT INTO meals (user_id, name, description, recipe, image_url, calories, proteins, fats, carbs, water)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id, created_at, updated_at
	`

	meal := &model.Meal{
		UserID:      req.UserID,
		Name:        req.Name,
		Description: req.Description,
		Recipe:      req.Recipe,
		ImageURL:    req.ImageURL,
		Calories:    req.Calories,
		Proteins:    req.Proteins,
		Fats:        req.Fats,
		Carbs:       req.Carbs,
		Water:       req.Water,
	}

	err = tx.QueryRow(ctx, query,
		meal.UserID, meal.Name, meal.Description, meal.Recipe, meal.ImageURL,
		meal.Calories, meal.Proteins, meal.Fats, meal.Carbs, meal.Water,
	).Scan(&meal.ID, &meal.CreatedAt, &meal.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create meal: %w", err)
	}

	// Сохраняем компоненты
	for i, c := range req.Components {
		if c.Amount <= 0 {
			c.Amount = 100
		}
		_, err = tx.Exec(ctx, `
			INSERT INTO meal_components (meal_id, component_meal_id, amount, sort_order)
			VALUES ($1, $2, $3, $4)
		`, meal.ID, c.ComponentMealID, c.Amount, i)
		if err != nil {
			return nil, fmt.Errorf("failed to insert component: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit: %w", err)
	}

	return meal, nil
}

func (r *MealRepo) GetMeal(ctx context.Context, mealID uuid.UUID) (*model.Meal, error) {
	query := `
		SELECT id, user_id, name, description, recipe, image_url, 
		       calories, proteins, fats, carbs, water, created_at, updated_at
		FROM meals WHERE id = $1
	`

	meal := &model.Meal{}
	err := r.pool.QueryRow(ctx, query, mealID).Scan(
		&meal.ID, &meal.UserID, &meal.Name, &meal.Description, &meal.Recipe, &meal.ImageURL,
		&meal.Calories, &meal.Proteins, &meal.Fats, &meal.Carbs, &meal.Water,
		&meal.CreatedAt, &meal.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get meal: %w", err)
	}

	components, err := r.getComponents(ctx, mealID)
	if err != nil {
		return nil, err
	}
	meal.Components = components

	return meal, nil
}

func (r *MealRepo) UpdateMeal(ctx context.Context, req *model.UpdateMealRequest) (*model.Meal, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	query := `
		UPDATE meals SET
			name = COALESCE($2, name),
			description = COALESCE($3, description),
			recipe = COALESCE($4, recipe),
			image_url = COALESCE($5, image_url),
			calories = COALESCE($6, calories),
			proteins = COALESCE($7, proteins),
			fats = COALESCE($8, fats),
			carbs = COALESCE($9, carbs),
			water = COALESCE($10, water),
			updated_at = NOW()
		WHERE id = $1
		RETURNING user_id, name, description, recipe, image_url,
		          calories, proteins, fats, carbs, water, created_at, updated_at
	`

	meal := &model.Meal{ID: req.ID}
	err = tx.QueryRow(ctx, query,
		req.ID, req.Name, req.Description, req.Recipe, req.ImageURL,
		req.Calories, req.Proteins, req.Fats, req.Carbs, req.Water,
	).Scan(
		&meal.UserID, &meal.Name, &meal.Description, &meal.Recipe, &meal.ImageURL,
		&meal.Calories, &meal.Proteins, &meal.Fats, &meal.Carbs, &meal.Water,
		&meal.CreatedAt, &meal.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update meal: %w", err)
	}

	// Обновляем компоненты
	if req.Components != nil {
		_, err = tx.Exec(ctx, "DELETE FROM meal_components WHERE meal_id = $1", req.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to delete old components: %w", err)
		}
		for i, c := range req.Components {
			if c.Amount <= 0 {
				c.Amount = 100
			}
			_, err = tx.Exec(ctx, `
				INSERT INTO meal_components (meal_id, component_meal_id, amount, sort_order)
				VALUES ($1, $2, $3, $4)
			`, req.ID, c.ComponentMealID, c.Amount, i)
			if err != nil {
				return nil, fmt.Errorf("failed to insert component: %w", err)
			}
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit: %w", err)
	}

	meal.Components, _ = r.getComponents(ctx, req.ID)
	return meal, nil
}

func (r *MealRepo) DeleteMeal(ctx context.Context, mealID uuid.UUID) error {
	_, err := r.pool.Exec(ctx, "DELETE FROM meals WHERE id = $1", mealID)
	return err
}

func (r *MealRepo) ListMeals(ctx context.Context, req *model.ListMealsRequest) ([]model.Meal, int, error) {
	validSort := map[string]bool{
		"name": true, "calories": true, "proteins": true,
		"fats": true, "carbs": true, "water": true, "created_at": true,
	}
	sortBy := "created_at"
	if validSort[req.SortBy] {
		sortBy = req.SortBy
	}
	sortOrder := "DESC"
	if strings.ToLower(req.SortOrder) == "asc" {
		sortOrder = "ASC"
	}
	if req.PageSize <= 0 {
		req.PageSize = 20
	}
	if req.Page <= 0 {
		req.Page = 1
	}
	offset := (req.Page - 1) * req.PageSize

	var total int
	err := r.pool.QueryRow(ctx, "SELECT COUNT(*) FROM meals WHERE user_id = $1", req.UserID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	query := fmt.Sprintf(`
		SELECT id, user_id, name, COALESCE(description, ''), COALESCE(recipe, ''), COALESCE(image_url, ''),
		       calories, proteins, fats, carbs, water, created_at, updated_at
		FROM meals WHERE user_id = $1
		ORDER BY %s %s LIMIT $2 OFFSET $3
	`, sortBy, sortOrder)

	rows, err := r.pool.Query(ctx, query, req.UserID, req.PageSize, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var meals []model.Meal
	for rows.Next() {
		var m model.Meal
		rows.Scan(&m.ID, &m.UserID, &m.Name, &m.Description, &m.Recipe, &m.ImageURL,
			&m.Calories, &m.Proteins, &m.Fats, &m.Carbs, &m.Water, &m.CreatedAt, &m.UpdatedAt)
		meals = append(meals, m)
	}
	return meals, total, nil
}

func (r *MealRepo) SearchMeals(ctx context.Context, req *model.SearchMealsRequest) ([]model.Meal, int, error) {
	validSort := map[string]bool{
		"name": true, "calories": true, "proteins": true,
		"fats": true, "carbs": true, "water": true, "created_at": true,
	}
	sortBy := "created_at"
	if validSort[req.SortBy] {
		sortBy = req.SortBy
	}
	sortOrder := "DESC"
	if strings.ToLower(req.SortOrder) == "asc" {
		sortOrder = "ASC"
	}
	if req.PageSize <= 0 {
		req.PageSize = 20
	}
	if req.Page <= 0 {
		req.Page = 1
	}
	offset := (req.Page - 1) * req.PageSize
	pattern := "%" + req.Query + "%"

	var total int
	r.pool.QueryRow(ctx, "SELECT COUNT(*) FROM meals WHERE user_id = $1 AND name ILIKE $2", req.UserID, pattern).Scan(&total)

	query := fmt.Sprintf(`
		SELECT id, user_id, name, COALESCE(description, ''), COALESCE(recipe, ''), COALESCE(image_url, ''),
		       calories, proteins, fats, carbs, water, created_at, updated_at
		FROM meals WHERE user_id = $1 AND name ILIKE $2
		ORDER BY %s %s LIMIT $3 OFFSET $4
	`, sortBy, sortOrder)

	rows, _ := r.pool.Query(ctx, query, req.UserID, pattern, req.PageSize, offset)
	defer rows.Close()

	var meals []model.Meal
	for rows.Next() {
		var m model.Meal
		rows.Scan(&m.ID, &m.UserID, &m.Name, &m.Description, &m.Recipe, &m.ImageURL,
			&m.Calories, &m.Proteins, &m.Fats, &m.Carbs, &m.Water, &m.CreatedAt, &m.UpdatedAt)
		meals = append(meals, m)
	}
	return meals, total, nil
}

func (r *MealRepo) getComponents(ctx context.Context, mealID uuid.UUID) ([]model.MealComponent, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT mc.component_meal_id, m.name, mc.amount,
		       m.calories, m.proteins, m.fats, m.carbs, m.water, mc.sort_order
		FROM meal_components mc
		JOIN meals m ON m.id = mc.component_meal_id
		WHERE mc.meal_id = $1
		ORDER BY mc.sort_order
	`, mealID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var components []model.MealComponent
	for rows.Next() {
		var c model.MealComponent
		rows.Scan(&c.ComponentMealID, &c.ComponentName, &c.Amount,
			&c.Calories, &c.Proteins, &c.Fats, &c.Carbs, &c.Water, &c.SortOrder)
		components = append(components, c)
	}
	return components, nil
}
