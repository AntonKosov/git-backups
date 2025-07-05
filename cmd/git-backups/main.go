package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/AntonKosov/git-backups/internal/clog"
	"github.com/AntonKosov/git-backups/internal/cmd"
	"github.com/AntonKosov/git-backups/internal/git"
)

func main() {
	h := clog.NewHandler(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	slog.SetDefault(slog.New(h))
	ctx := context.Background()

	cmd.Execute(ctx, "git", "status")
	git.Clone(ctx, "test.url", "this/is/path")
	// cmd.Execute(ctx, "git", "status")
}
