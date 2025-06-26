package opaqueinternalerror

import (
	"context"
	"errors"
	"log/slog"

	"connectrpc.com/connect"
	"github.com/getsentry/sentry-go"
)

func NewInterceptor() connect.Interceptor {
	return connect.UnaryInterceptorFunc(func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			// Mirrors sentryhttp behavior
			hub := sentry.GetHubFromContext(ctx)
			if hub == nil {
				hub = sentry.CurrentHub().Clone()
				ctx = sentry.SetHubOnContext(ctx, hub)
			}

			res, err := next(ctx, req)
			if err != nil {
				var connectErr *connect.Error
				if errors.As(err, &connectErr) {
					return nil, connectErr
				}

				slog.ErrorContext(ctx, "internal_error", "err", err)
				hub.CaptureException(err)

				return nil, connect.NewError(connect.CodeInternal, errors.New("internal server error"))
			}

			return res, nil
		}
	})
}
