package clog

import (
	"context"
	"fmt"
	"log/slog"
)

const (
	badKey       = "bad_key"
	missingValue = "missing_value"
)

type contextKey struct{}

func Add(ctx context.Context, kv ...any) context.Context {
	var attrs []slog.Attr
	if prevAttrs, ok := ctx.Value(contextKey{}).([]slog.Attr); ok {
		attrs = make([]slog.Attr, len(prevAttrs), len(prevAttrs)+(len(kv)+1)/2)
		copy(attrs, prevAttrs)
	} else {
		attrs = make([]slog.Attr, (len(kv)+1)/2)
	}

	for i := 0; i < len(kv); i += 2 {
		var key string
		switch item := kv[i].(type) {
		case string:
			key = item
		case fmt.Stringer:
			key = item.String()
		default:
			key = badKey
		}

		val := any(missingValue)
		if i+1 < len(kv) {
			val = kv[i+1]
		}

		attrs = append(attrs, slog.Any(key, val))
	}

	return context.WithValue(ctx, contextKey{}, attrs)
}
