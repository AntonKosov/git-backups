package launcher

import (
	"context"
	"errors"
	"fmt"
	"iter"
	"log/slog"
	"path"
	"strings"

	"github.com/AntonKosov/git-backups/internal/clog"
	"github.com/AntonKosov/git-backups/internal/config"
	"github.com/AntonKosov/git-backups/internal/github"
	"github.com/AntonKosov/git-backups/internal/slice"
)

//counterfeiter:generate . BackupService
type BackupService interface {
	Run(ctx context.Context, url, targetFolder string) error
}

//counterfeiter:generate . ReaderService
type ReaderService interface {
	AllRepos(ctx context.Context, token, affiliation string) iter.Seq2[github.Repo, error]
}

func Run(ctx context.Context, conf config.Config, backupService BackupService, readerService ReaderService) error {
	slog.InfoContext(ctx, "Beginning to backup generic repositories...")
	err := backupGeneric(ctx, conf.Repositories.Generic, backupService)
	if err != nil {
		return err
	}
	slog.InfoContext(ctx, "Backed up generic repositories")

	slog.InfoContext(ctx, "Beginning to backup github repositories...")
	err = backupGitHub(ctx, conf.Repositories.GitHub, backupService, readerService)
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
			select {
			case <-ctx.Done():
				return context.Canceled
			default:
				targetPath := path.Join(profile.RootFolder, target.Folder)
				ctx := clog.Add(ctx, "URL", target.URL, "Target folder", targetPath)
				if err := backupService.Run(ctx, target.URL, targetPath); err != nil {
					slog.ErrorContext(ctx, "Failed to backup", "error", err)
					return err
				}
			}
		}
	}

	return nil
}

func backupGitHub(ctx context.Context, githubConf []config.GitHubRepo, backupService BackupService, readerService ReaderService) error {
	for _, profile := range githubConf {
		ctx := clog.Add(ctx, "profile", profile)
		repos := include(
			profile.Include,
			exclude(
				profile.Exclude,
				readerService.AllRepos(ctx, profile.Token, profile.Affiliation),
			),
		)
		for repo, err := range repos {
			if err != nil {
				slog.ErrorContext(ctx, "Failed to read repositories", "error", err)
				return err
			}

			select {
			case <-ctx.Done():
				return context.Canceled
			default:
				ctx := clog.Add(ctx, "repo", repo.Name)

				urlWithToken, err := addTokenToGithubURL(repo.CloneURL, profile.Token)
				if err != nil {
					slog.ErrorContext(ctx, "Failed to add token to URL", "error", err)
					return err
				}

				if err := backupService.Run(ctx, urlWithToken, path.Join(profile.RootFolder, repo.Owner, repo.Name)); err != nil {
					slog.ErrorContext(ctx, "Failed to backup", "error", err)
					return err
				}
			}
		}
	}
	return nil
}

func include(toInclude []string, repos iter.Seq2[github.Repo, error]) iter.Seq2[github.Repo, error] {
	if toInclude == nil {
		return repos
	}

	includeMap := slice.Lookup(toInclude, func(name string) (string, bool) { return strings.ToLower(name), true })
	return func(yield func(github.Repo, error) bool) {
		for repo, err := range repos {
			if err != nil {
				yield(github.Repo{}, err)
				return
			}

			if !includeMap[strings.ToLower(repo.Name)] {
				continue
			}

			if !yield(repo, nil) {
				return
			}
		}
	}
}

func exclude(toExclude []string, repos iter.Seq2[github.Repo, error]) iter.Seq2[github.Repo, error] {
	if len(toExclude) == 0 {
		return repos
	}

	excludeMap := slice.Lookup(toExclude, func(name string) (string, bool) { return strings.ToLower(name), true })
	return func(yield func(github.Repo, error) bool) {
		for repo, err := range repos {
			if err != nil {
				yield(github.Repo{}, err)
				return
			}

			if excludeMap[strings.ToLower(repo.Name)] {
				continue
			}

			if !yield(repo, nil) {
				return
			}
		}
	}
}

func addTokenToGithubURL(url, token string) (string, error) {
	if !strings.HasPrefix(url, "https://") {
		return "", errors.New("unexpected URL prefix")
	}

	return fmt.Sprintf("https://oauth2:%v@%v", token, url[8:]), nil
}
