package prettysecret_test

import (
	"crypto/rand"
	"fmt"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/tesseral-labs/tesseral/internal/prettysecret"
)

func TestFormatParseErrors(t *testing.T) {
	format := func() (string, int) {
		prefix := "test_"
		length := prettysecret.SecretLen(prefix)
		return prefix, length
	}
	prefix, expectedLen := format()

	testCases := []struct {
		Name    string
		Input   string
		Prefix  string
		WantErr string
	}{
		{
			Name:    "wrong prefix",
			Input:   "badprefix_" + strings.Repeat("0", expectedLen-len("badprefix_")),
			Prefix:  prefix,
			WantErr: fmt.Sprintf("%q does not have expected prefix %q", "badprefix_"+strings.Repeat("0", expectedLen-len("badprefix_")), prefix),
		},
		{
			Name:    "wrong length",
			Input:   prefix + strings.Repeat("0", expectedLen-len(prefix)-1),
			Prefix:  prefix,
			WantErr: fmt.Sprintf("%q does not have expected length %v", prefix+strings.Repeat("0", expectedLen-len(prefix)-1), expectedLen),
		},
		{
			Name:    "invalid character",
			Input:   prefix + strings.Repeat("0", expectedLen-len(prefix)-1) + "!",
			Prefix:  prefix,
			WantErr: fmt.Sprintf("%q contains illegal char at position %v", prefix+strings.Repeat("0", expectedLen-len(prefix)-1)+"!", expectedLen-1),
		},
	}

	for _, tt := range testCases {
		t.Run(tt.Name, func(t *testing.T) {
			_, err := prettysecret.Parse(tt.Prefix, tt.Input)
			if err == nil {
				t.Fatalf("expected error, got nil")
			}
			if d := cmp.Diff(tt.WantErr, err.Error()); d != "" {
				t.Fatalf("unexpected error (-want +got):\n%s", d)
			}
		})
	}
}

func TestFormatParseRoundTrip(t *testing.T) {
	testCases := []struct {
		Name   string
		Secret [35]byte
	}{
		{
			Name:   "all zero bytes",
			Secret: [35]byte{},
		},
		{
			Name: "all 0xff bytes",
			Secret: func() (s [35]byte) {
				for i := range s {
					s[i] = 0xff
				}
				return
			}(),
		},
		{
			Name: "random secret",
			Secret: func() (s [35]byte) {
				if _, err := rand.Read(s[:]); err != nil {
					panic(fmt.Errorf("failed to generate random bytes: %w", err))
				}
				return
			}(),
		},
	}

	for _, tt := range testCases {
		t.Run(tt.Name, func(t *testing.T) {
			prefix := "test_"
			encoded := prettysecret.Format(prefix, tt.Secret)

			got, err := prettysecret.Parse(prefix, encoded)
			if err != nil {
				t.Fatalf("unexpected error from Parse: %v", err)
			}

			if d := cmp.Diff(tt.Secret[:], got[:]); d != "" {
				t.Fatalf("Parse(Format(secret)) mismatch (-want +got):\n%s", d)
			}
		})
	}
}
