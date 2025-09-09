package domain

// OrchestrateVideoUseCaseInterface defines the interface for video orchestration use case
type OrchestrateVideoUseCaseInterface interface {
	Execute(videoID string) error
	HandleTrimCompleted(videoID, filename string) error
	HandleEditCompleted(videoID, filename string) error
	HandleAudioRemovalCompleted(videoID, filename string) error
	HandleWatermarkingCompleted(videoID, filename string) error
	HandleGossipOpenCloseCompleted(videoID, filename string) error
	GetRetryDelayMinutes() int
}