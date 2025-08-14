package domain

import (
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func exampleBlocker() []Blocker {
	return []Blocker{
		{
			Domain:    "twitter.com",
			ForwardTo: net.IPv4(0, 0, 0, 0),
			Rules: []BlockRule{
				EveryDayRule{
					Op:   BlockOpsBlock,
					From: time.Date(0, 1, 1, 0, 0, 0, 0, time.UTC),
					To:   time.Date(0, 1, 1, 5, 0, 0, 0, time.UTC),
				},
				EveryDayRule{
					Op:   BlockOpsBlock,
					From: time.Date(0, 1, 1, 22, 0, 0, 0, time.UTC),
					To:   time.Date(0, 1, 1, 23, 59, 0, 0, time.UTC),
				},
				WeekdayRule{
					Op:       BlockOpsBlock,
					From:     time.Date(0, 1, 1, 8, 0, 0, 0, time.UTC),
					To:       time.Date(0, 1, 1, 18, 0, 0, 0, time.UTC),
					Weekdays: []time.Weekday{time.Monday, time.Tuesday, time.Wednesday, time.Thursday, time.Friday},
				},
			},
		},
		{
			Domain:    "x.com",
			ForwardTo: net.IPv4(0, 0, 0, 0),
			Rules: []BlockRule{
				EveryDayRule{
					Op:   BlockOpsBlock,
					From: time.Date(0, 1, 1, 0, 0, 0, 0, time.UTC),
					To:   time.Date(0, 1, 1, 5, 0, 0, 0, time.UTC),
				},
				WeekdayRule{
					Op:       BlockOpsBlock,
					From:     time.Date(0, 1, 1, 10, 0, 0, 0, time.UTC),
					To:       time.Date(0, 1, 1, 16, 0, 0, 0, time.UTC),
					Weekdays: []time.Weekday{time.Monday, time.Tuesday, time.Wednesday, time.Thursday, time.Friday},
				},
			},
		},
	}
}

func TestHostsGenerator_Gen(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		blockers []Blocker
		time     time.Time
		expected []HostsEntry
	}{
		{
			name:     "should generate hosts entries for active blockers",
			blockers: exampleBlocker(),
			time:     time.Date(2023, 10, 1, 2, 0, 0, 0, time.UTC),
			expected: []HostsEntry{
				{
					Domain: "twitter.com",
					IP:     net.IPv4(0, 0, 0, 0),
				},
				{
					Domain: "x.com",
					IP:     net.IPv4(0, 0, 0, 0),
				},
			},
		},
		{
			name:     "should not generate entries for inactive blockers",
			blockers: exampleBlocker(),
			time:     time.Date(2023, 10, 1, 7, 0, 0, 0, time.UTC),
			expected: []HostsEntry{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			generator := NewHostsGenerator(tt.blockers)
			entries := generator.Gen(tt.time)

			assert.Equal(t, len(tt.expected), len(entries), "expected %d entries but got %d", len(tt.expected), len(entries))

			assert.ElementsMatch(t, tt.expected, entries)
		})
	}
}
