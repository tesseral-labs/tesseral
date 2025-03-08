package authn

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/tesseral-labs/tesseral/internal/store/idformat"
)

type ctxKey struct{}

type ctxData struct {
	scimAPIKey *SCIMAPIKey
	projectID  string
}

type SCIMAPIKey struct {
	ID             string
	OrganizationID string
}

func NewContext(ctx context.Context, scimAPIKey *SCIMAPIKey, projectID string) context.Context {
	return context.WithValue(ctx, ctxKey{}, ctxData{
		scimAPIKey,
		projectID,
	})
}

func OrganizationID(ctx context.Context) uuid.UUID {
	v, ok := ctx.Value(ctxKey{}).(ctxData)
	if !ok {
		panic(fmt.Errorf("ctx does not carry authn data"))
	}

	if v.scimAPIKey == nil {
		panic(fmt.Errorf("ctx does not carry scimAPIKey"))
	}

	id, err := idformat.Organization.Parse(v.scimAPIKey.OrganizationID)
	if err != nil {
		panic(fmt.Errorf("parse organization id: %w", err))
	}

	return id
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
