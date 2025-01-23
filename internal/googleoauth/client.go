package googleoauth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type Client struct {
	HTTPClient *http.Client
}

type GetAuthorizeURLRequest struct {
	GoogleOAuthClientID string
	RedirectURI         string
	State               string
}

func GetAuthorizeURL(req *GetAuthorizeURLRequest) string {
	u, err := url.Parse("https://accounts.google.com/o/oauth2/v2/auth")
	if err != nil {
		panic(fmt.Errorf("parse authorize base url: %w", err))
	}

	q := url.Values{}
	q.Set("client_id", req.GoogleOAuthClientID)
	q.Set("redirect_uri", req.RedirectURI)
	q.Set("state", req.State)
	q.Set("response_type", "code")
	q.Set("scope", "https://www.googleapis.com/auth/userinfo.email")

	u.RawQuery = q.Encode()
	return u.String()
}

type RedeemCodeRequest struct {
	GoogleOAuthClientID     string
	GoogleOAuthClientSecret string
	RedirectURI             string
	Code                    string
}

type RedeemCodeResponse struct {
	GoogleUserID       string
	Email              string
	EmailVerified      bool
	GoogleHostedDomain string
}

func (c *Client) RedeemCode(ctx context.Context, req *RedeemCodeRequest) (*RedeemCodeResponse, error) {
	accessToken, err := c.token(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("redeem code: %w", err)
	}

	userInfo, err := c.userinfo(ctx, accessToken)
	if err != nil {
		return nil, fmt.Errorf("userinfo: %w", err)
	}

	if err := c.revoke(ctx, accessToken); err != nil {
		return nil, fmt.Errorf("revoke: %w", err)
	}

	return &RedeemCodeResponse{
		GoogleUserID:       userInfo.Sub,
		Email:              userInfo.Email,
		EmailVerified:      userInfo.EmailVerified,
		GoogleHostedDomain: userInfo.HD,
	}, nil
}

func (c *Client) token(ctx context.Context, req *RedeemCodeRequest) (string, error) {
	body := url.Values{}
	body.Set("client_id", req.GoogleOAuthClientID)
	body.Set("client_secret", req.GoogleOAuthClientSecret)
	body.Set("code", req.Code)
	body.Set("redirect_uri", req.RedirectURI)
	body.Set("grant_type", "authorization_code")

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://oauth2.googleapis.com/token", strings.NewReader(body.Encode()))
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
		AccessToken string `json:"access_token"`
	}
	if err := json.Unmarshal(resBody, &data); err != nil {
		return "", fmt.Errorf("unmarshal response: %w", err)
	}

	return data.AccessToken, nil
}

type userinfoResponse struct {
	Sub           string `json:"sub"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	HD            string `json:"hd"`
}

func (c *Client) userinfo(ctx context.Context, accessToken string) (*userinfoResponse, error) {
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://openidconnect.googleapis.com/v1/userinfo", nil)
	if err != nil {
		return nil, fmt.Errorf("new http request: %w", err)
	}

	httpReq.Header.Add("Authorization", "Bearer "+accessToken)

	httpRes, err := c.HTTPClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("send http request: %w", err)
	}
	defer func() { _ = httpRes.Body.Close() }()

	if httpRes.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad response status code: %s", httpRes.Status)
	}

	resBody, err := io.ReadAll(httpRes.Body)
	if err != nil {
		return nil, fmt.Errorf("read body: %w", err)
	}

	var data userinfoResponse
	if err := json.Unmarshal(resBody, &data); err != nil {
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}

	return &data, nil
}

func (c *Client) revoke(ctx context.Context, accessToken string) error {
	body := url.Values{}
	body.Set("token", accessToken)

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://oauth2.googleapis.com/revoke", strings.NewReader(body.Encode()))
	if err != nil {
		return fmt.Errorf("new http request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	httpRes, err := c.HTTPClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("send http request: %w", err)
	}
	defer func() { _ = httpRes.Body.Close() }()

	if httpRes.StatusCode != http.StatusOK {
		return fmt.Errorf("bad response status code: %s", httpRes.Status)
	}

	return nil
}
