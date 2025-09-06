package storage

import (
	"context"
	"io"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"api/internal/domain/interfaces"
	"api/internal/domain/requests"
	"api/internal/domain/responses"

	"github.com/google/uuid"
	minio "github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// MinioConfig holds the necessary configuration for connecting to MinIO.
type MinioConfig struct {
	Endpoint  string
	AccessKey string
	SecretKey string
	UseSSL    bool
	Bucket    string
}

type videoStorage struct {
	client   *minio.Client
	bucket   string
	endpoint string
	useSSL   bool
}

// NewMinioVideoStorage creates a new VideoStorage backed by MinIO.
func NewMinioVideoStorage(cfg MinioConfig) (interfaces.VideoStorage, error) {
	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: cfg.UseSSL,
	})
	if err != nil {
		return nil, err
	}
	return &videoStorage{client: client, bucket: cfg.Bucket, endpoint: cfg.Endpoint, useSSL: cfg.UseSSL}, nil
}

// Save uploads the provided video data to MinIO and returns the object name.
func (s *videoStorage) Save(ctx context.Context, objectName string, reader io.Reader, size int64, contentType string) (string, error) {
	_, err := s.client.PutObject(ctx, s.bucket, objectName, reader, size, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		return "", err
	}
	return objectName, nil
}

// PresignedPostPolicy builds and signs a POST policy for direct uploads.
func (s *videoStorage) PresignedPostPolicy(ctx context.Context, req requests.CreateUploadRequest) (*responses.CreateUploadResponsePostPolicy, error) {
	now := time.Now().UTC()
	key := buildObjectKey(req.Filename, now)

	expires := now.Add(15 * time.Minute)
	policy := minio.NewPostPolicy()
	_ = policy.SetBucket(s.bucket)
	_ = policy.SetKey(key)
	_ = policy.SetExpires(expires)

	var maxSize int64 = 100 * 1024 * 1024
	if req.SizeBytes > 0 {
		maxSize = req.SizeBytes
	}
	_ = policy.SetContentLengthRange(0, maxSize)
	_ = policy.SetContentType(req.MimeType)

	if strings.TrimSpace(req.Checksum) != "" {
		_ = policy.SetUserMetadata("sha256", req.Checksum)
	}

	u, formData, err := s.client.PresignedPostPolicy(ctx, policy)
	if err != nil {
		return nil, err
	}

	// Ensure fields include our desired values
	formData["key"] = key
	formData["Content-Type"] = req.MimeType
	if strings.TrimSpace(req.Checksum) != "" {
		formData["x-amz-meta-sha256"] = req.Checksum
	}
	formData["success_action_status"] = "201"

	resp := &responses.CreateUploadResponsePostPolicy{
		UploadURL:   u.String(),
		ResourceURL: "s3://" + s.bucket + "/" + key,
		ExpiresAt:   expires.Format(time.RFC3339),
		Form: responses.S3PostPolicyForm{
			Key:               formData["key"],
			Policy:            formData["policy"],
			Algorithm:         formData["x-amz-algorithm"],
			Credential:        formData["x-amz-credential"],
			Date:              formData["x-amz-date"],
			Signature:         formData["x-amz-signature"],
			ContentType:       formData["Content-Type"],
			SuccessActionCode: formData["success_action_status"],
		},
	}
	if v, ok := formData["x-amz-meta-sha256"]; ok {
		resp.Form.MetaSHA256 = v
	}
	return resp, nil
}

var sanitizeRe = regexp.MustCompile(`[^a-zA-Z0-9._-]+`)

func buildObjectKey(filename string, now time.Time) string {
	base := filepath.Base(filename)
	if base == "." || base == ".." || base == "" {
		base = "file"
	}
	base = strings.TrimSpace(base)
	base = strings.ReplaceAll(base, " ", "-")
	base = sanitizeRe.ReplaceAllString(base, "")
	if base == "" {
		base = "file"
	}
	id := uuid.New().String()
	yyyy, mm, dd := now.Date()
	return strings.Join([]string{
		"uploads",
		strconv.Itoa(yyyy),
		twoDigits(int(mm)),
		twoDigits(dd),
		id + "-" + base,
	}, "/")
}

func twoDigits(v int) string {
	if v < 10 {
		return "0" + strconv.Itoa(v)
	}
	return strconv.Itoa(v)
}
