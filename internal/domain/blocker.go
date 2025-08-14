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
