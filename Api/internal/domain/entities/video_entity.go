package entities

import "time"

type Video struct {
	VideoID       uint       `gorm:"column:video_id;primaryKey;autoIncrement"`
	UserID        uint       `gorm:"column:user_id;not null;index;constraint:OnDelete:CASCADE"`
	Title         string     `gorm:"column:title;size:255;not null"`
	OriginalFile  string     `gorm:"column:original_file;size:255;not null"`
	ProcessedFile *string    `gorm:"column:processed_file;size:255"`
	Status        string     `gorm:"column:status;type:video_status;not null;default:UPLOADED"`
	UploadedAt    time.Time  `gorm:"column:uploaded_at;not null;autoCreateTime"`
	ProcessedAt   *time.Time `gorm:"column:processed_at"`
}

func (Video) TableName() string { return "video" }
