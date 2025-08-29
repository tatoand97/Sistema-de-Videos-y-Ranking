package main_videork

import (
	"log"
	"main_videork/internal/application/useCase"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	postgresrepo "main_videork/internal/infrastructure/repository"
	"main_videork/internal/presentation"
)

func main() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://app_user:app_password@localhost:5432/videorank?sslmode=disable"
	}

	m, err := migrate.New("file://internal/infrastructure/migrations", dsn)
	if err != nil {
		log.Fatalf("migrate init failed: %v", err)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
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

	// Initialize repositories
	userRepo := postgresrepo.NewUserRepository(db)

	// Initialize services
	authService := useCase.NewAuthService(userRepo, jwtSecret)

	// Setup Gin router
	r := gin.Default()

	r.Static("/static", "./static")

	presentation.NewRouter(r, authService, jwtSecret)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	if err := r.Run(":" + port); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
