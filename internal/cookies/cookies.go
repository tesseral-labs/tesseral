package cookies

import (
	"context"
	"fmt"
	"strings"

	"connectrpc.com/connect"
	"github.com/google/uuid"
	"github.com/openauth/openauth/internal/projectid"
	"github.com/openauth/openauth/internal/store/idformat"
)

var errCookieNotFound = fmt.Errorf("cookie not found")

func BuildCookie(projectID uuid.UUID, cookieType string, accessToken string, secure bool) string {
	secureStr := ""
	if secure {
		secureStr = "Secure;"
	}

	maxAge := 60 * 60 * 24 * 7 // one week

	// TODO: Once domains are sorted out, we'll need to set the `Domain` attribute on the cookie.
	return fmt.Sprintf("tesseral_%s_%s=%s;SameSite=Lax;HttpOnly;MaxAge=%d;%s", idformat.Project.Format(projectID), cookieType, accessToken, maxAge, secureStr)
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
