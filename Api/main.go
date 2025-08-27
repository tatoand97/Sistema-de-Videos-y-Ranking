package main

import (
	"log"
	"net/http"

	httpadp "main_prj/internal/adapters/http"
	"main_prj/internal/adapters/postgres"
	"main_prj/internal/config"
	"main_prj/internal/services"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()
	cfg := config.FromEnv()

	db, err := postgres.Open(cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}

	ur := postgres.NewUserRepo(db)

	authSvc := services.NewAuthService(ur, cfg.JWTSecret)

	usersHandler := &httpadp.UsersHandler{Auth: authSvc, UploadDir: cfg.UploadDir}

	r := httpadp.SetupRouter(usersHandler, cfg.CORSOrigin, cfg.JWTSecret)

	log.Printf("API listening on :%s", cfg.Port)
	http.ListenAndServe(":"+cfg.Port, r)
}
