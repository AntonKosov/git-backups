package cmd

import (
	"context"
	"errors"
	"fmt"
	"io"
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

type Options struct {
	args         []string
	stdoutWriter io.Writer
}

type Option func(*Options)

func WithArguments(args ...string) Option {
	return func(o *Options) {
		o.args = args
	}
}

func WithStdoutWriter(writer io.Writer) Option {
	return func(o *Options) {
		o.stdoutWriter = writer
	}
}

func Execute(ctx context.Context, name string, opts ...Option) error {
	var options Options
	for _, opt := range opts {
		opt(&options)
	}

	args := options.args

	ctx = clog.Add(ctx, "name", name, "args", args)
	slog.DebugContext(ctx, "Executing application...")
	command := exec.Command(name, args...)

	var stderr strings.Builder
	command.Stderr = &stderr
	if w := options.stdoutWriter; w != nil {
		command.Stdout = w
	}

	if err := command.Run(); err != nil {
		err = errors.Join(err, CommandError{Name: name, Args: args, Err: stderr.String()})
		slog.ErrorContext(ctx, "Failed to run application", "error", err.Error())
		return err
	}

	slog.DebugContext(ctx, "Successfully ran application")
	return nil
}
