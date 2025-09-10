package responses_test

import (
	"testing"
	"time"

	"api/internal/domain/responses"

	"github.com/stretchr/testify/assert"
)

func TestVideoResponse_Fields(t *testing.T) {
	now := time.Now()
	processedURL := "https://example.com/video.mp4"

	response := responses.VideoResponse{
		VideoID:      "123",
		Title:        "Test Video",
		Status:       "processed",
		UploadedAt:   now,
		ProcessedAt:  &now,
		ProcessedURL: &processedURL,
	}

	assert.Equal(t, "123", response.VideoID)
	assert.Equal(t, "Test Video", response.Title)
	assert.Equal(t, "processed", response.Status)
	assert.Equal(t, now, response.UploadedAt)
	assert.Equal(t, &now, response.ProcessedAt)
	assert.Equal(t, &processedURL, response.ProcessedURL)
}

func TestPublicVideoResponse_Fields(t *testing.T) {
	processedURL := "https://example.com/public.mp4"
	city := "Bogotá"
	response := responses.PublicVideoResponse{
		VideoID:      1,
		Title:        "Public Video",
		ProcessedURL: &processedURL,
		City:         &city,
		Votes:        42,
	}

	assert.Equal(t, uint(1), response.VideoID)
	assert.Equal(t, "Public Video", response.Title)
	assert.Equal(t, &processedURL, response.ProcessedURL)
	assert.Equal(t, &city, response.City)
	assert.Equal(t, 42, response.Votes)
}

func TestRankingItem_Fields(t *testing.T) {
	city := "Bogotá"
	item := responses.RankingItem{
		Username: "user1",
		City:     &city,
		Votes:    100,
	}

	assert.Equal(t, "user1", item.Username)
	assert.Equal(t, &city, item.City)
	assert.Equal(t, 100, item.Votes)
}

func TestRankingEntry_Fields(t *testing.T) {
	city := "Medellín"
	entry := responses.RankingEntry{
		Position: 1,
		Username: "topuser",
		City:     &city,
		Votes:    500,
	}

	assert.Equal(t, 1, entry.Position)
	assert.Equal(t, "topuser", entry.Username)
	assert.Equal(t, &city, entry.City)
	assert.Equal(t, 500, entry.Votes)
}

func TestUserBasic_Fields(t *testing.T) {
	city := "Bogotá"
	user := responses.UserBasic{
		UserID:   123,
		Username: "johndoe",
		City:     &city,
	}

	assert.Equal(t, uint(123), user.UserID)
	assert.Equal(t, "johndoe", user.Username)
	assert.Equal(t, &city, user.City)
}

func TestCreateUploadResponsePostPolicy_Fields(t *testing.T) {
	form := responses.S3PostPolicyForm{
		Key:               "uploads/video.mp4",
		Policy:            "base64-encoded-policy",
		Algorithm:         "AWS4-HMAC-SHA256",
		Credential:        "credential",
		Date:              "20240101T000000Z",
		Signature:         "signature",
		ContentType:       "video/mp4",
		SuccessActionCode: "201",
	}

	response := responses.CreateUploadResponsePostPolicy{
		UploadURL:   "https://s3.amazonaws.com/bucket",
		ResourceURL: "https://s3.amazonaws.com/bucket/uploads/video.mp4",
		ExpiresAt:   "2024-01-01T01:00:00Z",
		Form:        form,
	}

	assert.Equal(t, "https://s3.amazonaws.com/bucket", response.UploadURL)
	assert.Equal(t, "https://s3.amazonaws.com/bucket/uploads/video.mp4", response.ResourceURL)
	assert.Equal(t, form, response.Form)
	assert.Equal(t, "uploads/video.mp4", response.Form.Key)
	assert.Equal(t, "base64-encoded-policy", response.Form.Policy)
}