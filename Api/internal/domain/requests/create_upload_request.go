package requests

// CreateUploadRequest defines the input to request a signed POST policy upload.
type CreateUploadRequest struct {
	Filename  string `json:"filename"`
	MimeType  string `json:"mimeType"`
	SizeBytes int64  `json:"sizeBytes"`
	Checksum  string `json:"checksum"`
}
