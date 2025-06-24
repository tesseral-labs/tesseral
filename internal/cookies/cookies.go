package cookies

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"connectrpc.com/connect"
	"github.com/google/uuid"
	commonstore "github.com/tesseral-labs/tesseral/internal/common/store"
	"github.com/tesseral-labs/tesseral/internal/store/idformat"
)

type Cookier struct {
	Store *commonstore.Store
}

func (c *Cookier) GetRefreshToken(projectID uuid.UUID, req connect.AnyRequest) (string, error) {
	return c.getCookie("refresh_token", projectID, req)
}

func (c *Cookier) GetAccessToken(projectID uuid.UUID, req connect.AnyRequest) (string, error) {
	return c.getCookie("access_token", projectID, req)
}

func (c *Cookier) GetIntermediateAccessToken(projectID uuid.UUID, req connect.AnyRequest) (string, error) {
	return c.getCookie("intermediate_access_token", projectID, req)
}

func (c *Cookier) GetOIDCIntermediateSessionToken(projectID uuid.UUID, req *http.Request) (string, error) {
	cookies := req.CookiesNamed(c.cookieName("oidc_intermediate_session", projectID))
	if len(cookies) != 1 {
		return "", fmt.Errorf("expected exactly one oidc_intermediate_session cookie, got %d", len(cookies))
	}
	return cookies[0].Value, nil
}

func (c *Cookier) getCookie(name string, projectID uuid.UUID, req connect.AnyRequest) (string, error) {
	cookieName := c.cookieName(name, projectID)

	var value string
	for _, h := range req.Header().Values("Cookie") {
		cookies, err := http.ParseCookie(h)
		if err != nil {
			return "", fmt.Errorf("parse cookie: %w", err)
		}

		for _, c := range cookies {
			if c.Name != cookieName {
				continue
			}
			value = c.Value
		}
	}

	return value, nil
}

func (c *Cookier) ExpiredRefreshToken(ctx context.Context, projectID uuid.UUID) (string, error) {
	return c.expiredCookie(ctx, "refresh_token", projectID)
}

func (c *Cookier) ExpiredAccessToken(ctx context.Context, projectID uuid.UUID) (string, error) {
	return c.expiredCookie(ctx, "access_token", projectID)
}

func (c *Cookier) ExpiredIntermediateAccessToken(ctx context.Context, projectID uuid.UUID) (string, error) {
	return c.expiredCookie(ctx, "intermediate_access_token", projectID)
}

func (c *Cookier) ExpiredOIDCIntermediateSessionToken(ctx context.Context, projectID uuid.UUID) (string, error) {
	return c.expiredCookie(ctx, "oidc_intermediate_session", projectID)
}

func (c *Cookier) expiredCookie(ctx context.Context, name string, projectID uuid.UUID) (string, error) {
	return c.newCookie(ctx, name, projectID, -1*time.Second, "", false)
}

func (c *Cookier) NewRefreshToken(ctx context.Context, projectID uuid.UUID, value string) (string, error) {
	return c.newCookie(ctx, "refresh_token", projectID, time.Hour*24*365, value, true)
}

func (c *Cookier) NewAccessToken(ctx context.Context, projectID uuid.UUID, value string) (string, error) {
	return c.newCookie(ctx, "access_token", projectID, 5*time.Minute, value, false)
}

func (c *Cookier) NewIntermediateAccessToken(ctx context.Context, projectID uuid.UUID, value string) (string, error) {
	return c.newCookie(ctx, "intermediate_access_token", projectID, 15*time.Minute, value, true)
}

func (c *Cookier) NewOIDCIntermediateSessionToken(ctx context.Context, projectID uuid.UUID, value string) (string, error) {
	return c.newCookie(ctx, "oidc_intermediate_session", projectID, 15*time.Minute, value, true)
}

func (c *Cookier) newCookie(ctx context.Context, name string, projectID uuid.UUID, maxAge time.Duration, value string, httpOnly bool) (string, error) {
	cookieDomain, err := c.Store.GetProjectCookieDomain(ctx, projectID)
	if err != nil {
		return "", fmt.Errorf("get project cookie domain: %w", err)
	}

	cookie := http.Cookie{
		Name:     c.cookieName(name, projectID),
		Value:    value,
		MaxAge:   int(maxAge.Seconds()),
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
		Secure:   true,
		HttpOnly: httpOnly,
		Domain:   cookieDomain,
	}
	return cookie.String(), nil
}

func (c *Cookier) cookieName(name string, projectID uuid.UUID) string {
	return fmt.Sprintf("tesseral_%s_%s", idformat.Project.Format(projectID), name)
}
