package config

import (
	"os"
	"strconv"
)

type Config struct {
	// Server
	GRPCPort string
	HTTPPort string

	// Database
	DBHost     string
	DBPort     int
	DBUser     string
	DBPassword string
	DBName     string

	// Redis
	RedisHost     string
	RedisPort     int
	RedisPassword string

	// JWT
	JWTSecret            string
	AccessTokenDuration  int // minutes
	RefreshTokenDuration int // days

	// SMTP
	SMTPHost     string
	SMTPPort     int
	SMTPUsername string
	SMTPPassword string
	SMTPSender   string
}

func Load() *Config {
	return &Config{
		GRPCPort: getEnv("BACKEND_GRPC_PORT", "50051"),
		HTTPPort: getEnv("BACKEND_HTTP_PORT", "8080"),

		DBHost:     getEnv("POSTGRES_HOST", "localhost"),
		DBPort:     getEnvInt("POSTGRES_PORT", 5432),
		DBUser:     getEnv("POSTGRES_USER", "eatfit_user"),
		DBPassword: getEnv("POSTGRES_PASSWORD", "eatfit_pass"),
		DBName:     getEnv("POSTGRES_DB", "eatfit"),

		RedisHost:     getEnv("REDIS_HOST", "localhost"),
		RedisPort:     getEnvInt("REDIS_PORT", 6379),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),

		JWTSecret:            getEnv("JWT_SECRET", "change_me"),
		AccessTokenDuration:  getEnvInt("ACCESS_TOKEN_DURATION", 15),
		RefreshTokenDuration: getEnvInt("REFRESH_TOKEN_DURATION", 30),

		SMTPHost:     getEnv("SMTP_HOST", ""),
		SMTPPort:     getEnvInt("SMTP_PORT", 587),
		SMTPUsername: getEnv("SMTP_USERNAME", ""),
		SMTPPassword: getEnv("SMTP_PASSWORD", ""),
		SMTPSender:   getEnv("SMTP_SENDER", ""),
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return fallback
}
