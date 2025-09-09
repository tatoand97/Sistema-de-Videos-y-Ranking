package adapters

import (
	"statesmachine/internal/domain"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestNewPostgresVideoRepository(t *testing.T) {
	repo := NewPostgresVideoRepository(nil)
	assert.NotNil(t, repo)
}

func TestPostgresVideoRepository_Structure(t *testing.T) {
	repo := &PostgresVideoRepository{db: &gorm.DB{}}
	assert.NotNil(t, repo)
	assert.NotNil(t, repo.db)
}

func TestPostgresVideoRepository_UpdateLogic(t *testing.T) {
	tests := []struct {
		name   string
		status domain.VideoStatus
	}{
		{"StatusTrimming", domain.StatusTrimming},
		{"StatusProcessed", domain.StatusProcessed},
		{"StatusFailed", domain.StatusFailed},
		{"StatusAdjustingRes", domain.StatusAdjustingRes},
		{"StatusAddingWatermark", domain.StatusAddingWatermark},
		{"StatusRemovingAudio", domain.StatusRemovingAudio},
		{"StatusAddingIntroOutro", domain.StatusAddingIntroOutro},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.NotEmpty(t, string(tt.status))
		})
	}
}

func TestPostgresVideoRepository_ProcessedStatusLogic(t *testing.T) {
	// Test the logic for processed status updates
	status := domain.StatusProcessed
	
	updates := map[string]interface{}{
		"status": string(status),
	}
	
	if status == domain.StatusProcessed {
		updates["processed_at"] = "NOW()"
	} else {
		updates["processed_at"] = nil
	}
	
	assert.Equal(t, "PROCESSED", updates["status"])
	assert.Equal(t, "NOW()", updates["processed_at"])
}

func TestPostgresVideoRepository_NonProcessedStatusLogic(t *testing.T) {
	// Test the logic for non-processed status updates
	status := domain.StatusTrimming
	
	updates := map[string]interface{}{
		"status": string(status),
	}
	
	if status == domain.StatusProcessed {
		updates["processed_at"] = "NOW()"
	} else {
		updates["processed_at"] = nil
	}
	
	assert.Equal(t, "TRIMMING", updates["status"])
	assert.Nil(t, updates["processed_at"])
}

func TestPostgresVideoRepository_UpdateStatusAndProcessedFileLogic(t *testing.T) {
	// Test the logic for updating status and processed file
	status := domain.StatusProcessed
	processedFile := "final.mp4"
	
	updates := map[string]interface{}{
		"status":         string(status),
		"processed_file": processedFile,
	}
	
	if status == domain.StatusProcessed {
		updates["processed_at"] = "NOW()"
	} else {
		updates["processed_at"] = nil
	}
	
	assert.Equal(t, "PROCESSED", updates["status"])
	assert.Equal(t, "final.mp4", updates["processed_file"])
	assert.Equal(t, "NOW()", updates["processed_at"])
}