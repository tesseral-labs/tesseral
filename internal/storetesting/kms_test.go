package storetesting

import (
	"crypto/rand"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewKMS(t *testing.T) {
	t.Parallel()

	client, cleanup := newKMS()
	t.Cleanup(cleanup)

	var plaintext [32]byte
	_, _ = rand.Read(plaintext[:]) // infallible
	encrypted, err := client.SessionSigningKeysKMS.Encrypt(t.Context(), plaintext[:])
	require.NoError(t, err)

	plaintext2, err := client.SessionSigningKeysKMS.Decrypt(t.Context(), encrypted)
	require.NoError(t, err)

	require.Equal(t, plaintext[:], plaintext2)
}
