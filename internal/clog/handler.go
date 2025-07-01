package clog

import (
	"context"
	"log/slog"
)

type Handler struct {
	nextHandler slog.Handler
}

func NewHandler(nextHandler slog.Handler) Handler {
	return Handler{nextHandler: nextHandler}
}

func (h Handler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.nextHandler.Enabled(ctx, level)
}

func (h Handler) Handle(ctx context.Context, rec slog.Record) error {
	if attrs, ok := ctx.Value(contextKey{}).([]slog.Attr); ok {
		rec.AddAttrs(attrs...)
	}

	return h.nextHandler.Handle(ctx, rec)
}

func (h Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return h.nextHandler.WithAttrs(attrs)
}

func (h Handler) WithGroup(name string) slog.Handler {
	return h.nextHandler.WithGroup(name)
}
