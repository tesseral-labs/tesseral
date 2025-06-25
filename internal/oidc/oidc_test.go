package oidc

import (
	"context"
	"encoding/base64"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestGetConfiguration_Google(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client := &Client{HTTPClient: http.DefaultClient}
	config, err := client.GetConfiguration(ctx, "https://accounts.google.com/.well-known/openid-configuration")
	require.NoError(t, err)
	require.NotNil(t, config)
	require.NotEmpty(t, config.AuthorizationEndpoint)
	require.NotEmpty(t, config.TokenEndpoint)
	require.NotEmpty(t, config.JwksURI)
	require.NotEmpty(t, config.GrantTypesSupported)
	require.NotEmpty(t, config.TokenEndpointAuthMethodsSupported)
}

func TestGenerateCodeVerifierAndChallenge(t *testing.T) {
	t.Parallel()

	client := &Client{}
	verifier, challenge, err := client.GenerateCodeVerifierAndChallenge()
	require.NoError(t, err)

	_, err = base64.RawURLEncoding.DecodeString(verifier)
	require.NoError(t, err)

	_, err = base64.RawURLEncoding.DecodeString(challenge)
	require.NoError(t, err)
}
