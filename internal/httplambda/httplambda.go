package httplambda

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// Handler converts an http.Handler into a lambda.Handler that supports
// APIGatewayV2 HTTP requests in BUFFERED mode.
func Handler(h http.Handler) lambda.Handler {
	return lambda.NewHandler(func(ctx context.Context, e events.LambdaFunctionURLRequest) (events.LambdaFunctionURLResponse, error) {
		req, err := httpRequest(ctx, e)
		if err != nil {
			return events.LambdaFunctionURLResponse{}, fmt.Errorf("http request from event: %w", err)
		}

		w := httptest.NewRecorder()
		h.ServeHTTP(w, req)
		res := httpResponseEvent(w)

		slog.Info("http_req_lambda", "req", req, "res", res)

		return res, nil
	})
}

func httpRequest(ctx context.Context, e events.LambdaFunctionURLRequest) (*http.Request, error) {
	u := url.URL{
		Scheme:   "https",
		Host:     e.Headers["Host"],
		Path:     e.RawPath,
		RawQuery: e.RawQueryString,
	}

	var body io.Reader
	if e.IsBase64Encoded {
		body = base64.NewDecoder(base64.StdEncoding, strings.NewReader(e.Body))
	} else {
		body = strings.NewReader(e.Body)
	}

	req, err := http.NewRequestWithContext(ctx, e.RequestContext.HTTP.Method, u.String(), body)
	if err != nil {
		return nil, fmt.Errorf("new request: %w", err)
	}

	for k, v := range e.Headers {
		req.Header.Add(k, v)
	}
	return req, nil
}

func httpResponseEvent(w *httptest.ResponseRecorder) events.LambdaFunctionURLResponse {
	res := w.Result()

	var cookies []string           // handled separately
	headers := map[string]string{} // comma-joined, because api gateway v2 does not support multivalued headers
	for k, vv := range res.Header {
		if k == "Set-Cookie" {
			cookies = append(cookies, vv...)
			continue
		}

		headers[k] = strings.Join(vv, ",")
	}

	body, _ := io.ReadAll(res.Body)
	return events.LambdaFunctionURLResponse{
		StatusCode:      res.StatusCode,
		Cookies:         cookies,
		Headers:         headers,
		Body:            base64.StdEncoding.EncodeToString(body),
		IsBase64Encoded: true,
	}
}
