package authn

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/openauth/openauth/internal/store/idformat"
)

type ctxKey struct{}

type ctxData struct {
	projectID string
}

func NewContext(ctx context.Context, projectID string) context.Context {
	return context.WithValue(ctx, ctxKey{}, ctxData{
		projectID,
	})
}

func ProjectID(ctx context.Context) uuid.UUID {
	v, ok := ctx.Value(ctxKey{}).(ctxData)
	if !ok {
		panic(fmt.Errorf("ctx does not carry authn data"))
	}

	id, err := idformat.Project.Parse(v.projectID)
	if err != nil {
		panic(fmt.Errorf("parse project id: %w", err))
	}

	return id
}
