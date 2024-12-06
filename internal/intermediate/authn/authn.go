package authn

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	intermediatev1 "github.com/openauth/openauth/internal/intermediate/gen/openauth/intermediate/v1"
	"github.com/openauth/openauth/internal/store/idformat"
)

type ctxData struct {
	intermediateSession *intermediatev1.IntermediateSession
}

type ctxKey struct{}

func NewContext(ctx context.Context, intermediateSession *intermediatev1.IntermediateSession) context.Context {
	return context.WithValue(ctx, ctxKey{}, ctxData{
		intermediateSession,
	})
}

func IntermediateSession(ctx context.Context) *intermediatev1.IntermediateSession {
	v, ok := ctx.Value(ctxKey{}).(ctxData)
	if !ok {
		panic(fmt.Errorf("ctx does not carry intermediate authn data"))
	}

	return v.intermediateSession
}

func IntermediateSessionID(ctx context.Context) uuid.UUID {
	id, err := idformat.IntermediateSession.Parse(IntermediateSession(ctx).Id)
	if err != nil {
		panic(fmt.Errorf("parse intermediate session id: %w", err))
	}
	return id
}

func ProjectID(ctx context.Context) uuid.UUID {
	id, err := idformat.Project.Parse(IntermediateSession(ctx).ProjectId)
	if err != nil {
		panic(fmt.Errorf("parse project id: %w", err))
	}
	return id
}
