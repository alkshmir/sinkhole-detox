package config

import (
	"context"

	"github.com/alkshmir/sinkhole-detox.git/internal/domain"
)

type BlockerFactory struct {
}

func (f *BlockerFactory) GenBlockers(ctx context.Context, configs []Blocker) ([]domain.Blocker, error) {
	var blockers []domain.Blocker
	for _, config := range configs {
		blocker, err := config.ToBlocker(ctx)
		if err != nil {
			return nil, err
		}
		blockers = append(blockers, blocker)
	}
	return blockers, nil
}
