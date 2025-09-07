package domain

type Video struct {
	ID           uint   `gorm:"column:video_id;primaryKey"`
	OriginalFile string `gorm:"column:original_file"`
	Status       string `gorm:"column:status"`
}

func (Video) TableName() string {
	return "video"
}

type VideoStatus string

const (
	StatusUploaded         VideoStatus = "UPLOADED"
	StatusTrimming         VideoStatus = "TRIMMING"
	StatusAdjustingRes     VideoStatus = "ADJUSTING_RESOLUTION"
	StatusAddingWatermark  VideoStatus = "ADDING_WATERMARK"
	StatusRemovingAudio    VideoStatus = "REMOVING_AUDIO"
	StatusAddingIntroOutro VideoStatus = "ADDING_INTRO_OUTRO"
	StatusProcessed        VideoStatus = "PROCESSED"
	StatusFailed           VideoStatus = "FAILED"
)