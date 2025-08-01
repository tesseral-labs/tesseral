package kms

import (
	"crypto/rand"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	awsKMSV1KeyID   = os.Getenv("TESSERAL_INTERNAL_KMS_TEST_AWS_V1_KMS_KEY_ID")
	gcpKMSV1KeyName = os.Getenv("TESSERAL_INTERNAL_KMS_TEST_GCP_V1_KMS_KEY_NAME")
)

func TestAWSKMS(t *testing.T) {
	if awsKMSV1KeyID == "" {
		t.Skip("AWS KMS key ID not set")
	}

	k, err := New(t.Context(), Config{
		Backend:       "aws_kms_v1",
		AWSKMSV1KeyID: awsKMSV1KeyID,
	})
	require.NoError(t, err)

	var plaintext [256]byte
	_, _ = rand.Read(plaintext[:]) // infallible

	ciphertext, err := k.Encrypt(t.Context(), plaintext[:])
	require.NoError(t, err)

	plaintext2, err := k.Decrypt(t.Context(), ciphertext)
	require.NoError(t, err)

	require.Equal(t, plaintext[:], plaintext2)
}

func TestGCPKMS(t *testing.T) {
	if gcpKMSV1KeyName == "" {
		t.Skip("GCP KMS key name not set")
	}

	k, err := New(t.Context(), Config{
		Backend:         "gcp_kms_v1",
		GCPKMSV1KeyName: gcpKMSV1KeyName,
	})
	require.NoError(t, err)

	var plaintext [256]byte
	_, _ = rand.Read(plaintext[:]) // infallible

	ciphertext, err := k.Encrypt(t.Context(), plaintext[:])
	require.NoError(t, err)

	plaintext2, err := k.Decrypt(t.Context(), ciphertext)
	require.NoError(t, err)

	require.Equal(t, plaintext[:], plaintext2)
}
