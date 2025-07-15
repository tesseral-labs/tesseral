package multislog

import (
	"context"
	"errors"
	"log/slog"
)

type Handler []slog.Handler

func (m Handler) Enabled(ctx context.Context, level slog.Level) bool {
	for _, handler := range m {
		if handler.Enabled(ctx, level) {
			return true
		}
	}
	return false
}

func (m Handler) Handle(ctx context.Context, record slog.Record) error {
	var errs []error
	for _, handler := range m {
		if handler.Enabled(ctx, record.Level) {
			errs = append(errs, handler.Handle(ctx, record))
		}
	}
	return errors.Join(errs...)
}

func (m Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	var newHandlers []slog.Handler
	for _, handler := range m {
		newHandlers = append(newHandlers, handler.WithAttrs(attrs))
	}
	return Handler(newHandlers)
}

func (m Handler) WithGroup(name string) slog.Handler {
	var newHandlers []slog.Handler
	for _, handler := range m {
		newHandlers = append(newHandlers, handler.WithGroup(name))
	}
	return Handler(newHandlers)
}
