package git

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/AntonKosov/git-backups/internal/clog"
	"github.com/AntonKosov/git-backups/internal/cmd"
)

type Git struct {
}

func (g Git) Clone(ctx context.Context, url, path string, privateSSHKey *string) error {
	ctx = clog.Add(ctx, "path", path)
	slog.InfoContext(ctx, "Cloning repository...")

	err := cmd.Execute(
		ctx,
		"git",
		argumentsWithSSHKey(privateSSHKey, "clone", "--bare", url, path),
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
		argumentsWithSSHKey(privateSSHKey, "-C", path, "--bare", "fetch"),
	)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to fetch", "error", err.Error())

		return err
	}

	slog.InfoContext(ctx, "Successfully fetched repository")
	return nil
}

func argumentsWithSSHKey(privateSSHKey *string, otherArgs ...string) cmd.Option {
	if privateSSHKey != nil {
		sshCommand := fmt.Sprintf(
			`core.sshCommand=ssh -i "%v" -o IdentitiesOnly=yes`,
			*privateSSHKey,
		)
		otherArgs = append([]string{"-c", sshCommand}, otherArgs...)
	}

	return cmd.WithArguments(otherArgs...)
}
