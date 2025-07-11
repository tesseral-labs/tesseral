package authn

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/tesseral-labs/tesseral/internal/oidc/store/queries"
)

type ctxData struct {
	intermediateSession *queries.IntermediateSession
	projectID           uuid.UUID
}

type ctxKey struct{}

func NewContext(ctx context.Context, intermediateSession *queries.IntermediateSession, projectID uuid.UUID) context.Context {
	return context.WithValue(ctx, ctxKey{}, ctxData{
		intermediateSession: intermediateSession,
		projectID:           projectID,
	})
}
func ProjectID(ctx context.Context) uuid.UUID {
	v, ok := ctx.Value(ctxKey{}).(ctxData)
	if !ok {
		panic(errors.New("ctx does not carry project ID data"))
	}

	return v.projectID
}

func IntermediateSession(ctx context.Context) *queries.IntermediateSession {
	v, ok := ctx.Value(ctxKey{}).(ctxData)
	if !ok {
		panic(errors.New("ctx does not carry intermediate session data"))
	}

	return v.intermediateSession
}
