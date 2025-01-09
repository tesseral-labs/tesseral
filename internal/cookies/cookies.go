package cookies

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"connectrpc.com/connect"
	"github.com/openauth/openauth/internal/projectid"
	"github.com/openauth/openauth/internal/store/idformat"
)

var errCookieNotFound = fmt.Errorf("cookie not found")

func BuildCookie(ctx context.Context, req connect.AnyRequest, cookieType string, value string) string {
	projectID := projectid.ProjectID(ctx)
	secure := req.Spec().Schema == "https"

	maxAge := 60 * 60 * 24 * 7 // one week
	if cookieType == "intermediateAccessToken" {
		maxAge = 60 * 15 // 15 minutes
	}

	// TODO: Once domains are sorted out, we'll need to set the `Domain` attribute on the cookie.
	cookie := http.Cookie{
		HttpOnly: true,
		MaxAge:   maxAge,
		Name:     fmt.Sprintf("tesseral_%s_%s", idformat.Project.Format(projectID), cookieType),
		Path:     "/",
		SameSite: http.SameSiteNoneMode,
		Secure:   secure,
		Value:    value,
	}

	return cookie.String()
}

func GetCookie(ctx context.Context, req connect.AnyRequest, cookieType string) (string, error) {
	cookie := req.Header().Get("Cookie")
	if cookie == "" {
		return "", errCookieNotFound
	}

	projectID := projectid.ProjectID(ctx)
	cookieName := fmt.Sprintf("tesseral_%s_%s", idformat.Project.Format(projectID), cookieType)

	value, found := extractCookieValue(cookie, cookieName)
	if !found {
		return "", errCookieNotFound
	}

	return value, nil
}

func extractCookieValue(cookie string, cookieName string) (string, bool) {
	parts := strings.Split(cookie, ";")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if strings.HasPrefix(part, cookieName) {
			return strings.CutPrefix(part, cookieName+"=")
		}
	}

	return "", false
}
