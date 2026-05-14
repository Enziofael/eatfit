// internal/config/config.go
package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Config содержит всю конфигурацию приложения
type Config struct {
	// Сервер
	GRPCPort string
	HTTPPort string

	// Базы данных
	PostgresDSN string
	RedisAddr   string
	RedisPass   string
	RedisDB     int

	// JWT
	JWTSecret            string
	AccessTokenDuration  time.Duration
	RefreshTokenDuration time.Duration

	// Email (SMTP)
	SMTPHost      string
	SMTPPort      int
	SMTPUser      string
	SMTPPass      string
	EmailFrom     string
	EmailFromName string

	// Приложение
	AppEnv   string // "development", "production"
	AppDebug bool
}

// Load загружает конфигурацию из .env файла и переменных окружения
func Load() (*Config, error) {
	// Загружаем .env из корня проекта, если существует
	// Игнорируем ошибку - .env опционален в production
	_ = godotenv.Load("../../.env")
	_ = godotenv.Load(".env")

	cfg := &Config{
		// Сервер
		GRPCPort: getEnv("BACKEND_GRPC_PORT", "50051"),
		HTTPPort: getEnv("BACKEND_HTTP_PORT", "8080"),

		// PostgreSQL
		PostgresDSN: buildPostgresDSN(),

		// Redis
		RedisAddr: fmt.Sprintf("%s:%s",
			getEnv("REDIS_HOST", "localhost"),
			getEnv("REDIS_PORT", "6379"),
		),
		RedisPass: getEnv("REDIS_PASSWORD", ""),
		RedisDB:   0,

		// JWT
		JWTSecret:            getEnv("JWT_SECRET", "default-secret-change-in-production"),
		AccessTokenDuration:  15 * time.Minute,
		RefreshTokenDuration: 30 * 24 * time.Hour, // 30 дней

		// Email
		SMTPHost:      getEnv("SMTP_HOST", "localhost"),
		SMTPPort:      getEnvAsInt("SMTP_PORT", 1025),
		SMTPUser:      getEnv("SMTP_USER", ""),
		SMTPPass:      getEnv("SMTP_PASS", ""),
		EmailFrom:     getEnv("EMAIL_FROM", "noreply@eatfit.local"),
		EmailFromName: "Eatfit",

		// Приложение
		AppEnv:   getEnv("APP_ENV", "development"),
		AppDebug: getEnv("APP_DEBUG", "true") == "true",
	}

	// Валидация обязательных полей
	if cfg.AppEnv == "production" {
		if cfg.JWTSecret == "default-secret-change-in-production" {
			return nil, fmt.Errorf("JWT_SECRET must be set in production")
		}
		if cfg.PostgresDSN == "" {
			return nil, fmt.Errorf("database configuration is required")
		}
	}

	return cfg, nil
}

// getEnv возвращает значение переменной окружения или значение по умолчанию
func getEnv(key, defaultVal string) string {
	if val, exists := os.LookupEnv(key); exists && val != "" {
		return val
	}
	return defaultVal
}

// buildPostgresDSN формирует строку подключения к PostgreSQL
func buildPostgresDSN() string {
	host := getEnv("POSTGRES_HOST", "localhost")
	port := getEnv("POSTGRES_PORT", "5432")
	user := getEnv("POSTGRES_USER", "eatfit_user")
	pass := getEnv("POSTGRES_PASSWORD", "eatfit_pass")
	dbname := getEnv("POSTGRES_DB", "eatfit")

	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		user, pass, host, port, dbname,
	)
}

func getEnvAsInt(key string, defaultVal int) int {
	if val, exists := os.LookupEnv(key); exists && val != "" {
		if intVal, err := strconv.Atoi(val); err == nil {
			return intVal
		}
	}
	return defaultVal
}
