package config

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/alkshmir/sinkhole-detox/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestBlockerFactory_GenBlockers(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		configs     []Blocker
		expected    []domain.Blocker
		expectError bool
	}{
		{
			name: "should generate blockers from valid configs",
			configs: []Blocker{
				{
					Name:   "twitter",
					Domain: "twitter.com",
					Rules: []Rule{
						{
							Type:  "everyday",
							Ops:   "block",
							Start: "00:00",
							End:   "05:00",
						},
						{
							Type:  "everyday",
							Ops:   "block",
							Start: "22:00",
							End:   "23:59",
						},
						{
							Type:     "weekday",
							Ops:      "block",
							Start:    "08:00",
							End:      "18:00",
							Weekdays: []int{1, 2, 3, 4, 5},
						},
					},
				},
				{
					Name:   "x",
					Domain: "x.com",
					Rules: []Rule{
						{
							Type:  "everyday",
							Ops:   "block",
							Start: "00:00",
							End:   "05:00",
						},
						{
							Type:     "weekday",
							Ops:      "block",
							Start:    "10:00",
							End:      "16:00",
							Weekdays: []int{1, 2, 3, 4, 5},
						},
					},
				},
			},
			expected: []domain.Blocker{
				{
					Domain:    "twitter.com",
					ForwardTo: net.IPv4(0, 0, 0, 0),
					Rules: []domain.BlockRule{
						domain.EveryDayRule{
							Op:   domain.BlockOpsBlock,
							From: time.Date(0, 1, 1, 0, 0, 0, 0, time.UTC),
							To:   time.Date(0, 1, 1, 5, 0, 0, 0, time.UTC),
						},
						domain.EveryDayRule{
							Op:   domain.BlockOpsBlock,
							From: time.Date(0, 1, 1, 22, 0, 0, 0, time.UTC),
							To:   time.Date(0, 1, 1, 23, 59, 0, 0, time.UTC),
						},
						domain.WeekdayRule{
							Op:       domain.BlockOpsBlock,
							From:     time.Date(0, 1, 1, 8, 0, 0, 0, time.UTC),
							To:       time.Date(0, 1, 1, 18, 0, 0, 0, time.UTC),
							Weekdays: []time.Weekday{time.Monday, time.Tuesday, time.Wednesday, time.Thursday, time.Friday},
						},
					},
				},
				{
					Domain:    "x.com",
					ForwardTo: net.IPv4(0, 0, 0, 0),
					Rules: []domain.BlockRule{
						domain.EveryDayRule{
							Op:   domain.BlockOpsBlock,
							From: time.Date(0, 1, 1, 0, 0, 0, 0, time.UTC),
							To:   time.Date(0, 1, 1, 5, 0, 0, 0, time.UTC),
						},
						domain.WeekdayRule{
							Op:       domain.BlockOpsBlock,
							From:     time.Date(0, 1, 1, 10, 0, 0, 0, time.UTC),
							To:       time.Date(0, 1, 1, 16, 0, 0, 0, time.UTC),
							Weekdays: []time.Weekday{time.Monday, time.Tuesday, time.Wednesday, time.Thursday, time.Friday},
						},
					},
				},
			},
			expectError: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			factory := BlockerFactory{}
			blockers, err := factory.GenBlockers(context.Background(), tt.configs)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, blockers, len(tt.expected))
				assert.Equal(t, tt.expected, blockers)
			}
		})
	}
}
