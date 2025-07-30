package internal_test

import (
	"io"
	"log/slog"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestInternal(t *testing.T) {
	slog.SetDefault(slog.New(slog.NewTextHandler(io.Discard, nil)))
	RegisterFailHandler(Fail)
	RunSpecs(t, "Internal Suite")
}
