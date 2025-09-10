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
	response := responses.PublicVideoResponse{
		VideoID:      1,
		Title:        "Public Video",
		Username:     "testuser",
		ProcessedURL: "https://example.com/public.mp4",
		Votes:        42,
	}

	assert.Equal(t, uint(1), response.VideoID)
	assert.Equal(t, "Public Video", response.Title)
	assert.Equal(t, "testuser", response.Username)
	assert.Equal(t, "https://example.com/public.mp4", response.ProcessedURL)
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
	user := responses.UserBasic{
		UserID:    123,
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john@example.com",
	}

	assert.Equal(t, uint(123), user.UserID)
	assert.Equal(t, "John", user.FirstName)
	assert.Equal(t, "Doe", user.LastName)
	assert.Equal(t, "john@example.com", user.Email)
}

func TestCreateUploadResponsePostPolicy_Fields(t *testing.T) {
	fields := map[string]string{
		"key":    "uploads/video.mp4",
		"policy": "base64-encoded-policy",
	}

	response := responses.CreateUploadResponsePostPolicy{
		URL:    "https://s3.amazonaws.com/bucket",
		Fields: fields,
	}

	assert.Equal(t, "https://s3.amazonaws.com/bucket", response.URL)
	assert.Equal(t, fields, response.Fields)
	assert.Equal(t, "uploads/video.mp4", response.Fields["key"])
	assert.Equal(t, "base64-encoded-policy", response.Fields["policy"])
}