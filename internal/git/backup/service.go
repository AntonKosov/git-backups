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
}

type Service struct {
	git Git
}

func NewService(git Git) Service {
	return Service{git: git}
}

func (s Service) Run(ctx context.Context, url, targetFolder string, privateSSHKey *string) error {
	ctx = clog.Add(ctx, "target folder", targetFolder)
	exists, err := folderExists(targetFolder)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to check folder", "error", err)
		return err
	}

	if exists {
		return s.git.Fetch(ctx, targetFolder, privateSSHKey)
	}

	return s.git.Clone(ctx, url, targetFolder, privateSSHKey)
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
