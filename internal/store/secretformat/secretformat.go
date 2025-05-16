package secretformat

import (
	"github.com/tesseral-labs/tesseral/internal/prettysecret"
)

const alphabet = "0123456789abcdefghijklmnopqrstuvwxyz"

var (
	APIKeySecretToken = prettysecret.MustNewFormat("api_secret_", alphabet)
)

func MustNewFormat(prefix string) prettysecret.Format {
	return prettysecret.MustNewFormat(prefix, alphabet)
}
