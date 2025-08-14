package domain

import (
	"fmt"
	"log/slog"
	"time"
)

type EveryDayRule struct {
	Op   BlockOps
	From time.Time // inclusive
	To   time.Time // exclusive
}

func NewEveryDayRule(ops string, from, to time.Time) (BlockRule, error) {
	r := EveryDayRule{
		Op:   BlockOps(ops),
		From: from,
		To:   to,
	}
	if err := r.Validate(); err != nil {
		return nil, fmt.Errorf("failed to create everyday rule: %w", err)
	}
	return r, nil
}

func (s EveryDayRule) Validate() error {
	if err := s.Op.Validate(); err != nil {
		return fmt.Errorf("invalid ops: %w", err)
	}
	if s.From.After(s.To) {
		return fmt.Errorf("from time must be before to time")
	}
	return nil

}

func (s EveryDayRule) Ops() BlockOps {
	return s.Op
}

func (s EveryDayRule) IsActive(t time.Time) bool {
	from := time.Date(t.Year(), t.Month(), t.Day(), s.From.Hour(), s.From.Minute(), 0, 0, t.Location())
	to := time.Date(t.Year(), t.Month(), t.Day(), s.To.Hour(), s.To.Minute(), 0, 0, t.Location())
	slog.Info("from/to", "from", from, "to", to, "t", t)
	if from.After(t) || to.Before(t) {
		return false
	}
	return true
}
