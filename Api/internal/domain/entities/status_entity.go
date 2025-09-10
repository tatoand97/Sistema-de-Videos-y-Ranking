package entities

type VideoStatus string

const (
	StatusUploaded         VideoStatus = "UPLOADED"
	StatusTrimming         VideoStatus = "TRIMMING"
	StatusAdjustingRes     VideoStatus = "ADJUSTING_RESOLUTION"
	StatusAddingWatermark  VideoStatus = "ADDING_WATERMARK"
	StatusRemovingAudio    VideoStatus = "REMOVING_AUDIO"
	StatusAddingIntroOutro VideoStatus = "ADDING_INTRO_OUTRO"
	StatusProcessed        VideoStatus = "PROCESSED"
	// StatusPublished indica que el video ya fue publicado para votación pública
	StatusPublished VideoStatus = "PUBLISHED"
	StatusFailed    VideoStatus = "FAILED"
)

func AllVideoStatuses() []VideoStatus {
	return []VideoStatus{
		StatusUploaded,
		StatusTrimming,
		StatusAdjustingRes,
		StatusAddingWatermark,
		StatusRemovingAudio,
		StatusAddingIntroOutro,
		StatusProcessed,
		StatusPublished,
		StatusFailed,
	}
}
