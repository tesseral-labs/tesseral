package authn

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/tesseral-labs/tesseral/internal/store/idformat"
)

type ContextData struct {
	ProjectAPIKey  *BackendAPIKeyContextData
	ConsoleSession *ConsoleSessionContextData
}

type BackendAPIKeyContextData struct {
	BackendAPIKeyID string
	ProjectID       string
}

// ConsoleSessionContextData contains data related to a user logged into
// app.tesseral.com.
type ConsoleSessionContextData struct {
	UserID    string
	SessionID string

	// ProjectID is the ID of the project the user is manipulating. This is
	// almost never the same thing as the console project.
	ProjectID string
}

type ctxKey struct{}

func NewBackendAPIKeyContext(ctx context.Context, projectAPIKey *BackendAPIKeyContextData) context.Context {
	return context.WithValue(ctx, ctxKey{}, ContextData{ProjectAPIKey: projectAPIKey})
}

func NewConsoleSessionContext(ctx context.Context, consoleSession ConsoleSessionContextData) context.Context {
	return context.WithValue(ctx, ctxKey{}, ContextData{ConsoleSession: &consoleSession})
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
	case v.ConsoleSession != nil:
		projectID = v.ConsoleSession.ProjectID
	default:
		panic("unsupported authn ctx data")
	}

	id, err := idformat.Project.Parse(projectID)
	if err != nil {
		panic(fmt.Errorf("parse project id: %w", err))
	}
	return id
}
