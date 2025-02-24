package cloudflaredoh

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type Client struct {
	HTTPClient *http.Client
}

type DNSQueryRequest struct {
	Name string
	Type string
}

type DNSQueryResponse struct {
	Answer []DNSQueryResponseAnswer
}

type DNSQueryResponseAnswer struct {
	Name string
	Type int32
	Data string
	TTL  uint32
}

func (c *Client) DNSQuery(ctx context.Context, req *DNSQueryRequest) (*DNSQueryResponse, error) {
	reqURL, err := url.Parse("https://cloudflare-dns.com/dns-query")
	if err != nil {
		return nil, fmt.Errorf("parse url: %w", err)
	}

	reqURL.RawQuery = url.Values{
		"name": {req.Name},
		"type": {req.Type},
	}.Encode()

	httpReq, err := http.NewRequestWithContext(ctx, "GET", reqURL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("new http request: %w", err)
	}

	httpReq.Header.Set("Accept", "application/dns-json")

	httpRes, err := c.HTTPClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("send http request: %w", err)
	}
	defer func() { _ = httpRes.Body.Close() }()

	if httpRes.StatusCode != http.StatusOK {
		buf, _ := io.ReadAll(httpRes.Body)
		fmt.Println(string(buf))
		return nil, fmt.Errorf("bad response status code: %s", httpRes.Status)
	}

	resBody, err := io.ReadAll(httpRes.Body)
	if err != nil {
		return nil, fmt.Errorf("read body: %w", err)
	}

	var data DNSQueryResponse
	if err := json.Unmarshal(resBody, &data); err != nil {
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}

	return &data, nil
}
