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
			res, err := next(ctx, req)
			if err != nil {
				var connectErr *connect.Error
				if errors.As(err, &connectErr) {
					return nil, connectErr
				}

				slog.ErrorContext(ctx, "internal_error", "err", err)
				sentry.CaptureException(err)
				return nil, connect.NewError(connect.CodeInternal, errors.New("internal server error"))
			}

			return res, nil
		}
	})
}
