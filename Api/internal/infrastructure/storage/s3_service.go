package storage

import (
	"context"
	"fmt"
	"io"

	"api/internal/domain/interfaces"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// S3Config holds the necessary configuration for connecting to Amazon S3 or an S3-compatible service.
type S3Config struct {
	Region          string
	AccessKey       string
	SecretKey       string
	Endpoint        string
	UsePathStyle    bool
	Bucket          string
	AnonymousAccess bool
}

type videoStorage struct {
	client *s3.Client
	bucket string
}

// NewS3VideoStorage creates a new VideoStorage backed by Amazon S3.
func NewS3VideoStorage(cfg S3Config) (interfaces.VideoStorage, error) {
	if cfg.Region == "" {
		return nil, fmt.Errorf("s3 region is required")
	}
	if cfg.Bucket == "" {
		return nil, fmt.Errorf("s3 bucket is required")
	}

	loadOpts := []func(*config.LoadOptions) error{
		config.WithRegion(cfg.Region),
	}

	var customCreds aws.CredentialsProvider

	// Allow anonymous credentials when explicitly requested, otherwise honor static credentials
	// and finally fall back to the default credential chain (e.g., IAM roles).
	if cfg.AnonymousAccess {
		customCreds = aws.AnonymousCredentials{}
	} else if cfg.AccessKey != "" && cfg.SecretKey != "" {
		customCreds = credentials.NewStaticCredentialsProvider(cfg.AccessKey, cfg.SecretKey, "")
	}

	awsCfg, err := config.LoadDefaultConfig(context.Background(), loadOpts...)
	if err != nil {
		return nil, fmt.Errorf("load aws config: %w", err)
	}

	client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		if customCreds != nil {
			o.Credentials = customCreds
		}
		if cfg.Endpoint != "" {
			o.EndpointResolver = s3.EndpointResolverFromURL(cfg.Endpoint)
		}
		if cfg.UsePathStyle {
			o.UsePathStyle = true
		}
	})

	return &videoStorage{client: client, bucket: cfg.Bucket}, nil
}

// Save uploads the provided video data to S3 and returns the object key.
func (s *videoStorage) Save(ctx context.Context, objectName string, reader io.Reader, size int64, contentType string) (string, error) {
	input := &s3.PutObjectInput{
		Bucket:        aws.String(s.bucket),
		Key:           aws.String(objectName),
		Body:          reader,
		ContentLength: aws.Int64(size),
	}
	if contentType != "" {
		input.ContentType = aws.String(contentType)
	}

	if _, err := s.client.PutObject(ctx, input); err != nil {
		return "", err
	}
	return objectName, nil
}
