package git

import (
	"context"
	"errors"
	"log/slog"
	"os"

	"github.com/AntonKosov/git-backups/internal/clog"
)

//counterfeiter:generate . GitWorker
type GitWorker interface {
	Clone(ctx context.Context, url, path string) error
	Fetch(ctx context.Context, path string) error
}

type Fetcher struct {
	worker GitWorker
}

func NewFetcher(worker GitWorker) Fetcher {
	return Fetcher{worker: worker}
}

func (f Fetcher) Run(ctx context.Context, url, targetFolder string) error {
	ctx = clog.Add(ctx, "URL", url, "target folder", targetFolder)
	exists, err := folderExists(targetFolder)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to check folder", "error", err)
		return err
	}

	if exists {
		return f.worker.Fetch(ctx, targetFolder)
	}

	return f.worker.Clone(ctx, url, targetFolder)
}

func folderExists(folder string) (bool, error) {
	_, err := os.Stat(folder)
	if err == nil {
		return true, nil
	}
	var pathError *os.PathError
	if errors.As(err, &pathError) {
		return false, nil
	}

	return false, err
}
