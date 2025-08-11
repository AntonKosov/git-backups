package backup_test

import (
	"context"
	"io"
	"log/slog"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var ctx context.Context

func TestBackup(t *testing.T) {
	slog.SetDefault(slog.New(slog.NewTextHandler(io.Discard, nil)))
	RegisterFailHandler(Fail)
	RunSpecs(t, "Backup Suite")
}

var _ = BeforeEach(func() {
	ctx = context.Background()
})
