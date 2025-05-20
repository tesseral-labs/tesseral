package prettysecret

import (
	"fmt"
	"math"
	"math/big"
	"strings"
)

const alphabet = "0123456789abcdefghijklmnopqrstuvwxyz"

// Format converts a secret to a pretty string.
func Format(prefix string, secret [35]byte) string {
	s := make([]byte, secret_len(prefix))
	copy(s, prefix)
	for i := len(prefix); i < secret_len(prefix); i++ {
		s[i] = alphabet[0]
	}

	var n big.Int
	n.SetBytes(secret[:])

	b := big.NewInt(int64(len(alphabet)))
	i := len(s) - 1
	for n.BitLen() > 0 {
		var r big.Int
		n.QuoRem(&n, b, &r)

		s[i] = alphabet[r.Int64()]
		i--
	}

	return string(s)
}

// Parse converts a pretty string to a [35]byte.
func Parse(prefix string, s string) ([35]byte, error) {
	if !strings.HasPrefix(s, prefix) {
		return [35]byte{}, fmt.Errorf("%q does not have expected prefix %q", s, prefix)
	}

	if len(s) != secret_len(prefix) {
		return [35]byte{}, fmt.Errorf("%q does not have expected length %v", s, secret_len(prefix))
	}

	var n big.Int
	for i := len(prefix); i < len(s); i++ {
		d := strings.IndexByte(alphabet, s[i])
		if d == -1 {
			return [35]byte{}, fmt.Errorf("%q contains illegal char at position %v", s, i)
		}

		n.Mul(&n, big.NewInt(int64(len(alphabet))))
		n.Add(&n, big.NewInt(int64(d)))
	}

	b := n.Bytes()
	if len(b) > 35 {
		panic(fmt.Errorf("prettsecret invariant failure: %d, %d", len(b), 35))
	}

	var out [35]byte
	copy(out[35-len(b):], b)
	return out, nil
}

func secret_len(prefix string) int {
	// digits required is ceil(log_{radix}(max_uuid))
	// log_{radix}(max_uuid) = log2(max_uuid) / log2(radix) = 128 / log2(radix)
	return len(prefix) + int(math.Ceil(280.0/math.Log2(float64(len(alphabet)))))
}

// SecretLen is a test-accessible version of the internal secret_len function.
func SecretLen(prefix string) int {
	return secret_len(prefix)
}
