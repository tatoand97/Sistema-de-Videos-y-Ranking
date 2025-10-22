package storage

import (
	"context"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// Config defines the minimal parameters needed to instantiate an S3 client.
type Config struct {
	Region       string
	AccessKey    string
	SecretKey    string
	Endpoint     string
	UsePathStyle bool
}

// Client is a thin wrapper over the AWS SDK S3 client to satisfy StorageService ports.
type Client struct {
	client *s3.Client
}

// NewClient builds a new S3 client. If AccessKey and SecretKey are empty the default provider chain is used.
func NewClient(cfg Config) (*Client, error) {
	if cfg.Region == "" {
		return nil, fmt.Errorf("s3 region is required")
	}

	loadOpts := []func(*config.LoadOptions) error{
		config.WithRegion(cfg.Region),
	}

	if cfg.AccessKey != "" && cfg.SecretKey != "" {
		loadOpts = append(loadOpts, config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(cfg.AccessKey, cfg.SecretKey, ""),
		))
	}

	awsCfg, err := config.LoadDefaultConfig(context.Background(), loadOpts...)
	if err != nil {
		return nil, fmt.Errorf("load aws config: %w", err)
	}

	client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		if cfg.Endpoint != "" {
			o.EndpointResolver = s3.EndpointResolverFromURL(cfg.Endpoint)
		}
		if cfg.UsePathStyle {
			o.UsePathStyle = true
		}
	})

	return &Client{client: client}, nil
}

// GetObject downloads an object from S3 and returns its body reader.
func (c *Client) GetObject(bucket, key string) (io.Reader, error) {
	if bucket == "" || key == "" {
		return nil, fmt.Errorf("bucket and key are required")
	}

	output, err := c.client.GetObject(context.Background(), &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, err
	}
	return output.Body, nil
}

// PutObject uploads the provided reader contents into S3.
func (c *Client) PutObject(bucket, key string, data io.Reader, size int64) error {
	if bucket == "" || key == "" {
		return fmt.Errorf("bucket and key are required")
	}
	if data == nil {
		return fmt.Errorf("data reader cannot be nil")
	}

	_, err := c.client.PutObject(context.Background(), &s3.PutObjectInput{
		Bucket:        aws.String(bucket),
		Key:           aws.String(key),
		Body:          data,
		ContentLength: aws.Int64(size),
	})
	return err
}
