package domain

import "time"

type Video struct {
	ID          string
	Filename    string
	Status      ProcessingStatus
	CreatedAt   time.Time
	ProcessedAt *time.Time
}

type ProcessingStatus string

const (
	StatusPending    ProcessingStatus = "pending"
	StatusProcessing ProcessingStatus = "processing"
	StatusCompleted  ProcessingStatus = "completed"
	StatusFailed     ProcessingStatus = "failed"
)
