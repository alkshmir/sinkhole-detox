package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/alkshmir/sinkhole-detox.git/internal/infra/config"
	"github.com/alkshmir/sinkhole-detox.git/internal/presentation"
)

var srv *presentation.Server

func init() {
	conf, err := config.LoadConfig()
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}
	slog.Debug("Configuration loaded", "config", conf)

	f := config.BlockerFactory{}
	blockers, err := f.GenBlockers(context.Background(), conf.Blockers)
	if err != nil {
		slog.Error("failed to create blockers from config", "error", err)
		os.Exit(1)
	}
	slog.Debug("Blockers created from config", "blockers", blockers)

	srv = presentation.NewServer(blockers, presentation.ServerConfig{
		Port: uint(conf.Server.Port),
	})
}

func main() {
	srv.Start()
}
