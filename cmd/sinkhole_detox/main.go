package main

import (
	"context"
	"log/slog"
	"os"
	"runtime/debug"

	_ "time/tzdata" // Load timezone data

	"github.com/alkshmir/sinkhole-detox/internal/infra/config"
	"github.com/alkshmir/sinkhole-detox/internal/presentation"
)

var srv *presentation.Server

func showVersion() {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		slog.Error("failed to read build info")
		return
	}
	slog.Debug("Build info", "info", info)
	var commit string
	for _, setting := range info.Settings {
		if setting.Key == "vcs.revision" {
			commit = setting.Value
			break
		}
	}
	slog.Info("Sinkhole-Detox", "GoVersion", info.GoVersion, "Version", info.Main.Version, "Commit", commit)
}

func init() {
	configPath := "config/config.yaml"
	if envPath := os.Getenv("CONFIG_FILE_PATH"); envPath != "" {
		configPath = envPath
	}

	conf, err := config.LoadConfig(configPath)
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
	showVersion()
	srv.Start()
}
