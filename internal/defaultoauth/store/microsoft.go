package store

import (
	"context"
	"crypto/sha256"
	"fmt"

	"github.com/tesseral-labs/tesseral/internal/defaultoauth/store/queries"
)

func (s *Store) GetVaultDomainByMicrosoftOAuthState(ctx context.Context, state string) (string, error) {
	q := queries.New(s.DB)
	stateSHA := sha256.Sum256([]byte(state))
	vaultDomain, err := q.GetVaultDomainByMicrosoftOAuthStateSHA256(ctx, stateSHA[:])
	if err != nil {
		return "", fmt.Errorf("get vault domain by microsoft oauth state sha256: %w", err)
	}
	return vaultDomain, nil
}
