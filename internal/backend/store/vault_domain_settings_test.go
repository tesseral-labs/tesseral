package store

import (
	"testing"

	"github.com/stretchr/testify/require"
	backendv1 "github.com/tesseral-labs/tesseral/internal/backend/gen/tesseral/backend/v1"
)

func TestGetVaultDomainSettings(t *testing.T) {
	t.Parallel()

	ctx, u := newTestUtil(t)

	resp, err := u.Store.GetVaultDomainSettings(ctx, &backendv1.GetVaultDomainSettingsRequest{})
	require.NoError(t, err)
	require.Nil(t, resp.VaultDomainSettings)
}

func TestUpdateVaultDomainSettings(t *testing.T) {
	t.Skip("cloudflare and ses are not available in test environment")
}

func TestEnableCustomVaultDomain(t *testing.T) {
	t.Skip("cloudflare and ses are not available in test environment")
}

func TestEnableEmailSendFromDomain(t *testing.T) {
	t.Skip("cloudflare and ses are not available in test environment")
}
