package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/AntonKosov/git-backups/internal/clog"
	"github.com/AntonKosov/git-backups/internal/cmd"
	"github.com/AntonKosov/git-backups/internal/config"
	yaml "github.com/goccy/go-yaml"
)

func main() {
	h := clog.NewHandler(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	slog.SetDefault(slog.New(h))
	ctx := context.Background()

	conf, err := readConfig()
	if err != nil {
		slog.ErrorContext(ctx, "Failed to read config", "error", err)
		os.Exit(1)
	}
	fmt.Printf("%+v\n", conf)

	cmd.Execute(ctx, "git", "status")
}

func readConfig() (config.V1, error) {
	configFile, err := os.ReadFile("config.yaml")
	if err != nil {
		return config.V1{}, err
	}

	var conf config.V1
	if err := yaml.Unmarshal(configFile, &conf); err != nil {
		return config.V1{}, err
	}

	return conf, nil
}
