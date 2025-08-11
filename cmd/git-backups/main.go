package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/AntonKosov/git-backups/internal/clog"
	"github.com/AntonKosov/git-backups/internal/config"
	"github.com/AntonKosov/git-backups/internal/git"
	"github.com/AntonKosov/git-backups/internal/git/backup"
	"github.com/AntonKosov/git-backups/internal/launcher"
)

func main() {
	h := clog.NewHandler(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	slog.SetDefault(slog.New(h))
	ctx := context.Background()

	conf, err := config.ReadConfig("config.yaml")
	if err != nil {
		slog.ErrorContext(ctx, "Failed to read config", "error", err)
		os.Exit(1)
	}

	err = launcher.Run(ctx, conf, backup.NewService(git.Git{}))
	if err != nil {
		slog.ErrorContext(ctx, "Failed to backup", "error", err)
		os.Exit(1)
	}
}
