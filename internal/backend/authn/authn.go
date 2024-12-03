package authn

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	backendv1 "github.com/openauth/openauth/internal/backend/gen/openauth/backend/v1"
	"github.com/openauth/openauth/internal/store/idformat"
)

type ContextData struct {
	ProjectAPIKey  *backendv1.ProjectAPIKey
	DogfoodSession *DogfoodSessionContextData
}

type DogfoodSessionContextData struct {
	UserID           string
	OrganizationID   string
	DogfoodProjectID string
}

type ctxKey struct{}

func NewContext(ctx context.Context, data ContextData) context.Context {
	return context.WithValue(ctx, ctxKey{}, data)
}

func NewProjectAPIKeyContext(ctx context.Context, projectAPIKey *backendv1.ProjectAPIKey) context.Context {
	return context.WithValue(ctx, ctxKey{}, ContextData{ProjectAPIKey: projectAPIKey})
}

func NewDogfoodSessionContext(ctx context.Context, dogfoodSession DogfoodSessionContextData) context.Context {
	return context.WithValue(ctx, ctxKey{}, ContextData{DogfoodSession: &dogfoodSession})
}

func ProjectID(ctx context.Context) uuid.UUID {
	v, ok := ctx.Value(ctxKey{}).(ContextData)
	if !ok {
		panic("ctx does not carry authn data")
	}

	var projectID string
	switch {
	case v.ProjectAPIKey != nil:
		projectID = v.ProjectAPIKey.ProjectId
	case v.DogfoodSession != nil:
		projectID = v.DogfoodSession.DogfoodProjectID
	default:
		panic("unsupported authn ctx data")
	}

	id, err := idformat.Project.Parse(projectID)
	if err != nil {
		panic(fmt.Errorf("parse project id: %w", err))
	}
	return id
}
