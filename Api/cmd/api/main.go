package main

import (
	"api/internal/presentation/handlers"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

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

func loadS3ConfigFromEnv() storage.S3Config {
	return storage.S3Config{
		Region:       os.Getenv("AWS_REGION"),
		AccessKey:    os.Getenv("AWS_ACCESS_KEY_ID"),
		SecretKey:    os.Getenv("AWS_SECRET_ACCESS_KEY"),
		Endpoint:     os.Getenv("S3_ENDPOINT"),
		UsePathStyle: strings.EqualFold(os.Getenv("S3_USE_PATH_STYLE"), "true"),
		Bucket:       os.Getenv("S3_BUCKET"),
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
func setupRedisCacheFromEnv() interfaces.Cache {
	addr := os.Getenv("REDIS_ADDR")
	if addr == "" {
		return nil
	}
	prefix := getEnvOrDefault("CACHE_PREFIX", "videorank:")
	rdb := infraCache.MustRedisClient(addr)
	return infraCache.NewRedisCache(rdb, prefix)
}

// Redis Aggregates removed; rankings served from DB.

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

	s3Cfg := loadS3ConfigFromEnv()
	videoStorage, err := storage.NewS3VideoStorage(s3Cfg)
	if err != nil {
		log.Fatalf("s3 storage init failed: %v", err)
	}

	authService := useCase.NewAuthService(userRepo, jwtSecret)
	userService := useCase.NewUserService(userRepo, locRepo)
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

	// Lightweight CORS middleware (avoids external deps)
	allowed := getEnvOrDefault("CORS_ORIGIN", "*")
	var allowAll bool
	var allowList []string
	if allowed == "*" || allowed == "" {
		allowAll = true
	} else {
		allowList = strings.Split(allowed, ",")
	}
	r.Use(func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		if allowAll {
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		} else {
			for _, o := range allowList {
				if strings.EqualFold(strings.TrimSpace(o), origin) {
					c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
					break
				}
			}
		}
		c.Writer.Header().Set("Vary", "Origin")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Authorization,Content-Type,Accept,Origin")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})
	r.Static("/static", "./static")
	statusService := useCase.NewStatusService()
	// Redis cache solo lectura
	cache := setupRedisCacheFromEnv()
	// Public service without Redis aggregates
	publicService := useCase.NewPublicService(publicRepo, voteRepo)

	processedBase := strings.TrimRight(os.Getenv("PROCESSED_VIDEO_BASE_URL"), "/")
	processedVideoURL := ""
	if processedBase != "" {
		processedVideoURL = fmt.Sprintf("%s/%%s", processedBase)
	} else {
		processedVideoURL = fmt.Sprintf("http://%s:%s/processed-videos/%%s",
			getEnvOrDefault("PROCESSED_VIDEO_HOST", "localhost"),
			getEnvOrDefault("PROCESSED_VIDEO_PORT", "8084"))
	}

	handlers.NewRouter(r, handlers.RouterConfig{
		AuthService:        authService,
		UserService:        userService,
		UploadsUC:          uploadsUC,
		PublicService:      publicService,
		StatusService:      statusService,
		JWTSecret:          jwtSecret,
		Cache:              cache,
		CacheSchemaVersion: getEnvOrDefault("SCHEMA_VERSION", "v2"),
		ProcessedVideoURL:  processedVideoURL,
	})

	port := getEnvOrDefault("PORT", "8080")
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
