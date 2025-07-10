package authn

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	intermediatev1 "github.com/tesseral-labs/tesseral/internal/intermediate/gen/tesseral/intermediate/v1"
	"github.com/tesseral-labs/tesseral/internal/store/idformat"
)

type ctxData struct {
	intermediateSession *intermediatev1.IntermediateSession
	projectID           uuid.UUID
}

type ctxKey struct{}

func NewContext(ctx context.Context, intermediateSession *intermediatev1.IntermediateSession, projectID uuid.UUID) context.Context {
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

func IntermediateSessionID(ctx context.Context) *uuid.UUID {
	v, ok := ctx.Value(ctxKey{}).(ctxData)
	if !ok {
		panic(errors.New("ctx does not carry intermediate session ID data"))
	}

	if v.intermediateSession == nil {
		return nil
	}

	id, err := idformat.IntermediateSession.Parse(v.intermediateSession.Id)
	if err != nil {
		panic(fmt.Errorf("parse intermediate session id: %w", err))
	}
	return (*uuid.UUID)(&id)
}
