package responses

// PublicVideoResponse representa el esquema de salida para /api/public/videos
type PublicVideoResponse struct {
	VideoID      uint    `json:"video_id"`
	Title        string  `json:"title"`
	ProcessedURL *string `json:"processed_url"`
	City         *string `json:"city"`
	Votes        int     `json:"votes"`
	OwnerUserID  uint    `json:"-" gorm:"column:owner_user_id"`
}
