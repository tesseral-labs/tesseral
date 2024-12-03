package authn

import (
	"context"
)

type ContextData struct {
	Session *SessionContextData
}

type SessionContextData struct {
	UserID         string
	OrganizationID string
	ProjectID      string
}

type ctxKey struct{}

func NewContext(ctx context.Context, data ContextData) context.Context {
	return context.WithValue(ctx, ctxKey{}, data)
}
