// internal/infrastructure/config.go
package infrastructure

import (
	"os"
	"strconv"
	"strings"
)

type Config struct {
	AWSRegion         string
	S3Endpoint        string
	S3AccessKey       string
	S3SecretKey       string
	S3SessionToken    string
	S3UsePathStyle    bool
	S3AnonymousAccess bool
	RawBucket         string
	ProcessedBucket   string
	SQSQueueURL       string
	MaxRetries        int
	MaxSeconds        int
	IntroSeconds      float64
	OutroSeconds      float64
	TargetWidth       int
	TargetHeight      int
	FPS               int
	LogoPath          string
}

func LoadConfig() *Config {
	maxRetries := getEnvInt("MAX_RETRIES", 5)
	maxSeconds := getEnvInt("MAX_SECONDS", 30)
	usePathStyle := strings.EqualFold(os.Getenv("S3_USE_PATH_STYLE"), "true")
	anonymousAccess := strings.EqualFold(os.Getenv("S3_ANONYMOUS_ACCESS"), "true")

	intro := getEnvFloat("INTRO_SECONDS", 2.5)
	outro := getEnvFloat("OUTRO_SECONDS", 2.5)
	tw := getEnvInt("TARGET_WIDTH", 1280)
	th := getEnvInt("TARGET_HEIGHT", 720)
	fps := getEnvInt("FPS", 30)

	logo := os.Getenv("LOGO_PATH")
	if logo == "" {
		logo = "./assets/nba-logo-removebg-preview.png"
	}

	return &Config{
		AWSRegion:         os.Getenv("AWS_REGION"),
		S3Endpoint:        os.Getenv("S3_ENDPOINT"),
		S3AccessKey:       os.Getenv("AWS_ACCESS_KEY_ID"),
		S3SecretKey:       os.Getenv("AWS_SECRET_ACCESS_KEY"),
		S3SessionToken:    os.Getenv("AWS_SESSION_TOKEN"),
		S3UsePathStyle:    usePathStyle,
		S3AnonymousAccess: anonymousAccess,
		RawBucket:         os.Getenv("S3_BUCKET_RAW"),
		ProcessedBucket:   os.Getenv("S3_BUCKET_PROCESSED"),
		SQSQueueURL:       os.Getenv("SQS_GOSSIP_QUEUE"),
		MaxRetries:        maxRetries,
		MaxSeconds:        maxSeconds,
		IntroSeconds:      intro,
		OutroSeconds:      outro,
		TargetWidth:       tw,
		TargetHeight:      th,
		FPS:               fps,
		LogoPath:          logo,
	}
}

func getEnvInt(key string, def int) int {
	if v := os.Getenv(key); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return def
}

func getEnvFloat(key string, def float64) float64 {
	if v := os.Getenv(key); v != "" {
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			return f
		}
	}
	return def
}
