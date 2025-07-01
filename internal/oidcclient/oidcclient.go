package oidcclient

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"slices"
	"strings"
)

type Client struct {
	HTTPClient *http.Client
}

type Configuration struct {
	Issuer                            string   `json:"issuer"`
	AuthorizationEndpoint             string   `json:"authorization_endpoint"`
	TokenEndpoint                     string   `json:"token_endpoint"`
	TokenEndpointAuthMethodsSupported []string `json:"token_endpoint_auth_methods_supported"`
	JWKSURI                           string   `json:"jwks_uri"`
	GrantTypesSupported               []string `json:"grant_types_supported"`
	CodeChallengeMethodsSupported     []string `json:"code_challenge_methods_supported"`
	IDTokenSigningAlgValuesSupported  []string `json:"id_token_signing_alg_values_supported"`
}

func (c *Client) GetConfiguration(ctx context.Context, configURL string) (*Configuration, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, configURL, nil)
	if err != nil {
		return nil, fmt.Errorf("create request for OIDC configuration: %w", err)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetch OIDC configuration: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("fetch OIDC configuration: unexpected status code %d", resp.StatusCode)
	}

	var config Configuration
	if err := json.NewDecoder(resp.Body).Decode(&config); err != nil {
		return nil, fmt.Errorf("decode OIDC configuration: %w", err)
	}

	return &config, nil
}

// Validate performs basic validation on the OIDC configuration.
func (c *Configuration) Validate() error {
	if c.AuthorizationEndpoint == "" {
		return fmt.Errorf("authorization endpoint is required")
	}
	if c.TokenEndpoint == "" {
		return fmt.Errorf("token endpoint is required")
	}
	if c.JWKSURI == "" {
		return fmt.Errorf("jwks uri is required")
	}

	// Sanity checks for downstream OIDC operations.
	//
	// If the OIDC configuration is not well-formed or missing these values, we are okay failing hard later.

	if len(c.GrantTypesSupported) != 0 {
		if !slices.Contains(c.GrantTypesSupported, "authorization_code") {
			return fmt.Errorf("grant type 'authorization_code' is required")
		}
	}
	if len(c.TokenEndpointAuthMethodsSupported) != 0 {
		if !slices.Contains(c.TokenEndpointAuthMethodsSupported, "client_secret_post") &&
			!slices.Contains(c.TokenEndpointAuthMethodsSupported, "client_secret_basic") {
			return fmt.Errorf("token endpoint auth method must be either 'client_secret_post' or 'client_secret_basic'")
		}
	}
	if len(c.CodeChallengeMethodsSupported) != 0 {
		if !slices.Contains(c.CodeChallengeMethodsSupported, "S256") {
			return fmt.Errorf("code challenge method 'S256' is required")
		}
	}
	if len(c.IDTokenSigningAlgValuesSupported) != 0 {
		if !slices.Contains(c.IDTokenSigningAlgValuesSupported, "RS256") {
			return fmt.Errorf("ID token signing algorithm 'RS256' is required")
		}
	}
	return nil
}

type ExchangeCodeRequest struct {
	TokenEndpoint   string
	Code            string
	RedirectURI     string
	ClientID        string
	ClientAuthBasic string
	ClientAuthPost  string
	CodeVerifier    *string // Optional, used for PKCE
}

type ExchangeCodeResponse struct {
	IDToken string `json:"id_token"`
}

func (c *Client) ExchangeCode(ctx context.Context, req ExchangeCodeRequest) (*ExchangeCodeResponse, error) {
	requestBody := url.Values{}
	requestBody.Set("grant_type", "authorization_code")
	requestBody.Set("code", req.Code)
	requestBody.Set("redirect_uri", req.RedirectURI)

	if req.ClientAuthPost != "" {
		requestBody.Set("client_id", req.ClientID)
		requestBody.Set("client_secret", req.ClientAuthPost)
	}
	if req.CodeVerifier != nil {
		requestBody.Set("code_verifier", *req.CodeVerifier)
	}

	httpRequest, err := http.NewRequestWithContext(ctx, http.MethodPost, req.TokenEndpoint, strings.NewReader(requestBody.Encode()))
	if err != nil {
		return nil, fmt.Errorf("create token request: %w", err)
	}
	httpRequest.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if req.ClientAuthBasic != "" {
		httpRequest.Header.Set("Authorization", "Basic "+req.ClientAuthBasic)
	}
	httpResponse, err := c.HTTPClient.Do(httpRequest)
	if err != nil {
		return nil, fmt.Errorf("send token request: %w", err)
	}
	defer func() {
		_ = httpResponse.Body.Close()
	}()
	if httpResponse.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(httpResponse.Body)
		return nil, fmt.Errorf("OIDC token exchange failed: %s\n%s", httpResponse.Status, body)
	}

	var response ExchangeCodeResponse
	if err := json.NewDecoder(httpResponse.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("decode token response: %w", err)
	}
	if response.IDToken == "" {
		return nil, fmt.Errorf("OIDC token response does not contain an ID token")
	}

	return &response, nil
}

func (c *Client) GenerateCodeVerifierAndChallenge() (string, string, error) {
	randomBytes := make([]byte, 64)
	if _, err := rand.Read(randomBytes); err != nil {
		return "", "", fmt.Errorf("failed to generate random bytes for code verifier: %w", err)
	}
	codeVerifier := base64.RawURLEncoding.EncodeToString(randomBytes)
	hash := sha256.Sum256([]byte(codeVerifier))
	codeChallenge := base64.RawURLEncoding.EncodeToString(hash[:])
	return codeVerifier, codeChallenge, nil
}
