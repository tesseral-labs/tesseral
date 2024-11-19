package crypto

import "crypto/sha256"

func StringToSha256(s string) []byte {
	hash := sha256.New()
	hash.Write([]byte(s))

	return hash.Sum(nil)
}