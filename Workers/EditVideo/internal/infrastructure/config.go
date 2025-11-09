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
	StateMachineQueue string
	MaxRetries        int
	MaxSeconds        int
}

func LoadConfig() *Config {
	maxRetries := 3
	if v := os.Getenv("MAX_RETRIES"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			maxRetries = n
		}
	}
	maxSeconds := 30
	if v := os.Getenv("MAX_SECONDS"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			maxSeconds = n
		}
	}
	usePathStyle := strings.EqualFold(os.Getenv("S3_USE_PATH_STYLE"), "true")
	anonymousAccess := strings.EqualFold(os.Getenv("S3_ANONYMOUS_ACCESS"), "true")
	stateMachineQueue := os.Getenv("SQS_STATE_MACHINE_QUEUE")
	if stateMachineQueue == "" {
		stateMachineQueue = os.Getenv("SQS_STATES_MACHINE_QUEUE")
	}

	return &Config{
		AWSRegion:         os.Getenv("AWS_REGION"),
		S3Endpoint:        os.Getenv("S3_ENDPOINT"),
		S3AccessKey:       os.Getenv("AWS_ACCESS_KEY_ID"),
		S3SecretKey:       os.Getenv("AWS_SECRET_ACCESS_KEY"),
		S3SessionToken:    os.Getenv("AWS_SESSION_TOKEN"),
		S3UsePathStyle:    usePathStyle,
		S3AnonymousAccess: anonymousAccess,
		RawBucket:         os.Getenv("RAW_BUCKET"),
		ProcessedBucket:   os.Getenv("PROCESSED_BUCKET"),
		SQSQueueURL:       os.Getenv("SQS_EDIT_VIDEO_QUEUE"),
		StateMachineQueue: stateMachineQueue,
		MaxRetries:        maxRetries,
		MaxSeconds:        maxSeconds,
	}
}
