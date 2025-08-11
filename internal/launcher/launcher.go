package launcher

import (
	"context"
	"log/slog"
	"path"

	"github.com/AntonKosov/git-backups/internal/clog"
	"github.com/AntonKosov/git-backups/internal/config"
)

//counterfeiter:generate . BackupService
type BackupService interface {
	Run(ctx context.Context, url, targetFolder string) error
}

func Run(ctx context.Context, conf config.Config, backupService BackupService) error {
	slog.InfoContext(ctx, "Beginning to backup generic repositories...")
	err := backupGeneric(ctx, conf.Repositories.Generic, backupService)
	if err != nil {
		return err
	}
	slog.InfoContext(ctx, "Backed up generic repositories")

	slog.InfoContext(ctx, "Beginning to backup github repositories...")
	err = backupGitHub(ctx, conf.Repositories.GitHub, backupService)
	if err != nil {
		return err
	}
	slog.InfoContext(ctx, "Backed up github repositories")

	return nil
}

func backupGeneric(ctx context.Context, genericConf []config.GenericRepo, backupService BackupService) error {
	for _, profile := range genericConf {
		ctx := clog.Add(ctx, "profile", profile)
		for _, target := range profile.Targets {
			targetPath := path.Join(profile.RootFolder, target.Folder)
			ctx := clog.Add(ctx, "URL", target.URL, "Target folder", targetPath)
			if err := backupService.Run(ctx, target.URL, targetPath); err != nil {
				slog.ErrorContext(ctx, "Failed to backup", "error", err)
				return err
			}
		}
	}

	return nil
}

func backupGitHub(context.Context, []config.GitHubRepo, BackupService) error {
	// TODO
	return nil
}
