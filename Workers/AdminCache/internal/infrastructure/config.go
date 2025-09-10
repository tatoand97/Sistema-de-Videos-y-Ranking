package infrastructure

import (
	"os"
	"strconv"
	"strings"
)

type Config struct {
    RedisAddr       string
    CachePrefix     string
    CacheTTLSeconds int

    PostgresDSN string

    RefreshIntervalSeconds int
    WarmPages              int
    WarmCities             []string
    PageSizeDefault        int
}

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" { return v }
	return def
}
func getenvInt(key string, def int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil { return n }
	}
	return def
}

func LoadConfig() Config {
    cfg := Config{
        RedisAddr:              getenv("REDIS_ADDR", "redis:6379"),
        CachePrefix:            getenv("CACHE_PREFIX", "videorank:"),
        CacheTTLSeconds:        getenvInt("CACHE_TTL_SECONDS", 120),
        PostgresDSN:            getenv("POSTGRES_DSN", "postgres://app_user:app_password@postgres:5432/videorank?sslmode=disable"),
        RefreshIntervalSeconds: getenvInt("REFRESH_INTERVAL_SECONDS", 60),
        WarmPages:              getenvInt("WARM_PAGES", 3),
        PageSizeDefault:        getenvInt("PAGE_SIZE_DEFAULT", 20),
    }
	cities := getenv("WARM_CITIES", "bogota,medellin,cali")
	if cities != "" {
		for _, c := range strings.Split(cities, ",") {
			c = strings.TrimSpace(c)
			if c != "" { cfg.WarmCities = append(cfg.WarmCities, c) }
		}
	}
	if cfg.CacheTTLSeconds <= 0 { cfg.CacheTTLSeconds = 120 }
	return cfg
}
