package git

import (
	"context"
	"log/slog"
	"strings"

	"github.com/AntonKosov/git-backups/internal/clog"
	"github.com/AntonKosov/git-backups/internal/cmd"
)

type Git struct {
}

func (g Git) Clone(ctx context.Context, url, path string) error {
	ctx = clog.Add(ctx, "URL", url, "path", path)
	slog.InfoContext(ctx, "Cloning repository...")

	if _, err := cmd.Execute(ctx, false, "git", "clone", "--bare", url, path); err != nil {
		slog.ErrorContext(ctx, "Failed to clone", "error", err.Error())

		return err
	}

	slog.InfoContext(ctx, "Successfully cloned repository")
	return nil
}

func (g Git) Fetch(ctx context.Context, path string) error {
	ctx = clog.Add(ctx, "path", path)
	slog.InfoContext(ctx, "Fetching repository...")

	if _, err := cmd.Execute(ctx, false, "git", "-C", path, "--bare", "fetch"); err != nil {
		slog.ErrorContext(ctx, "Failed to fetch", "error", err.Error())

		return err
	}

	slog.InfoContext(ctx, "Successfully fetched repository")
	return nil
}

func (g Git) GetRemoteURL(ctx context.Context, path string) (string, error) {
	url, err := cmd.Execute(ctx, true, "git", "-C", path, "remote", "get-url", "origin")
	return strings.Trim(url, " \n"), err
}

func (g Git) SetRemoteURL(ctx context.Context, path, url string) error {
	_, err := cmd.Execute(ctx, false, "git", "-C", path, "remote", "set-url", "origin", url)
	return err
}
