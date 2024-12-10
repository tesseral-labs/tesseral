package authn

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/openauth/openauth/internal/scim/store"
	"github.com/openauth/openauth/internal/store/idformat"
)

type ctxKey struct{}

func NewContext(ctx context.Context, scimAPIKey *store.SCIMAPIKey) context.Context {
	return context.WithValue(ctx, ctxKey{}, scimAPIKey)
}

func OrganizationID(ctx context.Context) uuid.UUID {
	v, ok := ctx.Value(ctxKey{}).(*store.SCIMAPIKey)
	if !ok {
		panic(fmt.Errorf("ctx does not carry authn data"))
	}

	id, err := idformat.Organization.Parse(v.OrganizationID)
	if err != nil {
		panic(fmt.Errorf("parse organization id: %w", err))
	}

	return id
}
