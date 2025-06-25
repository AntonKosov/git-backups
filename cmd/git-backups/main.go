package main

import (
	"context"

	"github.com/AntonKosov/git-backups/internal/cmd"
)

func main() {
	ctx := context.Background()
	cmd.Execute(ctx, "git", "status")
}
