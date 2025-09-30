package backup

import (
	"context"
	"errors"
	"log/slog"
	"os"

	"github.com/AntonKosov/git-backups/internal/clog"
)

//counterfeiter:generate . Git
type Git interface {
	Clone(ctx context.Context, url, path string, privateSSHKey *string) error
	Fetch(ctx context.Context, path string, privateSSHKey *string) error
	GetRemoteURL(ctx context.Context, path string) (string, error)
	SetRemoteURL(ctx context.Context, path, url string) error
}

type Service struct {
	git Git
}

func NewService(git Git) Service {
	return Service{git: git}
}

func (s Service) Run(ctx context.Context, url, targetFolder string, privateSSHKey *string) error {
	ctx = clog.Add(ctx, "URL", url, "target folder", targetFolder)
	exists, err := folderExists(targetFolder)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to check folder", "error", err)
		return err
	}

	if !exists {
		return s.git.Clone(ctx, url, targetFolder, privateSSHKey)
	}

	currentURL, err := s.git.GetRemoteURL(ctx, targetFolder)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to check current remote URL", "error", err)
		return err
	}

	if currentURL != url {
		err := s.git.SetRemoteURL(ctx, targetFolder, url)
		if err != nil {
			slog.ErrorContext(ctx, "Failed to replace outdated remote URL", "error", err)
			return err
		}
	}

	return s.git.Fetch(ctx, targetFolder, privateSSHKey)
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
