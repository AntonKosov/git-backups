package cmd

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os/exec"
	"strings"

	"github.com/AntonKosov/git-backups/internal/clog"
)

type CommandError struct {
	Name string
	Args []string
	Err  string
}

func (ce CommandError) Error() string {
	return fmt.Sprintf(`%v failed with "%v" (args: %v)`, ce.Name, ce.Err, ce.Args)
}

func Execute(ctx context.Context, name string, args ...string) error {
	ctx = clog.Add(ctx, "name", name, "args", args)
	slog.DebugContext(ctx, "Executing application")
	command := exec.Command(name, args...)

	var stderr strings.Builder
	command.Stderr = &stderr

	if err := command.Run(); err != nil {
		err = errors.Join(err, CommandError{Name: name, Args: args, Err: stderr.String()})
		slog.ErrorContext(ctx, "Failed to run application", "error", err.Error())
		return err
	}

	slog.DebugContext(ctx, "Successfully ran application")
	return nil
}
