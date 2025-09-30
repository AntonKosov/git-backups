package git

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/AntonKosov/git-backups/internal/clog"
	"github.com/AntonKosov/git-backups/internal/cmd"
)

type Git struct {
}

func (g Git) Clone(ctx context.Context, url, path string, privateSSHKey *string) error {
	ctx = clog.Add(ctx, "URL", url, "path", path)
	slog.InfoContext(ctx, "Cloning repository...")

	err := cmd.Execute(
		ctx,
		"git",
		cmd.WithArguments("clone", "--bare", url, path),
		sshKeyEnvVariableCommandOption(privateSSHKey),
	)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to clone", "error", err.Error())

		return err
	}

	slog.InfoContext(ctx, "Successfully cloned repository")
	return nil
}

func (g Git) Fetch(ctx context.Context, path string, privateSSHKey *string) error {
	ctx = clog.Add(ctx, "path", path)
	slog.InfoContext(ctx, "Fetching repository...")

	err := cmd.Execute(
		ctx,
		"git",
		cmd.WithArguments("-C", path, "--bare", "fetch"),
		sshKeyEnvVariableCommandOption(privateSSHKey),
	)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to fetch", "error", err.Error())

		return err
	}

	slog.InfoContext(ctx, "Successfully fetched repository")
	return nil
}

func (g Git) GetRemoteURL(ctx context.Context, path string) (string, error) {
	var url strings.Builder
	err := cmd.Execute(
		ctx,
		"git",
		cmd.WithArguments("-C", path, "remote", "get-url", "origin"),
		cmd.WithStdoutWriter(&url),
	)

	return strings.Trim(url.String(), " \n"), err
}

func (g Git) SetRemoteURL(ctx context.Context, path, url string) error {
	err := cmd.Execute(
		ctx,
		"git",
		cmd.WithArguments("-C", path, "remote", "set-url", "origin", url),
	)

	return err
}

func sshKeyEnvVariableCommandOption(privateSSHKey *string) cmd.Option {
	if privateSSHKey == nil {
		return nil
	}

	return cmd.WithEnvVariables(
		fmt.Sprintf(`GIT_SSH_COMMAND="ssh -i %v -o IdentitiesOnly=yes"`, *privateSSHKey),
	)
}
