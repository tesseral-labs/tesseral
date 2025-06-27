package totp

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base32"
	"encoding/binary"
	"fmt"
	"net/url"
	"time"
)

type Key struct {
	Secret []byte
}

func (k *Key) OTPAuthURI(issuer, user string) string {
	uri := url.URL{
		Scheme: "otpauth",
		Host:   "totp",
		Path:   fmt.Sprintf("%s:%s", issuer, user),
		RawQuery: url.Values{
			"secret": {base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(k.Secret)},
			"digits": {"6"},
			"period": {"30"},
		}.Encode(),
	}
	return uri.String()
}

func (k *Key) Validate(now time.Time, code string) error {
	if code != k.gen(now) {
		return fmt.Errorf("incorrect totp code")
	}
	return nil
}

func (k *Key) gen(now time.Time) string {
	counter := now.Unix() / 30

	mac := hmac.New(sha1.New, k.Secret)
	_ = binary.Write(mac, binary.BigEndian, counter)
	sum := mac.Sum(nil)

	offset := sum[len(sum)-1] & 0xf
	value := int64(((int(sum[offset]) & 0x7f) << 24) |
		((int(sum[offset+1] & 0xff)) << 16) |
		((int(sum[offset+2] & 0xff)) << 8) |
		(int(sum[offset+3]) & 0xff))

	return fmt.Sprintf("%06d", value%1_000_000)
}
