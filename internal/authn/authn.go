package authn

import (
	"context"

	"github.com/openauth-dev/openauth/internal/store"
)

type ContextData struct {
	IntermediateSession *store.IntermediateSessionJWTClaims
	Session             *store.SessionJWTClaims
}

type ctxKey struct{}

func NewContext(ctx context.Context, data ContextData) context.Context {
	return context.WithValue(ctx, ctxKey{}, data)
}
