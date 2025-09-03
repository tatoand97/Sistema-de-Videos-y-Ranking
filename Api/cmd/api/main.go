package main

import (
	"errors"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"main_videork/internal/application/useCase"
	postgresrepo "main_videork/internal/infrastructure/repository"
	"main_videork/internal/infrastructure/storage"
	"main_videork/internal/presentation"
)

func main() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://app_user:app_password@localhost:5432/videorank?sslmode=disable"
	}

	migPath := "file://Api/internal/infrastructure/migrations"

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

	uploadVideoUC := useCase.NewUploadVideoUseCase(videoRepo, videoStorage)
	authService := useCase.NewAuthService(userRepo, jwtSecret)
	userService := useCase.NewUserService(userRepo)

	r := gin.Default()
	r.Static("/static", "./static")
	presentation.NewRouter(r, authService, userService, jwtSecret, uploadVideoUC)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
