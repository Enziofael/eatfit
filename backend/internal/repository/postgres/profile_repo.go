package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/Enziofael/eatfit/backend/internal/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProfileRepo struct {
	pool *pgxpool.Pool
}

func NewProfileRepo(pool *pgxpool.Pool) *ProfileRepo {
	return &ProfileRepo{pool: pool}
}

// GetProfile получает профиль пользователя с актуальными нормами и весом
func (r *ProfileRepo) GetProfile(ctx context.Context, userID uuid.UUID) (*model.Profile, error) {
	query := `
		SELECT 
			p.user_id, u.login, u.email,
			p.first_name, p.last_name, p.avatar_url,
			p.height, p.birth_date, p.gender, p.bio,
			p.created_at, p.updated_at,
			(SELECT weight FROM weight_history WHERE user_id = $1 ORDER BY recorded_at DESC LIMIT 1) as current_weight,
			(SELECT calculate_age(p.birth_date)) as age,
			nn.calories, nn.proteins, nn.fats, nn.carbs, nn.water, nn.created_at as norms_updated_at
		FROM profiles p
		JOIN users u ON u.id = p.user_id
		LEFT JOIN LATERAL (
			SELECT * FROM nutrition_norms 
			WHERE user_id = $1 
			ORDER BY created_at DESC 
			LIMIT 1
		) nn ON true
		WHERE p.user_id = $1
	`

	profile := &model.Profile{}
	var norms model.NutritionNorms
	var normsUpdatedAt *time.Time
	var login, email string

	err := r.pool.QueryRow(ctx, query, userID).Scan(
		&profile.UserID,
		&login,
		&email,
		&profile.FirstName,
		&profile.LastName,
		&profile.AvatarURL,
		&profile.Height,
		&profile.BirthDate,
		&profile.Gender,
		&profile.Bio,
		&profile.CreatedAt,
		&profile.UpdatedAt,
		&profile.CurrentWeight,
		&profile.Age,
		&norms.Calories,
		&norms.Proteins,
		&norms.Fats,
		&norms.Carbs,
		&norms.Water,
		&normsUpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get profile: %w", err)
	}

	profile.Login = login
	profile.Email = email

	if normsUpdatedAt != nil {
		norms.UpdatedAt = *normsUpdatedAt
		profile.Norms = &norms
	}

	return profile, nil
}

// CreateProfile создаёт профиль (вызывается при регистрации)
func (r *ProfileRepo) CreateProfile(ctx context.Context, userID uuid.UUID) error {
	query := `INSERT INTO profiles (user_id) VALUES ($1)`
	_, err := r.pool.Exec(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("failed to create profile: %w", err)
	}
	return nil
}

// UpdateProfile обновляет профиль
func (r *ProfileRepo) UpdateProfile(ctx context.Context, userID uuid.UUID, req *model.UpdateProfileRequest) (*model.Profile, error) {
	query := `
		UPDATE profiles SET
			first_name = COALESCE($2, first_name),
			last_name = COALESCE($3, last_name),
			avatar_url = COALESCE($4, avatar_url),
			height = COALESCE($5, height),
			birth_date = COALESCE($6, birth_date),
			gender = COALESCE($7, gender),
			bio = COALESCE($8, bio),
			updated_at = NOW()
		WHERE user_id = $1
	`

	var birthDate *time.Time
	if req.BirthDate != nil {
		parsed, err := time.Parse("2006-01-02", *req.BirthDate)
		if err == nil {
			birthDate = &parsed
		}
	}

	_, err := r.pool.Exec(ctx, query,
		userID,
		req.FirstName,
		req.LastName,
		req.AvatarURL,
		req.Height,
		birthDate,
		req.Gender,
		req.Bio,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to update profile: %w", err)
	}

	return r.GetProfile(ctx, userID)
}

// AddWeightEntry добавляет запись об изменении веса
func (r *ProfileRepo) AddWeightEntry(ctx context.Context, userID uuid.UUID, weight float64, note *string) (float64, error) {
	query := `
		INSERT INTO weight_history (user_id, weight, note)
		VALUES ($1, $2, $3)
		RETURNING weight
	`

	var newWeight float64
	err := r.pool.QueryRow(ctx, query, userID, weight, note).Scan(&newWeight)
	if err != nil {
		return 0, fmt.Errorf("failed to add weight entry: %w", err)
	}

	return newWeight, nil
}

// UpdateNorms обновляет нормы КБЖУ
func (r *ProfileRepo) UpdateNorms(ctx context.Context, userID uuid.UUID, req *model.SetNormsRequest) (*model.NutritionNorms, error) {
	query := `
		INSERT INTO nutrition_norms (user_id, calories, proteins, fats, carbs, water)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING calories, proteins, fats, carbs, water, created_at
	`

	norms := &model.NutritionNorms{}
	err := r.pool.QueryRow(ctx, query,
		userID,
		req.Calories,
		req.Proteins,
		req.Fats,
		req.Carbs,
		req.Water,
	).Scan(
		&norms.Calories,
		&norms.Proteins,
		&norms.Fats,
		&norms.Carbs,
		&norms.Water,
		&norms.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to update norms: %w", err)
	}

	return norms, nil
}

// GetWeightHistory получает историю изменения веса
func (r *ProfileRepo) GetWeightHistory(ctx context.Context, userID uuid.UUID, limit int) ([]model.WeightEntry, error) {
	if limit <= 0 {
		limit = 30
	}

	query := `
		SELECT id, user_id, weight, note, recorded_at
		FROM weight_history
		WHERE user_id = $1
		ORDER BY recorded_at DESC
		LIMIT $2
	`

	rows, err := r.pool.Query(ctx, query, userID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get weight history: %w", err)
	}
	defer rows.Close()

	var entries []model.WeightEntry
	for rows.Next() {
		var entry model.WeightEntry
		err := rows.Scan(
			&entry.ID,
			&entry.UserID,
			&entry.Weight,
			&entry.Note,
			&entry.RecordedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan weight entry: %w", err)
		}
		entries = append(entries, entry)
	}

	return entries, nil
}
