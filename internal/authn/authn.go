package authn

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	backendv1 "github.com/openauth-dev/openauth/internal/gen/backend/v1"
	"github.com/openauth-dev/openauth/internal/jwt"
	"github.com/openauth-dev/openauth/internal/store/idformat"
)

type ContextData struct {
	IntermediateSession *jwt.IntermediateSessionJWTClaims
	Session             *jwt.SessionJWTClaims
	ProjectAPIKey       *backendv1.ProjectAPIKey
}

type ctxKey struct{}

func NewContext(ctx context.Context, data ContextData) context.Context {
	return context.WithValue(ctx, ctxKey{}, data)
}

func NewProjectAPIKeyContext(ctx context.Context, projectAPIKey *backendv1.ProjectAPIKey) context.Context {
	return context.WithValue(ctx, ctxKey{}, ContextData{ProjectAPIKey: projectAPIKey})
}

func ProjectID(ctx context.Context) uuid.UUID {
	v, ok := ctx.Value(ctxKey{}).(ContextData)
	if !ok {
		panic("ctx does not carry authn data")
	}

	var projectID string
	switch {
	case v.IntermediateSession != nil:
		projectID = v.IntermediateSession.ProjectID
	case v.Session != nil:
		projectID = v.Session.ProjectID
	case v.ProjectAPIKey != nil:
		projectID = v.ProjectAPIKey.ProjectId
	default:
		panic("unsupported authn ctx data")
	}

	id, err := idformat.Project.Parse(projectID)
	if err != nil {
		panic(fmt.Errorf("parse project id: %w", err))
	}
	return id
}
