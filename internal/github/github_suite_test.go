package github_test

import (
	"context"
	"io"
	"log/slog"
	"testing"

	"github.com/jarcoal/httpmock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var ctx context.Context

func TestGithub(t *testing.T) {
	slog.SetDefault(slog.New(slog.NewTextHandler(io.Discard, nil)))
	RegisterFailHandler(Fail)
	RunSpecs(t, "Github Suite")
}

var _ = BeforeSuite(func() {
	httpmock.Activate()
})

var _ = BeforeEach(func() {
	ctx = context.Background()
	httpmock.Reset()
})

var _ = AfterSuite(func() {
	httpmock.DeactivateAndReset()
})
