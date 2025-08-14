package config

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"os"
	"time"

	"github.com/alkshmir/sinkhole-detox.git/internal/domain"
	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig `mapstructure:"server"`
	Blockers []Blocker    `mapstructure:"blockers"`
}

type ServerConfig struct {
	Port int `mapstructure:"port"`
}

type Blocker struct {
	Name      string `mapstructure:"name"`
	Domain    string `mapstructure:"domain"`
	ForwardTo string `mapstructure:"forward_to"` // IP address to forward the request to this domain
	Rules     []Rule `mapstructure:"rules"`
}

type Rule struct {
	Type     string `mapstructure:"type"`     // "everyday" / "weekday"
	Ops      string `mapstructure:"ops"`      // "block" / "allow"
	Start    string `mapstructure:"start"`    // HH:MM
	End      string `mapstructure:"end"`      // HH:MM
	Weekdays []int  `mapstructure:"weekdays"` // 0=Sunday, 1=Monday, ..., 6=Saturday
}

func LoadConfig() (*Config, error) {
	configPath := "config/config.yaml"
	if envPath := os.Getenv("CONFIG_FILE_PATH"); envPath != "" {
		configPath = envPath
	}
	viper.SetConfigFile(configPath)

	slog.Info("Loading configuration from file", "path", configPath)
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}
	slog.Info("Configuration loaded successfully", "config", cfg)

	return &cfg, nil
}

func (b *Blocker) ToBlocker(ctx context.Context) (domain.Blocker, error) {
	forwardTo := net.ParseIP(b.ForwardTo)
	if forwardTo == nil {
		slog.Info("ForwardTo IP is invalid or not set, defaulting to 0.0.0.0", "forwardTo", b.ForwardTo)
		forwardTo = net.IPv4(0, 0, 0, 0) // Defaul
	}

	rules := make([]domain.BlockRule, len(b.Rules))
	for i, r := range b.Rules {
		start, err := parseTime(r.Start)
		if err != nil {
			return domain.Blocker{}, fmt.Errorf("failed to parse start time %s: %w", r.Start, err)
		}
		end, err := parseTime(r.End)
		if err != nil {
			return domain.Blocker{}, fmt.Errorf("failed to parse end time %s: %w", r.End, err)
		}

		var rule domain.BlockRule
		// TODO: rewrite to abstract factory pattern
		switch r.Type {
		case "everyday":
			rule, err = domain.NewEveryDayRule(
				r.Ops,
				start,
				end,
			)
			if err != nil {
				return domain.Blocker{}, fmt.Errorf("failed to create everyday rule: %w", err)
			}

		case "weekday":
			weekdays, err := parseWeekdays(r.Weekdays)
			if err != nil {
				return domain.Blocker{}, fmt.Errorf("failed to parse weekdays: %w", err)
			}
			rule, err = domain.NewWeekdayRule(
				r.Ops,
				start,
				end,
				weekdays,
			)
			if err != nil {
				return domain.Blocker{}, fmt.Errorf("failed to create weekday rule: %w", err)
			}
		default:
			return domain.Blocker{}, fmt.Errorf("unknown rule type: %s", r.Type)
		}
		rules[i] = rule
	}

	return domain.Blocker{
		Domain:    b.Domain,
		ForwardTo: forwardTo,
		Rules:     rules,
	}, nil
}

func parseTime(s string) (time.Time, error) {
	t, err := time.Parse("15:04", s)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to parse time %s: %w", s, err)
	}
	return t, nil
}

func parseWeekdays(weekdays []int) ([]time.Weekday, error) {
	var parsed []time.Weekday
	for _, wd := range weekdays {
		if wd < 0 || wd > 6 {
			return nil, fmt.Errorf("invalid weekday: %d", wd)
		}
		parsed = append(parsed, time.Weekday(wd))
	}
	return parsed, nil
}
