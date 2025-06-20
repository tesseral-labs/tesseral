package store

import (
	"testing"

	"github.com/stretchr/testify/require"
	backendv1 "github.com/tesseral-labs/tesseral/internal/backend/gen/tesseral/backend/v1"
)

func TestGetProjectEntitlements(t *testing.T) {
	t.Parallel()

	ctx, u := newTestUtil(t)

	resp, err := u.Store.GetProjectEntitlements(ctx, &backendv1.GetProjectEntitlementsRequest{})
	require.NoError(t, err)
	require.NotNil(t, resp)

	// Default entitlements from storetesting
	require.True(t, resp.EntitledCustomVaultDomains)
	require.True(t, resp.EntitledBackendApiKeys)
}

func TestCreateStripeCheckoutLink(t *testing.T) {
	t.Skip("stripe is not available in test environment")
}
