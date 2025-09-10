package mocks

import (
	"audioremoval/internal/domain"
	"errors"
)

type VideoRepositoryMock struct {
	FindByFilenameFunc func(filename string) (*domain.Video, error)
	UpdateStatusFunc   func(id string, status domain.ProcessingStatus) error
	Videos             map[string]*domain.Video
	StatusUpdates      map[string]domain.ProcessingStatus
}

func NewVideoRepositoryMock() *VideoRepositoryMock {
	return &VideoRepositoryMock{
		Videos:        make(map[string]*domain.Video),
		StatusUpdates: make(map[string]domain.ProcessingStatus),
	}
}

func (m *VideoRepositoryMock) FindByFilename(filename string) (*domain.Video, error) {
	if m.FindByFilenameFunc != nil {
		return m.FindByFilenameFunc(filename)
	}
	
	for _, video := range m.Videos {
		if video.Filename == filename {
			return video, nil
		}
	}
	
	return nil, errors.New("video not found")
}

func (m *VideoRepositoryMock) UpdateStatus(id string, status domain.ProcessingStatus) error {
	if m.UpdateStatusFunc != nil {
		return m.UpdateStatusFunc(id, status)
	}
	
	m.StatusUpdates[id] = status
	if video, exists := m.Videos[id]; exists {
		video.Status = status
	}
	
	return nil
}