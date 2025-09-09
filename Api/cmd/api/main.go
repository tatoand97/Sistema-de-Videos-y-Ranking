package main

import (
	"api/internal/presentation/handlers"
	"errors"
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
	infraMessaging "api/internal/infrastructure/messaging"
	postgresrepo "api/internal/infrastructure/repository"
	"api/internal/infrastructure/storage"
)

func main() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://app_user:app_password@localhost:5432/videorank?sslmode=disable"
	}

	migPath := "file://internal/infrastructure/migrations"

	m, err := migrate.New(migPath, dsn)
	if err != nil {
		log.Fatalf("migrate init failed: %v", err)
	}
	defer func() {
		srcErr, dbErr := m.Close()
		if srcErr != nil || dbErr != nil {
			log.Printf("migrate close warnings: src=%v db=%v", srcErr, dbErr)
		}
	}()

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatalf("migrate up failed: %v", err)
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("database connection failed: %v", err)
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "secret"
	}

	userRepo := postgresrepo.NewUserRepository(db)
	locRepo := postgresrepo.NewLocationRepository(db)
	videoRepo := postgresrepo.NewVideoRepository(db)

	minioCfg := storage.MinioConfig{
		Endpoint:  os.Getenv("MINIO_ENDPOINT"),   // e.g. "localhost:9000"
		AccessKey: os.Getenv("MINIO_ACCESS_KEY"), // e.g. "minio"
		SecretKey: os.Getenv("MINIO_SECRET_KEY"), // e.g. "minio12345"
		UseSSL:    os.Getenv("MINIO_USE_SSL") == "true",
		Bucket:    os.Getenv("MINIO_BUCKET"),
	}
	videoStorage, err := storage.NewMinioVideoStorage(minioCfg)
	if err != nil {
		log.Fatalf("minio storage init failed: %v", err)
	}

	authService := useCase.NewAuthService(userRepo, jwtSecret)
	userService := useCase.NewUserService(userRepo, locRepo)
	locationService := useCase.NewLocationService(locRepo)
	publicRepo := postgresrepo.NewPublicRepository(db)
	voteRepo := postgresrepo.NewVoteRepository(db)
	publicService := useCase.NewPublicService(publicRepo, voteRepo)

	// Optional RabbitMQ publisher for background processing
	var messagePublisher interfaces.MessagePublisher
	rabbitURL := os.Getenv("RABBITMQ_URL")
	audioQueue := os.Getenv("STATES_MACHINE_QUEUE")
	if audioQueue == "" {
		audioQueue = "states_machine_queue"
	}

	if rabbitURL != "" {
		p, err := infraMessaging.NewRabbitMQPublisher(rabbitURL)
		if err != nil {
			log.Printf("warning: rabbitmq publisher init failed: %v", err)
		} else {
			// Ensure queue infra roughly mirrors workers (DLQ + length)
			maxLen := 1000
			if v := os.Getenv("RABBITMQ_QUEUE_MAXLEN"); v != "" {
				if n, err := strconv.Atoi(v); err == nil {
					maxLen = n
				}
			}
			_ = p.EnsureQueue(audioQueue, maxLen, false)
			messagePublisher = p
			defer messagePublisher.Close()
		}
	}

	// Build use cases (inject publisher into use case, not handlers)
	uploadsUC := useCase.NewUploadsUseCase(videoRepo, videoStorage, messagePublisher, audioQueue)

	r := gin.Default()
	r.Static("/static", "./static")
	statusService := useCase.NewStatusService()
	handlers.NewRouter(r, authService, userService, locationService, jwtSecret, uploadsUC, publicService, statusService)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
