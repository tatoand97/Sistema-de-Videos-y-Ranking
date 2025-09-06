package domain

type Video struct {
	ID       string
	Filename string
	Status   VideoStatus
}

type VideoStatus string

const (
	StatusPending    VideoStatus = "pending"
	StatusProcessing VideoStatus = "processing"
	StatusCompleted  VideoStatus = "completed"
	StatusFailed     VideoStatus = "failed"
)