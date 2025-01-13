package authn

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/openauth/openauth/internal/store/idformat"
)

type ContextData struct {
	SessionID      string
	UserID         string
	OrganizationID string
	ProjectID      string
}

type ctxKey struct{}

func NewContext(ctx context.Context, data ContextData) context.Context {
	return context.WithValue(ctx, ctxKey{}, data)
}

func OrganizationID(ctx context.Context) uuid.UUID {
	v, ok := ctx.Value(ctxKey{}).(ContextData)
	if !ok {
		panic("ctx does not carry authn data")
	}
	orgID, err := idformat.Organization.Parse(v.OrganizationID)
	if err != nil {
		panic(fmt.Errorf("parse organization id: %w", err))
	}
	return orgID
}

func ProjectID(ctx context.Context) uuid.UUID {
	v, ok := ctx.Value(ctxKey{}).(ContextData)
	if !ok {
		panic("ctx does not carry authn data")
	}
	projectID, err := idformat.Project.Parse(v.ProjectID)
	if err != nil {
		panic(fmt.Errorf("parse project id: %w", err))
	}
	return projectID
}

func UserID(ctx context.Context) uuid.UUID {
	v, ok := ctx.Value(ctxKey{}).(ContextData)
	if !ok {
		panic("ctx does not carry authn data")
	}
	userID, err := idformat.User.Parse(v.UserID)
	if err != nil {
		panic(fmt.Errorf("parse user id: %w", err))
	}
	return userID
}
