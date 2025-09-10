package mocks

import "errors"

type NotificationServiceMock struct {
	NotifyVideoProcessedFunc    func(videoID, filename, bucketPath string) error
	NotifyProcessingCompleteFunc func(videoID string, success bool) error
	ShouldFail                  bool
	VideoProcessedCalls         []VideoProcessedCall
	ProcessingCompleteCalls     []ProcessingCompleteCall
}

type VideoProcessedCall struct {
	VideoID    string
	Filename   string
	BucketPath string
}

type ProcessingCompleteCall struct {
	VideoID string
	Success bool
}

func NewNotificationServiceMock() *NotificationServiceMock {
	return &NotificationServiceMock{
		VideoProcessedCalls:     make([]VideoProcessedCall, 0),
		ProcessingCompleteCalls: make([]ProcessingCompleteCall, 0),
	}
}

func (m *NotificationServiceMock) NotifyVideoProcessed(videoID, filename, bucketPath string) error {
	m.VideoProcessedCalls = append(m.VideoProcessedCalls, VideoProcessedCall{
		VideoID:    videoID,
		Filename:   filename,
		BucketPath: bucketPath,
	})
	
	if m.NotifyVideoProcessedFunc != nil {
		return m.NotifyVideoProcessedFunc(videoID, filename, bucketPath)
	}
	
	if m.ShouldFail {
		return errors.New("notification failed")
	}
	
	return nil
}

func (m *NotificationServiceMock) NotifyProcessingComplete(videoID string, success bool) error {
	m.ProcessingCompleteCalls = append(m.ProcessingCompleteCalls, ProcessingCompleteCall{
		VideoID: videoID,
		Success: success,
	})
	
	if m.NotifyProcessingCompleteFunc != nil {
		return m.NotifyProcessingCompleteFunc(videoID, success)
	}
	
	if m.ShouldFail {
		return errors.New("notification failed")
	}
	
	return nil
}