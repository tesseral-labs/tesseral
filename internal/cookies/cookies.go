package cookies

import (
	"fmt"
	"net/http"
	"time"

	"connectrpc.com/connect"
	"github.com/google/uuid"
	"github.com/tesseral-labs/tesseral/internal/store/idformat"
)

func ExpiredRefreshToken(projectID uuid.UUID) string {
	return newCookie("refresh_token", -1*time.Second, projectID, "", true)
}

func ExpiredAccessToken(projectID uuid.UUID) string {
	return newCookie("access_token", -1*time.Second, projectID, "", false)
}

func ExpiredIntermediateAccessToken(projectID uuid.UUID) string {
	return newCookie("intermediate_access_token", -1*time.Second, projectID, "", true)
}

func NewRefreshToken(projectID uuid.UUID, value string) string {
	return newCookie("refresh_token", time.Hour*24*365, projectID, value, true)
}

func NewAccessToken(projectID uuid.UUID, value string) string {
	return newCookie("access_token", 5*time.Minute, projectID, value, false)
}

func NewIntermediateAccessToken(projectID uuid.UUID, value string) string {
	return newCookie("intermediate_access_token", 15*time.Minute, projectID, value, true)
}

func newCookie(name string, maxAge time.Duration, projectID uuid.UUID, value string, httpOnly bool) string {
	c := http.Cookie{
		Name:     fmt.Sprintf("tesseral_%s_%s", idformat.Project.Format(projectID), name),
		Value:    value,
		MaxAge:   int(maxAge.Seconds()),
		Path:     "/",
		SameSite: http.SameSiteNoneMode,
		Secure:   true,
		HttpOnly: httpOnly,
	}
	return c.String()
}

func GetRefreshToken(projectID uuid.UUID, req connect.AnyRequest) (string, error) {
	return getCookie("refresh_token", projectID, req)
}

func GetAccessToken(projectID uuid.UUID, req connect.AnyRequest) (string, error) {
	return getCookie("access_token", projectID, req)
}

func GetIntermediateAccessToken(projectID uuid.UUID, req connect.AnyRequest) (string, error) {
	return getCookie("intermediate_access_token", projectID, req)
}

func getCookie(name string, projectID uuid.UUID, req connect.AnyRequest) (string, error) {
	cookieName := fmt.Sprintf("tesseral_%s_%s", idformat.Project.Format(projectID), name)

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
