package hibp

import (
	"bufio"
	"context"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type Client struct {
	HTTPClient *http.Client
}

func (c *Client) Pwned(ctx context.Context, password string) (bool, error) {
	sha1Bytes := sha1.Sum([]byte(password))
	sha1Hex := strings.ToUpper(hex.EncodeToString(sha1Bytes[:]))
	sha1HexPrefix, sha1HexSuffix := sha1Hex[:5], sha1Hex[5:]

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("https://api.pwnedpasswords.com/range/%s", sha1HexPrefix), nil)
	if err != nil {
		return false, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Add-Padding", "true")

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return false, fmt.Errorf("do request: %w", err)
	}
	defer func() { _ = res.Body.Close() }()

	if res.StatusCode != http.StatusOK {
		return false, fmt.Errorf("bad response status code: %d", res.StatusCode)
	}

	s := bufio.NewScanner(res.Body)
	var match bool
	for s.Scan() {
		line := s.Text()
		lineMatch, err := checkLineMatch(sha1HexSuffix, line)
		if err != nil {
			return false, err
		}

		if lineMatch {
			match = true
			break
		}
	}

	if err := s.Err(); err != nil {
		return false, fmt.Errorf("scan response: %w", err)
	}

	return match, nil
}

func checkLineMatch(sha1HexSuffix, line string) (bool, error) {
	parts := strings.SplitN(line, ":", 2)
	if len(parts) != 2 {
		return false, fmt.Errorf("bad response line: %s", line)
	}

	suffix := parts[0]
	count, err := strconv.Atoi(parts[1])
	if err != nil {
		return false, fmt.Errorf("parse line count: %s: %w", line, err)
	}

	// match, and it's not a padding result
	if suffix == sha1HexSuffix && count > 0 {
		return true, nil
	}

	return false, nil
}
