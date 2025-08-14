package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestEveryDayRule_Validate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		rule      EveryDayRule
		expectErr bool
	}{
		{
			name: "should validate correct rule",
			rule: EveryDayRule{
				Op:   BlockOpsBlock,
				From: time.Date(0, 1, 1, 22, 0, 0, 0, time.UTC),
				To:   time.Date(0, 1, 1, 23, 59, 59, 999999999, time.UTC),
			},
			expectErr: false,
		},
		{
			name: "should return error when from time is after to time",
			rule: EveryDayRule{
				Op:   BlockOpsBlock,
				From: time.Date(0, 1, 1, 23, 0, 0, 0, time.UTC),
				To:   time.Date(0, 1, 1, 22, 59, 59, 999999999, time.UTC),
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.expectErr {
				assert.Error(t, tt.rule.Validate(), "expected error but got none")
			} else {
				assert.NoError(t, tt.rule.Validate(), "expected no error but got one")
			}
		})
	}
}

func TestEveryDayRule_IsActive(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		rule     EveryDayRule
		time     time.Time
		expected bool
	}{
		{
			name: "should return true if time is within the range",
			rule: EveryDayRule{
				Op:   BlockOpsBlock,
				From: time.Date(0, 0, 0, 22, 0, 0, 0, time.UTC),
				To:   time.Date(0, 0, 0, 23, 59, 59, 999999999, time.UTC),
			},
			time:     time.Date(2025, 1, 1, 22, 30, 0, 0, time.UTC),
			expected: true,
		},
		{
			name: "should return false if time is outside the range",
			rule: EveryDayRule{
				Op:   BlockOpsBlock,
				From: time.Date(0, 0, 0, 22, 0, 0, 0, time.UTC),
				To:   time.Date(0, 0, 0, 23, 59, 59, 999999999, time.UTC),
			},
			time:     time.Date(2025, 1, 1, 21, 30, 0, 0, time.UTC),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.rule.IsActive(tt.time)
			assert.Equal(t, tt.expected, result, "expected %v but got %v", tt.expected, result)
		})
	}
}
