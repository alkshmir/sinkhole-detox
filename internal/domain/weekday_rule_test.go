package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestWeekdayRule_Validate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		rule      WeekdayRule
		expectErr bool
	}{
		{
			name: "should validate correct rule",
			rule: WeekdayRule{
				Op:       BlockOpsBlock,
				From:     time.Date(0, 1, 1, 8, 0, 0, 0, time.UTC),
				To:       time.Date(0, 1, 1, 18, 0, 0, 0, time.UTC),
				Weekdays: []time.Weekday{time.Monday, time.Tuesday, time.Wednesday, time.Thursday, time.Friday},
			},
			expectErr: false,
		},
		{
			name: "should return error when from time is after to time",
			rule: WeekdayRule{
				Op:       BlockOpsBlock,
				From:     time.Date(0, 1, 1, 23, 0, 0, 0, time.UTC),
				To:       time.Date(0, 1, 1, 22, 59, 59, 999999999, time.UTC),
				Weekdays: []time.Weekday{time.Monday, time.Tuesday, time.Wednesday, time.Thursday, time.Friday},
			},
			expectErr: true,
		},
		{
			name: "should return error when no weekdays are specified",
			rule: WeekdayRule{
				Op:       BlockOpsBlock,
				From:     time.Date(0, 1, 1, 8, 0, 0, 0, time.UTC),
				To:       time.Date(0, 1, 1, 18, 0, 0, 0, time.UTC),
				Weekdays: []time.Weekday{},
			},
			expectErr: true,
		},
		{
			name: "should return error for invalid weekday",
			rule: WeekdayRule{
				Op:       BlockOpsBlock,
				From:     time.Date(0, 1, 1, 8, 0, 0, 0, time.UTC),
				To:       time.Date(0, 1, 1, 18, 0, 0, 0, time.UTC),
				Weekdays: []time.Weekday{time.Sunday, 7}, // 7 is invalid
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.rule.Validate()
			if tt.expectErr {
				assert.Error(t, err, "expected error but got none")
			} else {
				assert.NoError(t, err, "expected no error but got one")
			}
		})
	}
}

func TestWeekdayRule_IsActive(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		rule     WeekdayRule
		time     time.Time
		expected bool
	}{
		{
			name: "should return true if time is within the range on a valid weekday",
			rule: WeekdayRule{
				Op:       BlockOpsBlock,
				From:     time.Date(0, 0, 0, 8, 0, 0, 0, time.UTC),
				To:       time.Date(0, 0, 0, 18, 0, 0, 0, time.UTC),
				Weekdays: []time.Weekday{time.Wednesday},
			},
			time:     time.Date(2025, 1, 1, 10, 30, 0, 0, time.UTC), // Wednesday
			expected: true,
		},
		{
			name: "should return false if time is outside the range on a valid weekday",
			rule: WeekdayRule{
				Op:       BlockOpsBlock,
				From:     time.Date(0, 0, 0, 8, 0, 0, 0, time.UTC),
				To:       time.Date(0, 0, 0, 18, 0, 0, 0, time.UTC),
				Weekdays: []time.Weekday{time.Monday},
			},
			time:     time.Date(2025, 1, 6, 19, 30, 0, 0, time.UTC), // Monday
			expected: false,
		},
		{
			name: "should return false if weekday is not in the specified weekdays",
			rule: WeekdayRule{
				Op:       BlockOpsBlock,
				From:     time.Date(0, 1, 1, 8, 0, 0, 0, time.UTC),
				To:       time.Date(0, 1, 1, 18, 0, 0, 0, time.UTC),
				Weekdays: []time.Weekday{time.Monday},
			},
			time:     time.Date(2025, 1, 1, 10, 30, 0, 0, time.UTC), // Wednesday
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.rule.IsActive(tt.time)
			assert.Equal(t, tt.expected, result)
		})
	}
}
