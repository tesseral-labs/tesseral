package totp_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/tesseral-labs/tesseral/internal/totp"
)

func TestKey_Validate(t *testing.T) {
	// https://datatracker.ietf.org/doc/html/rfc6238#appendix-B
	k := totp.Key{Secret: []byte("12345678901234567890")}

	testCases := []struct {
		unix int64
		code string
	}{
		{59, "287082"},
		{1111111109, "081804"},
		{1111111111, "050471"},
		{1234567890, "005924"},
		{2000000000, "279037"},
		{20000000000, "353130"},
	}

	for _, tt := range testCases {
		t.Run(fmt.Sprintf("%d", tt.unix), func(t *testing.T) {
			assert.NoError(t, k.Validate(time.Unix(tt.unix, 0), tt.code))
			assert.Error(t, k.Validate(time.Unix(tt.unix, 0), "000000"))
		})
	}
}
