package trusteddomains

import (
	"fmt"
	"net/url"
	"strings"
)

func IsTrustedDomain(trustedDomains []string, origin string) (bool, error) {
	originUrl, err := url.Parse(origin)
	if err != nil {
		return false, fmt.Errorf("parse origin: %w", err)
	}

	for _, trustedDomain := range trustedDomains {
		domain := trustedDomain
		domainParts := strings.Split(domain, ":")
		if len(domainParts) > 1 {
			domain = domainParts[0] // Remove port if present
		}

		trustedUrl := url.URL{Scheme: "https", Host: domain}

		if trustedUrl.Hostname() == "localhost" {
			trustedUrl = url.URL{Scheme: "http", Host: domain}
		}

		// Check in a port-agnostic way, so that
		// https://example.com:443 and https://example.com are considered the same.
		// Also allow subdomains of the trusted origin.
		if (trustedUrl.Hostname() == originUrl.Hostname() || strings.HasSuffix(originUrl.Hostname(), fmt.Sprintf(".%s", trustedUrl.Hostname())) || strings.HasPrefix(originUrl.Hostname(), fmt.Sprintf("%s:", trustedUrl.Hostname()))) && trustedUrl.Scheme == originUrl.Scheme {
			return true, nil
		}
	}

	return false, nil
}
