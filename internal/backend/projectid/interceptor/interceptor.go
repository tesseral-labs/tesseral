package interceptor

import (
	"context"

	"connectrpc.com/connect"
	"github.com/openauth/openauth/internal/backend/projectid"
	"github.com/openauth/openauth/internal/backend/store"
)

func New(store *store.Store) connect.UnaryInterceptorFunc {
	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			hostHeader := req.Header().Get("Host")
			projectID, err := store.GetProjectIDByDomain(ctx, hostHeader)
			if err != nil {
				return nil, connect.NewError(connect.CodeInvalidArgument, err)
			}

			ctx = projectid.NewContext(ctx, *projectID)

			return next(ctx, req)
		}
	}
}
