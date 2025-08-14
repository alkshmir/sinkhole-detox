package domain

import (
	"fmt"
	"log/slog"
	"net"
	"time"
)

type Blocker struct {
	Domain string // RFC 1035
	// ForwardTo is IP address to forward requests to if the domain is blocked.
	// Usually this is 0.0.0.0
	ForwardTo net.IP
	// Rules is a list of blocking Rules. Latter Rules take precedence over earlier ones.
	Rules []BlockRule
}

func (b *Blocker) IsBlocked(t time.Time) bool {
	slog.Info("evaluating blocker for domain", "domain", b.Domain, "time", t)
	blocked := false
	for _, rule := range b.Rules {
		if rule.IsActive(t) {
			switch rule.Ops() {
			case BlockOpsBlock:
				blocked = true
				// case BlockOpsAllow:
				// 	blocked = false
			}
		}
	}
	slog.Info("blocker evaluation result", "domain", b.Domain, "blocked", blocked)
	return blocked
}

type BlockRule interface {
	Ops() BlockOps
	IsActive(time.Time) bool
}

type EveryDayRule struct {
	ops  BlockOps
	from time.Time // inclusive
	to   time.Time // exclusive
}

func NewEveryDayRule(ops string, from, to time.Time) (BlockRule, error) {
	r := EveryDayRule{
		ops:  BlockOps(ops),
		from: from,
		to:   to,
	}
	if err := r.Validate(); err != nil {
		return nil, fmt.Errorf("failed to create everyday rule: %w", err)
	}
	return r, nil
}

func (s EveryDayRule) Validate() error {
	if err := s.ops.Validate(); err != nil {
		return fmt.Errorf("invalid ops: %w", err)
	}
	if s.from.After(s.to) {
		return fmt.Errorf("from time must be before to time")
	}
	return nil

}

func (s EveryDayRule) Ops() BlockOps {
	return s.ops
}

func (s EveryDayRule) IsActive(t time.Time) bool {
	from := time.Date(t.Year(), t.Month(), t.Day(), s.from.Hour(), s.from.Minute(), 0, 0, t.Location())
	to := time.Date(t.Year(), t.Month(), t.Day(), s.to.Hour(), s.to.Minute(), 0, 0, t.Location())
	slog.Info("from/to", "from", from, "to", to, "t", t)
	if from.After(t) || to.Before(t) {
		return false
	}
	return true
}

type WeekdayRule struct {
	ops  BlockOps
	from time.Time // inclusive
	to   time.Time // exclusive
	// weekdays is a list of days of the week when this rule is active.
	weekdays []time.Weekday
}

func NewWeekdayRule(ops string, from, to time.Time, weekdays []time.Weekday) (BlockRule, error) {
	r := WeekdayRule{
		ops:      BlockOps(ops),
		from:     from,
		to:       to,
		weekdays: weekdays,
	}
	if err := r.Validate(); err != nil {
		return nil, fmt.Errorf("failed to create weekday rule: %w", err)
	}
	return r, nil
}

func (s WeekdayRule) Validate() error {
	if err := s.ops.Validate(); err != nil {
		return fmt.Errorf("invalid ops: %w", err)
	}
	if s.from.After(s.to) {
		return fmt.Errorf("from time must be before to time")
	}
	if len(s.weekdays) == 0 {
		return fmt.Errorf("weekdays cannot be empty")
	}
	for _, day := range s.weekdays {
		if day < time.Sunday || day > time.Saturday {
			return fmt.Errorf("invalid weekday: %v", day)
		}
	}
	return nil
}

func (s WeekdayRule) Ops() BlockOps {
	return s.ops
}

func (s WeekdayRule) IsActive(t time.Time) bool {
	from := time.Date(t.Year(), t.Month(), t.Day(), s.from.Hour(), s.from.Minute(), 0, 0, t.Location())
	to := time.Date(t.Year(), t.Month(), t.Day(), s.to.Hour(), s.to.Minute(), 0, 0, t.Location())

	w := t.Weekday()
	for _, day := range s.weekdays {
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

type BlockOps string

const (
	BlockOpsBlock BlockOps = "block"
	// BlockOpsAllow BlockOps = "allow"
)

func (o BlockOps) Validate() error {
	switch o {
	case BlockOpsBlock:
		return nil
	// case BlockOpsAllow:
	// 	return nil
	default:
		return fmt.Errorf("unknown block operation: %s", o)
	}
}
