package responses

import "time"

// VideoResponse aligns with the OpenAPI Video schema for private listing
type VideoResponse struct {
	VideoID      string     `json:"video_id"`
	Title        string     `json:"title"`
	Status       string     `json:"status"`
	UploadedAt   time.Time  `json:"uploaded_at"`
	ProcessedAt  *time.Time `json:"processed_at,omitempty"`
	OriginalURL  *string    `json:"original_url,omitempty"`
	ProcessedURL *string    `json:"processed_url,omitempty"`
}
