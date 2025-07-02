package accesstoken

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/tesseral-labs/tesseral/internal/common/store"
)

type Issuer struct {
	store *store.Store
}

func NewIssuer(store *store.Store) *Issuer {
	return &Issuer{
		store: store,
	}
}

func (i *Issuer) NewAccessToken(ctx context.Context, projectID uuid.UUID, refreshToken string) (string, error) {
	res, err := i.store.IssueAccessToken(ctx, projectID, refreshToken)
	if err != nil {
		return "", fmt.Errorf("issue access token: %w", err)
	}
	return res, nil
}
