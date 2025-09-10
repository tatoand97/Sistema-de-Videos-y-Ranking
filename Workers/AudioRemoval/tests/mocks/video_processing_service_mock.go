package mocks

import "errors"

type VideoProcessingServiceMock struct {
	RemoveAudioFunc func(inputData []byte) ([]byte, error)
	ProcessedData   []byte
	ShouldFail      bool
	CallCount       int
}

func NewVideoProcessingServiceMock() *VideoProcessingServiceMock {
	return &VideoProcessingServiceMock{
		ProcessedData: []byte("processed video data"),
	}
}

func (m *VideoProcessingServiceMock) RemoveAudio(inputData []byte) ([]byte, error) {
	m.CallCount++
	
	if m.RemoveAudioFunc != nil {
		return m.RemoveAudioFunc(inputData)
	}
	
	if m.ShouldFail {
		return nil, errors.New("processing failed")
	}
	
	return m.ProcessedData, nil
}