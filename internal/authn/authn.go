package authn

import (
	"context"

	"github.com/openauth-dev/openauth/internal/jwt"
)

type ContextData struct {
	IntermediateSession *jwt.IntermediateSessionJWTClaims
	Session             *jwt.SessionJWTClaims
}

type ctxKey struct{}

func NewContext(ctx context.Context, data ContextData) context.Context {
	return context.WithValue(ctx, ctxKey{}, data)
}
