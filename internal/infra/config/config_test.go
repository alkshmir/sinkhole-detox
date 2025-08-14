package config

import (
	"context"
	"net"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/alkshmir/sinkhole-detox.git/internal/domain"
	"github.com/stretchr/testify/assert"
)

func getTestFilePath(filename string) string {
	_, file, _, _ := runtime.Caller(0)
	dir := filepath.Dir(file)
	return filepath.Join(dir, filename)
}

func TestLoadConfig(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		path           string
		expectedConfig *Config
		expectError    bool
	}{
		{
			name:        "should return config with valid path",
			path:        getTestFilePath("test.yaml"),
			expectError: false,
			expectedConfig: &Config{
				Server: ServerConfig{
					Port: 8080,
				},
				Blockers: []Blocker{
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
			},
		},
		{
			name:        "should return error for invalid path",
			path:        "./invalid.yaml",
			expectError: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			config, err := LoadConfig(tt.path)
			if tt.expectError {
				assert.Error(t, err, "expected error but got none")
			} else {
				assert.NoError(t, err, "expected no error but got one")
				assert.Equal(t, tt.expectedConfig, config)
			}
		})
	}
}

func TestBlocker_ToBlocker(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		blocker     Blocker
		expectError bool
		expected    domain.Blocker
	}{
		{
			name: "should convert Blocker to domain.Blocker successfully",
			blocker: Blocker{
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
			expectError: false,
			expected: domain.Blocker{
				Domain: "twitter.com",
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
				ForwardTo: net.IPv4(0, 0, 0, 0), // Default forward IP
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			blocker, err := tt.blocker.ToBlocker(context.Background())
			if tt.expectError {
				assert.Error(t, err, "expected error but got none")
			} else {
				assert.NoError(t, err, "expected no error but got one")
				assert.Equal(t, tt.expected, blocker)
			}
		})
	}
}
