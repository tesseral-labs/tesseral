package accesstoken

import (
	"context"
	"fmt"

	"github.com/openauth/openauth/internal/common/store"
)

type Issuer struct {
	store *store.Store
}

func NewIssuer(store *store.Store) *Issuer {
	return &Issuer{
		store: store,
	}
}

func (i *Issuer) NewAccessToken(ctx context.Context, refreshToken string) (string, error) {
	res, err := i.store.IssueAccessToken(ctx, refreshToken)
	if err != nil {
		return "", fmt.Errorf("issue access token: %w", err)
	}
	return res, nil
}
