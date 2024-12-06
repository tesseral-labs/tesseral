package projectid

import (
	"context"
	"errors"

	"connectrpc.com/connect"
	"github.com/google/uuid"
	"github.com/openauth/openauth/internal/store/idformat"
)

type ctxData struct {
	projectID uuid.UUID
}

type ctxKey struct{}

var ErrProjectIDHeaderRequired = errors.New("X-TODO-OpenAuth-Project-ID header is required")

func newContext(ctx context.Context, projectID uuid.UUID) context.Context {
	return context.WithValue(ctx, ctxKey{}, ctxData{
		projectID,
	})
}

func NewInterceptor() connect.UnaryInterceptorFunc {
	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			// Extract project ID from the request

			projectIDHeader := req.Header().Get("X-TODO-OpenAuth-Project-ID")
			if projectIDHeader == "" {
				return nil, connect.NewError(connect.CodeInvalidArgument, ErrProjectIDHeaderRequired)
			}

			projectID, err := idformat.Project.Parse(projectIDHeader)
			if err != nil {
				return nil, connect.NewError(connect.CodeInvalidArgument, err)
			}

			ctx = newContext(ctx, projectID)

			return next(ctx, req)
		}
	}
}

func ProjectID(ctx context.Context) uuid.UUID {
	v, ok := ctx.Value(ctxKey{}).(ctxData)
	if !ok {
		panic(errors.New("ctx does not carry project ID data"))
	}

	return v.projectID
}
