package config

import "os"

type Config struct {
	Port          string
	JWTSecret     string
	JWTTTLMinutes string
	UploadDir     string
	DatabaseURL   string
	CORSOrigin    string
}

func FromEnv() Config {
	return Config{
		Port:          get("API_PORT", "8080"),
		JWTSecret:     get("JWT_SECRET", "dev"),
		UploadDir:     get("UPLOAD_DIR", "/app/uploads"),
		DatabaseURL:   get("DATABASE_URL", "postgres://todo_user:todo_pass@db:5432/todo_db?sslmode=disable"),
		CORSOrigin:    get("CORS_ORIGIN", "*"),
		JWTTTLMinutes: get("JWT_TTL_MINUTES", "60"),
	}
}

func get(k, def string) string {
	v := os.Getenv(k)
	if v == "" {
		return def
	}
	return v
}
