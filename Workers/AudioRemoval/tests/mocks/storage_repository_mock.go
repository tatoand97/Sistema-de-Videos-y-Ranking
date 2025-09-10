package mocks

import "errors"

type StorageRepositoryMock struct {
	DownloadFunc func(bucket, filename string) ([]byte, error)
	UploadFunc   func(bucket, filename string, data []byte) error
	Files        map[string][]byte
	UploadCalls  []UploadCall
}

type UploadCall struct {
	Bucket   string
	Filename string
	Data     []byte
}

func NewStorageRepositoryMock() *StorageRepositoryMock {
	return &StorageRepositoryMock{
		Files:       make(map[string][]byte),
		UploadCalls: make([]UploadCall, 0),
	}
}

func (m *StorageRepositoryMock) Download(bucket, filename string) ([]byte, error) {
	if m.DownloadFunc != nil {
		return m.DownloadFunc(bucket, filename)
	}
	
	key := bucket + "/" + filename
	if data, exists := m.Files[key]; exists {
		return data, nil
	}
	
	return nil, errors.New("file not found")
}

func (m *StorageRepositoryMock) Upload(bucket, filename string, data []byte) error {
	if m.UploadFunc != nil {
		return m.UploadFunc(bucket, filename, data)
	}
	
	key := bucket + "/" + filename
	m.Files[key] = data
	m.UploadCalls = append(m.UploadCalls, UploadCall{
		Bucket:   bucket,
		Filename: filename,
		Data:     data,
	})
	
	return nil
}