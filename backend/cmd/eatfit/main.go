// cmd/eatfit/main.go
package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Enziofael/eatfit/backend/internal/config"
	grpcHandler "github.com/Enziofael/eatfit/backend/internal/handler/grpc"
	"github.com/Enziofael/eatfit/backend/internal/repository/postgres"
	redisRepo "github.com/Enziofael/eatfit/backend/internal/repository/redis"
	"github.com/Enziofael/eatfit/backend/internal/service"
	"github.com/Enziofael/eatfit/backend/pkg/password"

	pb "github.com/Enziofael/eatfit/backend/api/gen/eatfit/v1"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	// Загружаем конфигурацию
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	log.Printf("Starting Eatfit Backend Server")
	log.Printf("Environment: %s", cfg.AppEnv)
	log.Printf("gRPC port: %s", cfg.GRPCPort)

	// Инициализируем подключение к PostgreSQL
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pgPool, err := pgxpool.New(ctx, cfg.PostgresDSN)
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}
	defer pgPool.Close()

	if err := pgPool.Ping(ctx); err != nil {
		log.Fatalf("Failed to ping PostgreSQL: %v", err)
	}
	log.Println("Connected to PostgreSQL")

	// Инициализируем подключение к Redis
	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPass,
		DB:       cfg.RedisDB,
	})

	if err := redisClient.Ping(ctx).Err(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	log.Println("Connected to Redis")

	// ============================================
	// Инициализируем репозитории
	// ============================================
	accountRepo := postgres.NewAccountRepo(pgPool)
	tokenRepo := redisRepo.NewTokenRepo(redisClient)
	profileRepo := postgres.NewProfileRepo(pgPool)

	// ============================================
	// Инициализируем сервисы
	// ============================================
	jwtService := service.NewJWTService(
		cfg.JWTSecret,
		cfg.AccessTokenDuration,
		cfg.RefreshTokenDuration,
	)
	emailService := service.NewEmailService(cfg)
	passwordHasher := password.NewHasher(12)
	profileService := service.NewProfileService(profileRepo)

	authService := service.NewAuthService(
		accountRepo,
		tokenRepo,
		jwtService,
		emailService,
		passwordHasher,
		profileService, // ← передаём profileService
	)

	// ============================================
	// Создаём gRPC сервер
	// ============================================
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			loggingInterceptor(),
			recoveryInterceptor(),
		),
	)

	// ============================================
	// Регистрируем сервисы
	// ============================================
	authHandler := grpcHandler.NewAuthHandler(authService)
	pb.RegisterAuthServiceServer(grpcServer, authHandler)

	profileHandler := grpcHandler.NewProfileHandler(profileService)
	pb.RegisterProfileServiceServer(grpcServer, profileHandler) // ← теперь grpcServer существует

	// ============================================
	// Reflection (только для разработки)
	// ============================================
	if cfg.AppEnv == "development" {
		reflection.Register(grpcServer)
		log.Println("gRPC reflection enabled")
	}

	// ============================================
	// Запускаем сервер
	// ============================================
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.GRPCPort))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		sig := <-sigChan
		log.Printf("Received signal: %v, shutting down gracefully...", sig)

		grpcServer.GracefulStop()
		log.Println("gRPC server stopped")

		redisClient.Close()
		log.Println("Redis connection closed")

		pgPool.Close()
		log.Println("PostgreSQL connection closed")

		os.Exit(0)
	}()

	log.Printf("gRPC server listening on :%s", cfg.GRPCPort)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

// loggingInterceptor логирует все входящие gRPC запросы
func loggingInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		start := time.Now()
		log.Printf("gRPC Request: %s", info.FullMethod)

		resp, err := handler(ctx, req)

		duration := time.Since(start)
		if err != nil {
			log.Printf("gRPC Error: %s - %v (took %s)", info.FullMethod, err, duration)
		} else {
			log.Printf("gRPC Success: %s (took %s)", info.FullMethod, duration)
		}

		return resp, err
	}
}

// recoveryInterceptor перехватывает паники и возвращает Internal ошибку
func recoveryInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("PANIC recovered: %v", r)
				err = fmt.Errorf("internal server error")
			}
		}()
		return handler(ctx, req)
	}
}
