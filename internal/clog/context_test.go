package clog_test

import (
	"context"
	"encoding/json"
	"log/slog"
	"strings"

	"github.com/AntonKosov/git-backups/internal/clog"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Context tests", func() {
	var (
		ctx           context.Context
		defaultLogger *slog.Logger
		output        strings.Builder
	)

	writeLog := func() map[string]any {
		slog.InfoContext(ctx, "test message")
		result := output.String()
		decoder := json.NewDecoder(strings.NewReader(result))
		values := map[string]any{}
		err := decoder.Decode(&values)
		Expect(err).NotTo(HaveOccurred())

		return values
	}

	verifyLogs := func(values, attrs map[string]any) {
		Expect(values).To(HaveLen(3 + len(attrs)))
		Expect(values["time"]).ToNot(BeEmpty())
		Expect(values["level"]).To(Equal("INFO"))
		Expect(values["msg"]).To(Equal("test message"))
		for key, val := range attrs {
			Expect(values).To(HaveKey(key))
			Expect(values[key]).To(BeEquivalentTo(val))
		}
	}

	BeforeEach(func() {
		output = strings.Builder{}
		handler := clog.NewHandler(slog.NewJSONHandler(&output, &slog.HandlerOptions{Level: slog.LevelDebug}))
		defaultLogger = slog.Default()
		slog.SetDefault(slog.New(handler))

		ctx = context.Background()
	})

	AfterEach(func() {
		slog.SetDefault(defaultLogger)
	})

	It("has no attributes", func() {
		values := writeLog()
		verifyLogs(values, map[string]any{})
	})

	It("two correct attributes", func() {
		ctx = clog.Add(ctx, "key1", 1, "key2", 2.0)
		values := writeLog()
		verifyLogs(values, map[string]any{
			"key1": 1,
			"key2": 2.0,
		})
	})

	It("received attributes twice", func() {
		ctx = clog.Add(ctx, "key1", 1, "key2", 2.0)
		ctx = clog.Add(ctx, "key3", "abc", "key4", 'R')
		values := writeLog()
		verifyLogs(values, map[string]any{
			"key1": 1,
			"key2": 2.0,
			"key3": "abc",
			"key4": 'R',
		})
	})

	It("received bad key", func() {
		ctx = clog.Add(ctx, struct{}{}, 1)
		values := writeLog()
		verifyLogs(values, map[string]any{
			"bad_key": 1,
		})
	})

	It("received key without value", func() {
		ctx = clog.Add(ctx, "key")
		values := writeLog()
		verifyLogs(values, map[string]any{
			"key": "missing_value",
		})
	})

	It("received Stringer key", func() {
		var sb strings.Builder
		sb.WriteString("stringer")
		ctx = clog.Add(ctx, &sb, 1)
		values := writeLog()
		verifyLogs(values, map[string]any{
			"stringer": 1,
		})
	})
})
