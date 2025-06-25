package cmd

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os/exec"
	"strings"
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
	slog.DebugContext(ctx, "Executing application", "name", name, "args", args)
	command := exec.Command(name, args...)

	var stderr strings.Builder
	command.Stderr = &stderr

	if err := command.Run(); err != nil {
		err = errors.Join(err, CommandError{Name: name, Args: args, Err: stderr.String()})
		slog.ErrorContext(ctx, "Failed to run application", "error", err.Error())
		return err
	}

	slog.DebugContext(ctx, "Successfully ran application", "name", name, "args", args)
	return nil
}
