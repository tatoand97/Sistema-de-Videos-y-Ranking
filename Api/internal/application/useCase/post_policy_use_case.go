package useCase

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"api/internal/domain"
	"api/internal/domain/interfaces"
	"api/internal/domain/requests"
	"api/internal/domain/responses"
)

type PostPolicyUseCase struct {
	storage interfaces.VideoStorage
}

func NewPostPolicyUseCase(storage interfaces.VideoStorage) *PostPolicyUseCase {
	return &PostPolicyUseCase{storage: storage}
}

var sha256HexRe = regexp.MustCompile(`^[A-Fa-f0-9]{64}$`)

func (uc *PostPolicyUseCase) Execute(ctx context.Context, req requests.CreateUploadRequest) (*responses.CreateUploadResponsePostPolicy, error) {
	// Basic validations moved from handler
	if strings.TrimSpace(req.Filename) == "" || strings.TrimSpace(req.MimeType) == "" {
		return nil, fmt.Errorf("%w: filename and mimeType are required", domain.ErrInvalid)
	}
	if req.SizeBytes < 0 {
		return nil, fmt.Errorf("%w: sizeBytes must be >= 0", domain.ErrInvalid)
	}
	if req.Checksum != "" && !sha256HexRe.MatchString(req.Checksum) {
		return nil, fmt.Errorf("%w: checksum must be SHA-256 hex", domain.ErrInvalid)
	}

	return uc.storage.PresignedPostPolicy(ctx, req)
}
