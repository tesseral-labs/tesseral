package passkey_test

import (
	"fmt"
	"testing"

	"github.com/openauth/openauth/internal/passkey"
)

func TestParse(t *testing.T) {
	c, err := passkey.Parse(&passkey.ParseRequest{
		RPID:              "localhost",
		AttestationObject: "o2NmbXRkbm9uZWdhdHRTdG10oGhhdXRoRGF0YViYSZYN5YgOjGh0NBcPZHZgW4_krrmihjLHmVzzuoMdl2NdAAAAAPv8MAcVTk7MjAtuAgVX170AFFOOoV7GrKnKnccDm-8m0dTm_yQFpQECAyYgASFYILSxMv1cg3WN5GkouhyLJOXrIbBgSi9yAjI_QrC-IhAIIlggTDqhDZTjsPzx-dq2lkfu2AiieZuIpPUpQtOBYrX8gJw",
	})
	if err != nil {
		panic(err)
	}

	fmt.Println(c.Verify(&passkey.VerifyRequest{
		RPID:              "localhost",
		Origin:            "http://localhost:3002",
		ChallengeSHA256:   []byte{132, 143, 185, 51, 158, 98, 37, 12, 23, 156, 66, 204, 255, 170, 216, 93, 168, 10, 69, 31, 108, 79, 71, 89, 15, 138, 213, 219, 29, 51, 128, 200},
		ClientDataJSON:    "eyJ0eXBlIjoid2ViYXV0aG4uZ2V0IiwiY2hhbGxlbmdlIjoiZVZCb2NkcnU3cm1VdzhJVloyVW1JdVB6cXp0NEx0VnZnU2JpcGdGOWRGQSIsIm9yaWdpbiI6Imh0dHA6Ly9sb2NhbGhvc3Q6MzAwMiIsImNyb3NzT3JpZ2luIjpmYWxzZSwib3RoZXJfa2V5c19jYW5fYmVfYWRkZWRfaGVyZSI6ImRvIG5vdCBjb21wYXJlIGNsaWVudERhdGFKU09OIGFnYWluc3QgYSB0ZW1wbGF0ZS4gU2VlIGh0dHBzOi8vZ29vLmdsL3lhYlBleCJ9",
		AuthenticatorData: "SZYN5YgOjGh0NBcPZHZgW4_krrmihjLHmVzzuoMdl2MdAAAAAA",
		Signature:         "MEQCICUsfxpP1H2YjKM3PUwdX6rlTcIkrSUtsggWnqyEHE2NAiAwvtKsHzJtzE9ITWTP4rvIvkYoGss3Dg_a3RNkoNoXSg",
	}))
}
