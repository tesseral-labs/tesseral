package main

import (
	"context"
	"errors"
	"log/slog"
)

type multiOutputHandler []slog.Handler

func (m multiOutputHandler) Enabled(ctx context.Context, level slog.Level) bool {
	for _, handler := range m {
		if handler.Enabled(ctx, level) {
			return true
		}
	}
	return false
}

func (m multiOutputHandler) Handle(ctx context.Context, record slog.Record) error {
	var errs []error
	for _, handler := range m {
		if handler.Enabled(ctx, record.Level) {
			errs = append(errs, handler.Handle(ctx, record))
		}
	}
	return errors.Join(errs...)
}

func (m multiOutputHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	var newHandlers []slog.Handler
	for _, handler := range m {
		newHandlers = append(newHandlers, handler.WithAttrs(attrs))
	}
	return multiOutputHandler(newHandlers)
}

func (m multiOutputHandler) WithGroup(name string) slog.Handler {
	var newHandlers []slog.Handler
	for _, handler := range m {
		newHandlers = append(newHandlers, handler.WithGroup(name))
	}
	return multiOutputHandler(newHandlers)
}
