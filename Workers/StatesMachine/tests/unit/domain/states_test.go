package domain_test

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

type VideoState string

const (
	StateUploaded    VideoState = "UPLOADED"
	StateTrimming    VideoState = "TRIMMING"
	StateEditing     VideoState = "EDITING"
	StateProcessing  VideoState = "PROCESSING"
	StateCompleted   VideoState = "COMPLETED"
	StateFailed      VideoState = "FAILED"
)

type StateTransition struct {
	From VideoState
	To   VideoState
}

func TestVideoState_Constants(t *testing.T) {
	assert.Equal(t, "UPLOADED", string(StateUploaded))
	assert.Equal(t, "TRIMMING", string(StateTrimming))
	assert.Equal(t, "EDITING", string(StateEditing))
	assert.Equal(t, "PROCESSING", string(StateProcessing))
	assert.Equal(t, "COMPLETED", string(StateCompleted))
	assert.Equal(t, "FAILED", string(StateFailed))
}

func TestStateTransition_Creation(t *testing.T) {
	transition := StateTransition{
		From: StateUploaded,
		To:   StateTrimming,
	}
	
	assert.Equal(t, StateUploaded, transition.From)
	assert.Equal(t, StateTrimming, transition.To)
}

func TestStateTransition_Validation(t *testing.T) {
	validTransitions := map[VideoState][]VideoState{
		StateUploaded:   {StateTrimming, StateFailed},
		StateTrimming:   {StateEditing, StateFailed},
		StateEditing:    {StateProcessing, StateFailed},
		StateProcessing: {StateCompleted, StateFailed},
	}
	
	tests := []struct {
		name       string
		transition StateTransition
		valid      bool
	}{
		{
			name:       "valid transition",
			transition: StateTransition{From: StateUploaded, To: StateTrimming},
			valid:      true,
		},
		{
			name:       "invalid transition",
			transition: StateTransition{From: StateCompleted, To: StateUploaded},
			valid:      false,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validTargets, exists := validTransitions[tt.transition.From]
			if !exists {
				assert.False(t, tt.valid)
				return
			}
			
			valid := false
			for _, target := range validTargets {
				if target == tt.transition.To {
					valid = true
					break
				}
			}
			assert.Equal(t, tt.valid, valid)
		})
	}
}