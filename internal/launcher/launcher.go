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
	Run(ctx context.Context, url, targetFolder string, privateSSHKey *string) error
}

//counterfeiter:generate . ReaderService
type ReaderService interface {
	AllRepos(ctx context.Context, token, affiliation string) iter.Seq2[github.Repo, error]
}

func Run(ctx context.Context, conf config.Config, backupService BackupService, readerService ReaderService) error {
	slog.InfoContext(ctx, "Beginning to backup generic repositories...")
	err := backupGenericProfiles(ctx, conf.Profiles.GenericProfiles, backupService)
	slog.InfoContext(ctx, "Backed up generic repositories")

	slog.InfoContext(ctx, "Beginning to backup github repositories...")
	err = errors.Join(err, backupGitHubProfiles(ctx, conf.Profiles.GitHubProfiles, backupService, readerService))
	slog.InfoContext(ctx, "Backed up github repositories")

	return err
}

func backupGenericProfiles(ctx context.Context, genericProfiles []config.GenericProfile, backupService BackupService) (backupErrors error) {
	for _, profile := range genericProfiles {
		ctx := clog.Add(ctx, "profile", profile)
		for _, target := range profile.Targets {
			select {
			case <-ctx.Done():
				return errors.Join(backupErrors, context.Canceled)
			default:
				targetPath := path.Join(profile.RootFolder, target.Folder)
				ctx := clog.Add(ctx, "URL", target.URL, "Target folder", targetPath)
				if err := backupService.Run(ctx, target.URL, targetPath, profile.PrivateSSHKey); err != nil {
					slog.ErrorContext(ctx, "Failed to backup", "error", err)
					backupErrors = errors.Join(backupErrors, fmt.Errorf("failed to backup repository %v from profile %v: %w", target.URL, profile.Name, err))
				}
			}
		}
	}

	return backupErrors
}

func backupGitHubProfiles(ctx context.Context, githubProfiles []config.GitHubProfile, backupService BackupService, readerService ReaderService) (backupErrors error) {
	for _, profile := range githubProfiles {
		ctx := clog.Add(ctx, "profile", profile)
		privateSSHKey := profile.PrivateSSHKey
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
				backupErrors = errors.Join(backupErrors, errors.Join(backupErrors, fmt.Errorf("failed to read repositories: %w", err)))
				break
			}

			select {
			case <-ctx.Done():
				return errors.Join(backupErrors, context.Canceled)
			default:
				ctx := clog.Add(ctx, "repo", repo.Name)

				url := repo.GitURL
				if privateSSHKey == nil {
					var err error
					url, err = addTokenToGithubURL(repo.CloneURL, profile.Token)
					if err != nil {
						slog.ErrorContext(ctx, "Failed to add token to URL", "error", err)
						backupErrors = errors.Join(backupErrors, fmt.Errorf("failed to add token to clone URL %v from profile %v: %w", repo.CloneURL, profile.Name, err))
						break
					}
				}

				err := backupService.Run(
					ctx,
					url,
					path.Join(profile.RootFolder, repo.Owner, repo.Name),
					privateSSHKey,
				)
				if err != nil {
					slog.ErrorContext(ctx, "Failed to backup", "error", err)
					backupErrors = errors.Join(backupErrors, fmt.Errorf("failed to backup repository %v from profile %v: %w", repo.CloneURL, profile.Name, err))
				}
			}
		}
	}

	return backupErrors
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
