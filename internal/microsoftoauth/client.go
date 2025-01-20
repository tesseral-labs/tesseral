package microsoftoauth

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// https://learn.microsoft.com/en-us/entra/identity-platform/id-token-claims-reference:
//
// "For sign-ins to the personal Microsoft account tenant (services like Xbox,
// Teams for Life, or Outlook), the value is
// 9188040d-6c67-4c5b-b112-36a304b66dad."
const microsoftPersonalTenantID = "9188040d-6c67-4c5b-b112-36a304b66dad"

type Client struct {
	HTTPClient *http.Client
}

type GetAuthorizeURLRequest struct {
	MicrosoftOAuthClientID string
	RedirectURI            string
	State                  string
}

func GetAuthorizeURL(req *GetAuthorizeURLRequest) string {
	u, err := url.Parse("https://login.microsoftonline.com/common/oauth2/v2.0/authorize")
	if err != nil {
		panic(fmt.Errorf("parse authorize base url: %w", err))
	}

	q := url.Values{}
	q.Set("client_id", req.MicrosoftOAuthClientID)
	q.Set("redirect_uri", req.RedirectURI)
	q.Set("state", req.State)
	q.Set("response_type", "code")
	q.Set("scope", "openid email profile") // for id_token, email, and oid

	u.RawQuery = q.Encode()
	return u.String()
}

type RedeemCodeRequest struct {
	MicrosoftOAuthClientID     string
	MicrosoftOAuthClientSecret string
	RedirectURI                string
	Code                       string
}

type RedeemCodeResponse struct {
	MicrosoftUserID   string
	Email             string
	MicrosoftTenantID string
}

func (c *Client) RedeemCode(ctx context.Context, req *RedeemCodeRequest) (*RedeemCodeResponse, error) {
	idToken, err := c.getIDToken(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("redeem code: %w", err)
	}

	res, err := parseIDToken(idToken)
	if err != nil {
		return nil, fmt.Errorf("parse id_token: %w", err)
	}

	return res, nil
}

func (c *Client) getIDToken(ctx context.Context, req *RedeemCodeRequest) (string, error) {
	body := url.Values{}
	body.Set("client_id", req.MicrosoftOAuthClientID)
	body.Set("client_secret", req.MicrosoftOAuthClientSecret)
	body.Set("code", req.Code)
	body.Set("redirect_uri", req.RedirectURI)
	body.Set("grant_type", "authorization_code")

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://login.microsoftonline.com/common/oauth2/v2.0/token", strings.NewReader(body.Encode()))
	if err != nil {
		return "", fmt.Errorf("new http request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	httpRes, err := c.HTTPClient.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("send http request: %w", err)
	}
	defer func() { _ = httpRes.Body.Close() }()

	if httpRes.StatusCode != http.StatusOK {
		return "", fmt.Errorf("bad response status code: %s", httpRes.Status)
	}

	resBody, err := io.ReadAll(httpRes.Body)
	if err != nil {
		return "", fmt.Errorf("read body: %w", err)
	}

	var data struct {
		IDToken string `json:"id_token"`
	}
	if err := json.Unmarshal(resBody, &data); err != nil {
		return "", fmt.Errorf("unmarshal response: %w", err)
	}

	return data.IDToken, nil
}

func parseIDToken(idToken string) (*RedeemCodeResponse, error) {
	// We get id_token from Microsoft's OAuth endpoint; we have already
	// established its authenticity that way. We don't need to do JWKS-based
	// authentication of id_token here.
	claimsBase64 := strings.Split(idToken, ".")[1]
	claimsJSON, err := base64.RawURLEncoding.DecodeString(claimsBase64)
	if err != nil {
		return nil, fmt.Errorf("base64 decode claims: %w", err)
	}

	var claims struct {
		OID   string `json:"oid"`
		Email string `json:"email"`
		TID   string `json:"tid"`
	}
	if err := json.Unmarshal(claimsJSON, &claims); err != nil {
		return nil, fmt.Errorf("unmarshal claims: %w", err)
	}

	// For safety-by-default for callers, treat the personal tenant as a
	// non-tenant.
	if claims.TID == microsoftPersonalTenantID {
		claims.TID = ""
	}

	return &RedeemCodeResponse{
		MicrosoftUserID:   claims.OID,
		Email:             claims.Email,
		MicrosoftTenantID: claims.TID,
	}, nil
}
