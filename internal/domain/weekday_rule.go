package domain

import (
	"fmt"
	"time"
)

type WeekdayRule struct {
	Op   BlockOps
	From time.Time // inclusive
	To   time.Time // exclusive
	// Weekdays is a list of days of the week when this rule is active.
	Weekdays []time.Weekday
}

func NewWeekdayRule(ops string, from, to time.Time, weekdays []time.Weekday) (BlockRule, error) {
	r := WeekdayRule{
		Op:       BlockOps(ops),
		From:     from,
		To:       to,
		Weekdays: weekdays,
	}
	if err := r.Validate(); err != nil {
		return nil, fmt.Errorf("failed to create weekday rule: %w", err)
	}
	return r, nil
}

func (s WeekdayRule) Validate() error {
	if err := s.Op.Validate(); err != nil {
		return fmt.Errorf("invalid ops: %w", err)
	}
	if s.From.After(s.To) {
		return fmt.Errorf("from time must be before to time")
	}
	if len(s.Weekdays) == 0 {
		return fmt.Errorf("weekdays cannot be empty")
	}
	for _, day := range s.Weekdays {
		if day < time.Sunday || day > time.Saturday {
			return fmt.Errorf("invalid weekday: %v", day)
		}
	}
	return nil
}

func (s WeekdayRule) Ops() BlockOps {
	return s.Op
}

func (s WeekdayRule) IsActive(t time.Time) bool {
	from := time.Date(t.Year(), t.Month(), t.Day(), s.From.Hour(), s.From.Minute(), 0, 0, t.Location())
	to := time.Date(t.Year(), t.Month(), t.Day(), s.To.Hour(), s.To.Minute(), 0, 0, t.Location())

	w := t.Weekday()
	for _, day := range s.Weekdays {
		if day != w {
			continue
		}
		if from.After(t) || to.Before(t) {
			return false
		}
		return true
	}
	return false
}
