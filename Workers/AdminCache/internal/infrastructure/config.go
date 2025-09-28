package infrastructure

import (
	"os"
	"strconv"
	"strings"
)

type Config struct {
	RedisAddr        string
	CachePrefix      string
	SchemaVersion    string
	TTLFreshSeconds  int
	MaxStaleSeconds  int
	LockLeaseSeconds int
	JitterPercent    int

	PostgresDSN          string
	DBReadTimeoutSeconds int
	DBMaxRetries         int

	RefreshIntervalSeconds int
	CityBatchSize          int
	WarmCities             []string

	MaxTopUsers int
}

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func getenvInt(key string, def int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return def
}

func LoadConfig() Config {
	cfg := Config{
		RedisAddr:              getenv("REDIS_ADDR", "redis:6379"),
		CachePrefix:            getenv("CACHE_PREFIX", "videorank:"),
		SchemaVersion:          getenv("SCHEMA_VERSION", "v2"),
		TTLFreshSeconds:        getenvInt("CACHE_TTL_FRESH_SECONDS", 900),
		MaxStaleSeconds:        getenvInt("CACHE_MAX_STALE_SECONDS", 600),
		LockLeaseSeconds:       getenvInt("CACHE_LOCK_LEASE_SECONDS", 10),
		JitterPercent:          getenvInt("CACHE_JITTER_PERCENT", 10),
		PostgresDSN:            getenv("POSTGRES_DSN", "postgres://app_user:app_password@postgres:5432/videorank?sslmode=disable"),
		DBReadTimeoutSeconds:   getenvInt("DB_READ_TIMEOUT_SECONDS", 3),
		DBMaxRetries:           getenvInt("DB_MAX_RETRIES", 3),
		RefreshIntervalSeconds: getenvInt("REFRESH_INTERVAL_SECONDS", 300),
		CityBatchSize:          getenvInt("BATCH_SIZE_CITIES", 50),
		MaxTopUsers:            getenvInt("CACHE_MAX_TOP_USERS", 10),
	}

	cities := getenv("WARM_CITIES", "")
	if cities != "" {
		for _, c := range strings.Split(cities, ",") {
			c = strings.TrimSpace(c)
			if c != "" {
				cfg.WarmCities = append(cfg.WarmCities, c)
			}
		}
	}

	if cfg.TTLFreshSeconds <= 0 {
		cfg.TTLFreshSeconds = 900
	}
	if cfg.MaxStaleSeconds < 0 {
		cfg.MaxStaleSeconds = 0
	}
	if cfg.LockLeaseSeconds <= 0 {
		cfg.LockLeaseSeconds = 10
	}
	if cfg.JitterPercent < 0 {
		cfg.JitterPercent = 0
	}
	if cfg.DBReadTimeoutSeconds <= 0 {
		cfg.DBReadTimeoutSeconds = 3
	}
	if cfg.DBMaxRetries < 1 {
		cfg.DBMaxRetries = 1
	}
	if cfg.RefreshIntervalSeconds <= 0 {
		cfg.RefreshIntervalSeconds = 300
	}
	if cfg.CityBatchSize <= 0 {
		cfg.CityBatchSize = 50
	}
	if cfg.MaxTopUsers <= 0 || cfg.MaxTopUsers > 10 {
		cfg.MaxTopUsers = 10
	}

	return cfg
}
