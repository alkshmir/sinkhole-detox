package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type MockRule struct {
	Active bool
}

var _ BlockRule = &MockRule{}

func (m *MockRule) Ops() BlockOps {
	return BlockOpsBlock
}

func (m *MockRule) IsActive(t time.Time) bool {
	return m.Active
}

func TestBlocker_IsBlocked(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		blocker  Blocker
		time     time.Time
		expected bool
	}{
		{
			name: "should return true if any rule is active",
			blocker: Blocker{
				Domain: "example.com",
				Rules: []BlockRule{
					&MockRule{Active: true},
					&MockRule{Active: false},
				},
			},
			time:     time.Now(),
			expected: true,
		},
		{
			name: "should return false if no rules are active",
			blocker: Blocker{
				Domain: "example.com",
				Rules: []BlockRule{
					&MockRule{Active: false},
					&MockRule{Active: false},
				},
			},
			time:     time.Now(),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.blocker.IsBlocked(tt.time)
			assert.Equal(t, tt.expected, result, "expected %v but got %v", tt.expected, result)
		})
	}
}

func TestBlockOps_Validate(t *testing.T) {
	tests := []struct {
		name    string
		ops     BlockOps
		wantErr bool
	}{
		{
			name: "should return no error for valid ops",
			ops:  BlockOpsBlock,
		},
		{
			name:    "should return error for invalid ops",
			ops:     "invalid",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.ops.Validate()
			if tt.wantErr {
				assert.Error(t, err, "expected error but got none")
			} else {
				assert.NoError(t, err, "expected no error but got one")
			}
		})
	}
}
