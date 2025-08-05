package git

import (
	"context"
	"log/slog"

	"github.com/AntonKosov/git-backups/internal/clog"
	"github.com/AntonKosov/git-backups/internal/cmd"
)

type Git struct {
}

func (g Git) Clone(ctx context.Context, url, path string) error {
	ctx = clog.Add(ctx, "URL", url, "path", path)
	slog.InfoContext(ctx, "Cloning repository...")

	if err := cmd.Execute(ctx, "git", "clone", "--bare", url, path); err != nil {
		slog.ErrorContext(ctx, "Failed to clone", "error", err.Error())

		return err
	}

	slog.InfoContext(ctx, "Successfully cloned repository")
	return nil
}

func (g Git) Fetch(ctx context.Context, path string) error {
	ctx = clog.Add(ctx, "path", path)
	slog.InfoContext(ctx, "Fetching repository...")

	if err := cmd.Execute(ctx, "git", "-C", path, "--bare", "fetch"); err != nil {
		slog.ErrorContext(ctx, "Failed to fetch", "error", err.Error())

		return err
	}

	slog.InfoContext(ctx, "Successfully fetched repository")
	return nil
}
