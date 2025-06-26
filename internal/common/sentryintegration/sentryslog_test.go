package sentryintegration

import (
	"context"
	"encoding/json"
	"log/slog"
	"strings"
	"testing"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/stretchr/testify/require"
)

type customValue string

// LogValue implements slog.LogValuer.
func (c customValue) LogValue() slog.Value {
	return slog.StringValue(strings.Repeat(string(c), 2))
}

func TestSlogHandler(t *testing.T) {
	t.Parallel()

	handler := NewSlogHandler(slog.NewTextHandler(new(strings.Builder), nil))

	record := slog.NewRecord(time.Now(), slog.LevelInfo, "test message", 0)

	record.AddAttrs(
		slog.String("key1", "value1"),
		slog.Int("key2", 42),
		slog.Bool("key3", true),
		slog.Group("nested_group",
			slog.String("nested_key1", "nested_value1"),
			slog.Int("nested_key2", 100),
		),
		slog.Duration("duration_key", 1234567890),
		slog.Float64("float_key", 3.14),
		slog.Int64("int64_key", 1234567890123456789),
		slog.Uint64("uint64_key", 1234567890123456789),
		slog.Any("any_key", map[string]any{
			"foo": "bar",
			"baz": 123,
		}),
	)
	record.Add("custom_value_key", customValue("custom_value"))

	values := make(map[string]any)
	record.Attrs(func(a slog.Attr) bool {
		values[a.Key] = attrValue(a.Value)
		return true
	})
	expected := map[string]any{
		"key1": "value1",
		"key2": int64(42),
		"key3": true,
		"nested_group": map[string]any{
			"nested_key1": "nested_value1",
			"nested_key2": int64(100),
		},
		"duration_key": time.Duration(1234567890),
		"float_key":    3.14,
		"int64_key":    int64(1234567890123456789),
		"uint64_key":   uint64(1234567890123456789),
		"any_key": map[string]any{
			"foo": "bar",
			"baz": 123,
		},
		"custom_value_key": "custom_valuecustom_value",
	}
	require.Equal(t, expected, values)

	// Ensure the record can be marshaled to JSON as required by Sentry.
	_, err := json.Marshal(values)
	require.NoError(t, err)

	// Test record handling
	hub := sentry.NewHub(nil, sentry.NewScope())
	ctx := sentry.SetHubOnContext(context.Background(), hub)

	err = handler.Handle(ctx, record)
	require.NoError(t, err)
}
