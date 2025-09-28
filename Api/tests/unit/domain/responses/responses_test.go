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
	city := "BogotÃ¡"
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
	city := "BogotÃ¡"
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
	city := "MedellÃ­n"
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
	city := "BogotÃ¡"
	user := responses.UserBasic{
		UserID:   123,
		Username: "johndoe",
		City:     &city,
	}

	assert.Equal(t, uint(123), user.UserID)
	assert.Equal(t, "johndoe", user.Username)
	assert.Equal(t, &city, user.City)
}
