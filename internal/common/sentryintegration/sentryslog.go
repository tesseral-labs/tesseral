package sentryintegration

import (
	"context"
	"log/slog"
	"runtime"
	"strings"

	"github.com/getsentry/sentry-go"
)

type sentryHandler struct {
	base slog.Handler
}

var _ slog.Handler = (*sentryHandler)(nil)

// NewSlogHandler creates an slog handler which inserts breadcrumbs into the current breadcrumb stack for each log record.
//
// These breadcrumbs are associated with a panic or manual Sentry event if they occur and are otherwise discarded.
func NewSlogHandler(base slog.Handler) slog.Handler {
	return &sentryHandler{base: base}
}

func sentryLevel(level slog.Level) sentry.Level {
	levels := strings.Split(level.String(), "+")
	return sentry.Level(strings.ToLower(levels[0]))
}

func attrValue(value slog.Value) any {
	switch value.Kind() {
	case slog.KindGroup:
		group := value.Group()
		groupAttrs := make(map[string]any, len(group))
		for _, a := range group {
			groupAttrs[a.Key] = attrValue(a.Value)
		}
		return groupAttrs
	case slog.KindLogValuer:
		return attrValue(value.Resolve())
	default:
		anyValue := value.Any()
		switch v := anyValue.(type) {
		case error:
			return v.Error() // Convert error to string
		case []byte:
			return string(v) // Convert byte slice to string
		default:
			return anyValue
		}
	}
}

// Enabled implements slog.Handler.
func (s *sentryHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return s.base.Enabled(ctx, level)
}

// Handle implements slog.Handler.
func (s *sentryHandler) Handle(ctx context.Context, record slog.Record) error {
	attrs := make(map[string]any)
	record.Attrs(func(a slog.Attr) bool {
		attrs[a.Key] = attrValue(a.Value)
		return true
	})
	addSource(record, attrs)

	hub := sentry.GetHubFromContext(ctx)

	// Don't add breadcrumbs to the global hub.
	if hub != nil {
		hub.AddBreadcrumb(&sentry.Breadcrumb{
			Type:      "default",
			Data:      attrs,
			Category:  "slog",
			Message:   record.Message,
			Level:     sentryLevel(record.Level),
			Timestamp: record.Time,
		}, nil)
	}

	return s.base.Handle(ctx, record)
}

// WithAttrs implements slog.Handler.
func (s *sentryHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &sentryHandler{base: s.base.WithAttrs(attrs)}
}

// WithGroup implements slog.Handler.
func (s *sentryHandler) WithGroup(name string) slog.Handler {
	return &sentryHandler{base: s.base.WithGroup(name)}
}

func addSource(r slog.Record, attrs map[string]any) {
	fs := runtime.CallersFrames([]uintptr{r.PC})
	f, _ := fs.Next()

	source := make(map[string]any, 3)
	if f.File != "" {
		source["file"] = f.File
	}
	if f.Line > 0 {
		source["line"] = f.Line
	}
	if f.Function != "" {
		source["function"] = f.Function
	}
	if len(source) == 0 {
		return
	}
	attrs["source"] = source
}
