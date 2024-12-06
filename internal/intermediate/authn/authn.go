package authn

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	intermediatev1 "github.com/openauth/openauth/internal/intermediate/gen/openauth/intermediate/v1"
)

type ctxData struct {
	intermediateSession *intermediatev1.IntermediateSession
	projectID           uuid.UUID
}

type ctxKey struct{}

func NewContext(ctx context.Context, intermediateSession *intermediatev1.IntermediateSession, projectID uuid.UUID) context.Context {
	return context.WithValue(ctx, ctxKey{}, ctxData{
		intermediateSession,
		projectID,
	})
}

func IntermediateSession(ctx context.Context) *intermediatev1.IntermediateSession {
	v, ok := ctx.Value(ctxKey{}).(ctxData)
	if !ok {
		panic(fmt.Errorf("ctx does not carry intermediate authn data"))
	}

	return v.intermediateSession
}

func ProjectID(ctx context.Context) uuid.UUID {
	v, ok := ctx.Value(ctxKey{}).(ctxData)
	if !ok {
		panic(fmt.Errorf("ctx does not carry intermediate authn data"))
	}

	return v.projectID
}

// TODO we will likely want a convenience ProjectID(ctx) uuid.UUID method here,
// as well as one for IntermediateSessionID
