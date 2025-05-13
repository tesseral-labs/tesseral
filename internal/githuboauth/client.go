package githuboauth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
)

type Client struct {
	HTTPClient *http.Client
}

type GetAuthorizeURLRequest struct {
	GithubOAuthClientID string
	RedirectURI         string
	State               string
}

func GetAuthorizeURL(req *GetAuthorizeURLRequest) string {
	u, err := url.Parse("https://github.com/login/oauth/authorize")
	if err != nil {
		panic(fmt.Errorf("parse authorize base url: %w", err))
	}

	q := url.Values{}
	q.Set("client_id", req.GithubOAuthClientID)
	q.Set("redirect_uri", req.RedirectURI)
	q.Set("state", req.State)
	q.Set("scope", "read:user user:email")

	u.RawQuery = q.Encode()
	return u.String()
}

type RedeemCodeRequest struct {
	GithubOAuthClientID     string
	GithubOAuthClientSecret string
	RedirectURI             string
	Code                    string
}

type RedeemCodeResponse struct {
	GithubUserID      string
	Email             string
	EmailVerified     bool
	DisplayName       string
	ProfilePictureURL string
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

	email, verified, err := c.primaryEmail(ctx, accessToken)
	if err != nil {
		return nil, fmt.Errorf("primary email: %w", err)
	}

	return &RedeemCodeResponse{
		GithubUserID:      fmt.Sprint(userInfo.ID),
		Email:             email,
		EmailVerified:     verified,
		DisplayName:       userInfo.Name,
		ProfilePictureURL: userInfo.AvatarURL,
	}, nil
}

func (c *Client) token(ctx context.Context, req *RedeemCodeRequest) (string, error) {
	body := url.Values{}
	body.Set("client_id", req.GithubOAuthClientID)
	body.Set("client_secret", req.GithubOAuthClientSecret)
	body.Set("code", req.Code)
	body.Set("redirect_uri", req.RedirectURI)

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://github.com/login/oauth/access_token", strings.NewReader(body.Encode()))
	if err != nil {
		return "", fmt.Errorf("new http request: %w", err)
	}
	httpReq.Header.Set("Accept", "application/json")
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

type githubUserInfo struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	AvatarURL string `json:"avatar_url"`
}

func (c *Client) userinfo(ctx context.Context, accessToken string) (*githubUserInfo, error) {
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.github.com/user", nil)
	if err != nil {
		return nil, fmt.Errorf("new http request: %w", err)
	}
	httpReq.Header.Set("Authorization", "Bearer "+accessToken)
	httpReq.Header.Set("Accept", "application/vnd.github+json")

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

	slog.InfoContext(ctx, "github_userinfo", "response_body", string(resBody))

	var data githubUserInfo
	if err := json.Unmarshal(resBody, &data); err != nil {
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}

	return &data, nil
}

func (c *Client) primaryEmail(ctx context.Context, accessToken string) (string, bool, error) {
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.github.com/user/emails", nil)
	if err != nil {
		return "", false, fmt.Errorf("new http request: %w", err)
	}
	httpReq.Header.Set("Authorization", "Bearer "+accessToken)
	httpReq.Header.Set("Accept", "application/vnd.github+json")

	httpRes, err := c.HTTPClient.Do(httpReq)
	if err != nil {
		return "", false, fmt.Errorf("send http request: %w", err)
	}
	defer func() { _ = httpRes.Body.Close() }()

	if httpRes.StatusCode != http.StatusOK {
		return "", false, fmt.Errorf("bad response status code: %s", httpRes.Status)
	}

	resBody, err := io.ReadAll(httpRes.Body)
	if err != nil {
		return "", false, fmt.Errorf("read body: %w", err)
	}

	var emails []struct {
		Email    string `json:"email"`
		Primary  bool   `json:"primary"`
		Verified bool   `json:"verified"`
	}
	if err := json.Unmarshal(resBody, &emails); err != nil {
		return "", false, fmt.Errorf("unmarshal response: %w", err)
	}

	for _, e := range emails {
		if e.Primary {
			return e.Email, e.Verified, nil
		}
	}

	return "", false, fmt.Errorf("no verified primary email found")
}
