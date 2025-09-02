package entities

import "time"

// Video representa la entidad de video según la migración SQL.
type Video struct {
	VideoID       uint      `gorm:"primaryKey;autoIncrement"`
	UserID        uint      `gorm:"not null;index"`
	Title         string    `gorm:"size:255;not null"`
	OriginalFile  string    `gorm:"size:255;not null"`
	ProcessedFile *string   `gorm:"size:255"`
	StatusID      uint      `gorm:"not null;index"`
	UploadedAt    time.Time `gorm:"not null;autoCreateTime"`
	ProcessedAt   *time.Time
}
