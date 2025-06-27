package httplog

import (
	"context"
	"log/slog"

	"connectrpc.com/connect"
)

func NewInterceptor() connect.Interceptor {
	return connect.UnaryInterceptorFunc(func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			slog.InfoContext(ctx, "http_request", "x_tesseral_host", req.Header().Get("X-Tesseral-Host"), "rpc", req.Spec().Procedure, "user_agent", req.Header().Get("User-Agent"), "traceparent", req.Header().Get("traceparent"))
			res, err := next(ctx, req)

			var errorCode string
			if err != nil {
				errorCode = connect.CodeOf(err).String()
			}

			// for convenience, log request details here too
			slog.InfoContext(ctx, "http_response", "error_code", errorCode, "x_tesseral_host", req.Header().Get("X-Tesseral-Host"), "rpc", req.Spec().Procedure, "user_agent", req.Header().Get("User-Agent"))
			return res, err
		}
	})
}
