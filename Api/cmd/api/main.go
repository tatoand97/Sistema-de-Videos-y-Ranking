package main

import (
	"api/internal/presentation/handlers"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"api/internal/application/useCase"
	"api/internal/domain/interfaces"
	infraCache "api/internal/infrastructure/cache"
	infraMessaging "api/internal/infrastructure/messaging"
	postgresrepo "api/internal/infrastructure/repository"
	"api/internal/infrastructure/storage"
)

// runMigrations initializes and applies DB migrations, ensuring proper cleanup.
func runMigrations(dsn, migPath string) error {
	m, err := migrate.New(migPath, dsn)
	if err != nil {
		return fmt.Errorf("migrate init failed: %w", err)
	}
	defer func() {
		srcErr, dbErr := m.Close()
		if srcErr != nil || dbErr != nil {
			log.Printf("migrate close warnings: src=%v db=%v", srcErr, dbErr)
		}
	}()

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("migrate up failed: %w", err)
	}
	return nil
}

func openDB(dsn string) (*gorm.DB, error) {
	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}

func getEnvOrDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func atoiOrDefault(s string, def int) int {
	if s == "" {
		return def
	}
	if n, err := strconv.Atoi(s); err == nil {
		return n
	}
	return def
}

func loadMinioConfigFromEnv() storage.MinioConfig {
	return storage.MinioConfig{
		Endpoint:  os.Getenv("MINIO_ENDPOINT"),
		AccessKey: os.Getenv("MINIO_ACCESS_KEY"),
		SecretKey: os.Getenv("MINIO_SECRET_KEY"),
		UseSSL:    os.Getenv("MINIO_USE_SSL") == "true",
		Bucket:    os.Getenv("MINIO_BUCKET"),
	}
}

func setupRabbitPublisher(rabbitURL, queue string) interfaces.MessagePublisher {
	if rabbitURL == "" {
		return nil
	}
	p, err := infraMessaging.NewRabbitMQPublisher(rabbitURL)
	if err != nil {
		log.Printf("warning: rabbitmq publisher init failed: %v", err)
		return nil
	}
	maxLen := atoiOrDefault(os.Getenv("RABBITMQ_QUEUE_MAXLEN"), 1000)
	_ = p.EnsureQueue(queue, maxLen, false)
	return p
}

// setupRedisCacheFromEnv initializes a Redis-backed cache if REDIS_ADDR is set.
// It uses the same env var names as Workers/AdminCache for consistency:
// - REDIS_ADDR (e.g., "redis:6379")
// - CACHE_PREFIX (default: "videorank:")
// - CACHE_TTL_SECONDS (default: 120)
func setupRedisCacheFromEnv() interfaces.Cache {
	addr := os.Getenv("REDIS_ADDR")
	if addr == "" {
		return nil
	}
	prefix := getEnvOrDefault("CACHE_PREFIX", "videorank:")
	ttl := atoiOrDefault(os.Getenv("CACHE_TTL_SECONDS"), 120)
	rdb := infraCache.MustRedisClient(addr)
	return infraCache.NewRedisCache(rdb, prefix, ttl)
}

// setupAggregatesFromEnv initializes Redis Aggregates for leaderboards/stats.
// Uses the same REDIS_ADDR and CACHE_PREFIX. Optional tenant via TENANT_ID.
func setupAggregatesFromEnv() interfaces.Aggregates {
	addr := os.Getenv("REDIS_ADDR")
	if addr == "" {
		return nil
	}
	prefix := getEnvOrDefault("CACHE_PREFIX", "videorank:")
	tenant := os.Getenv("TENANT_ID")
	rdb := infraCache.MustRedisClient(addr)
	return infraCache.NewRedisAggregates(rdb, prefix, tenant)
}

func main() {
	dsn := os.Getenv("DATABASE_URL")
	if err := runMigrations(dsn, "file://internal/infrastructure/migrations"); err != nil {
		log.Fatalf("%v", err)
	}

	db, err := openDB(dsn)
	if err != nil {
		log.Fatalf("database connection failed: %v", err)
	}

	jwtSecret := getEnvOrDefault("JWT_SECRET", "secret")

	userRepo := postgresrepo.NewUserRepository(db)
	locRepo := postgresrepo.NewLocationRepository(db)
	videoRepo := postgresrepo.NewVideoRepository(db)

	minioCfg := loadMinioConfigFromEnv()
	videoStorage, err := storage.NewMinioVideoStorage(minioCfg)
	if err != nil {
		log.Fatalf("minio storage init failed: %v", err)
	}

	authService := useCase.NewAuthService(userRepo, jwtSecret)
	userService := useCase.NewUserService(userRepo, locRepo)
	locationService := useCase.NewLocationService(locRepo)
	publicRepo := postgresrepo.NewPublicRepository(db)
	voteRepo := postgresrepo.NewVoteRepository(db)

	audioQueue := getEnvOrDefault("STATES_MACHINE_QUEUE", "states_machine_queue")
	messagePublisher := setupRabbitPublisher(os.Getenv("RABBITMQ_URL"), audioQueue)
	defer func() {
		if messagePublisher != nil {
			messagePublisher.Close()
		}
	}()

	// Build use cases (inject publisher into use case, not handlers)
	uploadsUC := useCase.NewUploadsUseCase(videoRepo, videoStorage, messagePublisher, audioQueue)

	r := gin.Default()
	r.Static("/static", "./static")
	statusService := useCase.NewStatusService()
	// Redis cache and idempotency TTL
	cache := setupRedisCacheFromEnv()
	idemTTL := atoiOrDefault(os.Getenv("IDEMPOTENCY_TTL_SECONDS"), 3600)
	// Redis aggregates (leaderboards/stats)
	aggregates := setupAggregatesFromEnv()
	publicService := useCase.NewPublicServiceWithAgg(publicRepo, voteRepo, aggregates)

	handlers.NewRouter(r, handlers.RouterConfig{
		AuthService:     authService,
		UserService:     userService,
		LocationService: locationService,
		UploadsUC:       uploadsUC,
		PublicService:   publicService,
		StatusService:   statusService,
		JWTSecret:       jwtSecret,
		Cache:           cache,
		IdemTTLSeconds:  idemTTL,
		Aggregates:      aggregates,
	})

	port := getEnvOrDefault("PORT", "8080")
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
