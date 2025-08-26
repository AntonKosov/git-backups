package launcher_test

import (
	"context"
	"io"
	"log/slog"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var (
	ctx       context.Context
	ctxCancel context.CancelFunc
)

func TestLauncher(t *testing.T) {
	slog.SetDefault(slog.New(slog.NewTextHandler(io.Discard, nil)))
	RegisterFailHandler(Fail)
	RunSpecs(t, "Launcher Suite")
}

var _ = BeforeEach(func() {
	ctx, ctxCancel = context.WithCancel(context.Background())
})
