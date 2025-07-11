package trusteddomains_test

import (
	"testing"

	"github.com/tesseral-labs/tesseral/internal/common/trusteddomains"
)

func Test_IsTrustedDomain(t *testing.T) {
	trustedDomains := []string{
		"example.subdomain.com",
		"example.com",
		"localhost",
		"test.local:9999",
	}

	testCases := []struct {
		Name      string
		Origin    string
		IsTrusted bool
	}{
		{
			Name:      "Exact match",
			Origin:    "https://example.com",
			IsTrusted: true,
		},
		{
			Name:      "Subdomain match",
			Origin:    "https://sub.example.com",
			IsTrusted: true,
		},
		{
			Name:      "Parent domain",
			Origin:    "https://subdomain.com",
			IsTrusted: false,
		},
		{
			Name:      "Localhost with port",
			Origin:    "http://localhost:8080",
			IsTrusted: true,
		},
		{
			Name:      "Without port",
			Origin:    "https://test.local",
			IsTrusted: true,
		},
		{
			Name:      "With another port",
			Origin:    "https://test.local:1111",
			IsTrusted: true,
		},
		{
			Name:      "Non-trusted domain",
			Origin:    "https://nottrusted.com",
			IsTrusted: false,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.Name, func(t *testing.T) {
			isTrusted, err := trusteddomains.IsTrustedDomain(trustedDomains, tt.Origin)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if isTrusted != tt.IsTrusted {
				t.Errorf("expected %v, got %v", tt.IsTrusted, isTrusted)
			}
		})
	}
}
