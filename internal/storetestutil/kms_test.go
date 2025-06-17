package storetestutil

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go-v2/service/kms/types"
	"github.com/stretchr/testify/require"
)

func TestNewKMS(t *testing.T) {
	t.Parallel()

	client, cleanup := newKMS()
	t.Cleanup(cleanup)

	_, err := client.CreateKey(t.Context(), &kms.CreateKeyInput{
		KeySpec:  types.KeySpecEccNistP256,
		KeyUsage: types.KeyUsageTypeSignVerify,
	})
	require.NoError(t, err, "failed to create KMS key")
}
