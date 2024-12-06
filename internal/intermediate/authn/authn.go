package authn

import (
	"context"
	"fmt"

	intermediatev1 "github.com/openauth/openauth/internal/intermediate/gen/openauth/intermediate/v1"
)

type ctxData struct {
	intermediateSession *intermediatev1.IntermediateSession
}

type ctxKey struct{}

func NewContext(ctx context.Context, intermediateSession *intermediatev1.IntermediateSession) context.Context {
	return context.WithValue(ctx, ctxKey{}, ctxData{intermediateSession})
}

func IntermediateSession(ctx context.Context) *intermediatev1.IntermediateSession {
	v, ok := ctx.Value(ctxKey{}).(ctxData)
	if !ok {
		panic(fmt.Errorf("ctx does not carry intermediate authn data"))
	}

	return v.intermediateSession
}

// TODO we will likely want a convenience ProjectID(ctx) uuid.UUID method here,
// as well as one for IntermediateSessionID
