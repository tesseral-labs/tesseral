package authn

import (
	"context"
)

type ContextData struct {
	IntermediateSession *IntermediateSessionContextData
}

type IntermediateSessionContextData struct {
	OrganizationID string
	ProjectID      string
}

type ctxKey struct{}

func NewContext(ctx context.Context, data ContextData) context.Context {
	return context.WithValue(ctx, ctxKey{}, data)
}
