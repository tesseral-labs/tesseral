package secretformat

import (
	"fmt"
	"math"
	"math/big"
	"strings"
)

// Format converts between secrets of an arbitrary length and their pretty formats.
//
// A Format keeps track of a prefix and an alphabet. The alphabet is used as the set of digits in a base-len(alphabet)
// numeral system. The prefix is prepended to all formatted strings, and expected in all parsed strings.
//
// For example, here a few common alphabets:
//
//	binary: 01
//	hexadecimal: 0123456789abcdef
//	raw unpadded base64: ABCDEFGHJIKLMNOPQRSTUVWYZabcdefghijklmnopqrstuvwxyz0123456789+/
type Format struct {
	prefix   string
	alphabet string
}

// NewFormat creates a Format with the given prefix and alphabet.
//
// For example, to create a Format where UUIDs are prefixed with "invoice_" and are base36-encoded, use:
//
//	NewFormat("invoice_", "0123456789abcdefghijklmnopqrstuvwxyz")
//
// NewFormat returns an error if len(alphabet) < 2 or if alphabet contains duplicate characters. Any prefix, including
// an empty prefix, is valid.
func NewFormat(prefix, alphabet string) (Format, error) {
	if len(alphabet) < 2 {
		return Format{}, fmt.Errorf("alphabet must have len >= 2: %q", alphabet)
	}
	seen := map[byte]struct{}{}
	for i := 0; i < len(alphabet); i++ {
		if _, ok := seen[alphabet[i]]; ok {
			return Format{}, fmt.Errorf("alphabet must not contain duplicate chars: %q", alphabet)
		}
		seen[alphabet[i]] = struct{}{}
	}
	return Format{prefix: prefix, alphabet: alphabet}, nil
}

// MustNewFormat is like NewFormat but panics if alphabet isn't valid.
func MustNewFormat(prefix, alphabet string) Format {
	format, err := NewFormat(prefix, alphabet)
	if err != nil {
		panic(fmt.Errorf("NewFormat(%q, %q) err: %w", prefix, alphabet, err))
	}
	return format
}

// Format converts a UUID to a pretty string.
func (f *Format) Format(uuid [16]byte) string {
	s := make([]byte, f.len())
	copy(s, f.prefix)
	for i := len(f.prefix); i < f.len(); i++ {
		s[i] = f.alphabet[0]
	}

	var n big.Int
	n.SetBytes(uuid[:])

	b := big.NewInt(int64(len(f.alphabet)))
	i := len(s) - 1
	for n.BitLen() > 0 {
		var r big.Int
		n.QuoRem(&n, b, &r)

		s[i] = f.alphabet[r.Int64()]
		i--
	}

	return string(s)
}

// Parse converts a pretty string to a UUID.
func (f *Format) Parse(s string) ([16]byte, error) {
	if !strings.HasPrefix(s, f.prefix) {
		return [16]byte{}, fmt.Errorf("%q does not have expected prefix %q", s, f.prefix)
	}

	if len(s) != f.len() {
		return [16]byte{}, fmt.Errorf("%q does not have expected length %v", s, f.len())
	}

	var n big.Int
	for i := len(f.prefix); i < len(s); i++ {
		d := strings.IndexByte(f.alphabet, s[i])
		if d == -1 {
			return [16]byte{}, fmt.Errorf("%q contains illegal char at position %v", s, i)
		}

		n.Mul(&n, big.NewInt(int64(len(f.alphabet))))
		n.Add(&n, big.NewInt(int64(d)))
	}

	b := n.Bytes()
	if len(b) > 16 {
		panic(fmt.Errorf("prettyuuid invariant failure"))
	}

	var out [16]byte
	copy(out[16-len(b):], b)
	return out, nil
}

func (f *Format) len() int {
	// digits required is ceil(log_{radix}(max_uuid))
	// log_{radix}(max_uuid) = log2(max_uuid) / log2(radix) = 128 / log2(radix)
	return len(f.prefix) + int(math.Ceil(128.0/math.Log2(float64(len(f.alphabet)))))
}
