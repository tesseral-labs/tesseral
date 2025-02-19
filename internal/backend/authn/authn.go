package authn

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/tesseral-labs/tesseral/internal/store/idformat"
)

type ContextData struct {
	ProjectAPIKey  *ProjectAPIKeyContextData
	DogfoodSession *DogfoodSessionContextData
}

type ProjectAPIKeyContextData struct {
	ProjectAPIKeyID string
	ProjectID       string
}

// DogfoodSessionContextData contains data related to a user logged into
// app.tesseral.com.
type DogfoodSessionContextData struct {
	UserID string

	// ProjectID is the ID of the project the user is manipulating. This is
	// almost never the same thing as the dogfood project.
	ProjectID string
}

type ctxKey struct{}

func NewProjectAPIKeyContext(ctx context.Context, projectAPIKey *ProjectAPIKeyContextData) context.Context {
	return context.WithValue(ctx, ctxKey{}, ContextData{ProjectAPIKey: projectAPIKey})
}

func NewDogfoodSessionContext(ctx context.Context, dogfoodSession DogfoodSessionContextData) context.Context {
	return context.WithValue(ctx, ctxKey{}, ContextData{DogfoodSession: &dogfoodSession})
}

func GetContextData(ctx context.Context) ContextData {
	v, ok := ctx.Value(ctxKey{}).(ContextData)
	if !ok {
		panic("ctx does not carry authn data")
	}
	return v
}

func ProjectID(ctx context.Context) uuid.UUID {
	v, ok := ctx.Value(ctxKey{}).(ContextData)
	if !ok {
		panic("ctx does not carry authn data")
	}

	var projectID string
	switch {
	case v.ProjectAPIKey != nil:
		projectID = v.ProjectAPIKey.ProjectID
	case v.DogfoodSession != nil:
		projectID = v.DogfoodSession.ProjectID
	default:
		panic("unsupported authn ctx data")
	}

	id, err := idformat.Project.Parse(projectID)
	if err != nil {
		panic(fmt.Errorf("parse project id: %w", err))
	}
	return id
}
