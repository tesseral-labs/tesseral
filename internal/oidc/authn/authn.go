package authn

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

type ctxData struct {
	projectID uuid.UUID
}

type ctxKey struct{}

func NewContext(ctx context.Context, projectID uuid.UUID) context.Context {
	return context.WithValue(ctx, ctxKey{}, ctxData{projectID})
}

func ProjectID(ctx context.Context) uuid.UUID {
	v, ok := ctx.Value(ctxKey{}).(ctxData)
	if !ok {
		panic(errors.New("ctx does not carry project ID data"))
	}

	return v.projectID
}
