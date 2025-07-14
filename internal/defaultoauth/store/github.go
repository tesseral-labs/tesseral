package store

import (
	"context"
	"crypto/sha256"
	"fmt"

	"github.com/tesseral-labs/tesseral/internal/defaultoauth/store/queries"
)

func (s *Store) GetVaultDomainByGitHubOAuthState(ctx context.Context, state string) (string, error) {
	q := queries.New(s.DB)
	stateSHA := sha256.Sum256([]byte(state))
	vaultDomain, err := q.GetVaultDomainByGitHubOAuthStateSHA256(ctx, stateSHA[:])
	if err != nil {
		return "", fmt.Errorf("get vault domain by github oauth state sha256: %w", err)
	}
	return vaultDomain, nil
}