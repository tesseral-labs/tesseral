package symmetrickeys

import (
	"crypto/rand"
	"encoding/base64"
)

func GenerateSymmetricKey() (string, error) {
    key := make([]byte, 32) // 256-bit key for HS256
    _, err := rand.Read(key)

    if err != nil {
      return "", err
    }

    return base64.URLEncoding.EncodeToString(key), nil
}