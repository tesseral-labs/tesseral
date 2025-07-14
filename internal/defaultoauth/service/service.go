package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"connectrpc.com/connect"
	"github.com/tesseral-labs/tesseral/internal/defaultoauth/store"
	"google.golang.org/protobuf/encoding/protojson"
)

type Service struct {
	Store *store.Store
}

func (s *Service) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.Handle("GET /google-oauth-callback", withErr(s.googleOAuthCallback))
	mux.Handle("GET /microsoft-oauth-callback", withErr(s.microsoftOAuthCallback))
	mux.Handle("GET /github-oauth-callback", withErr(s.gitHubOAuthCallback))
	return mux
}

func withErr(f func(w http.ResponseWriter, r *http.Request) error) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			// For rough consistency with the frontend/backend/intermediate
			// APIs, we emulate connect-go's handling of connect.Error.
			var connectErr *connect.Error
			if errors.As(err, &connectErr) {
				writeConnectError(w, connectErr)
				return
			}

			http.Error(w, "", http.StatusInternalServerError)
			panic(err)
		}
	})
}

// Adapted from https://github.com/connectrpc/connect-go/blob/d7c0966751650b41a9f1794513592e81b9beed45/protocol_connect.go

func writeConnectError(w http.ResponseWriter, connectErr *connect.Error) {
	wireError := struct {
		Code    int               `json:"code"`
		Message string            `json:"message"`
		Details []json.RawMessage `json:"details"`
	}{
		Code:    connectCodeToHTTP(connectErr.Code()),
		Message: connectErr.Message(),
	}

	for _, detail := range connectErr.Details() {
		value, _ := detail.Value()
		detailJSON, _ := protojson.Marshal(value)
		wireError.Details = append(wireError.Details, detailJSON)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(wireError.Code)
	if err := json.NewEncoder(w).Encode(wireError); err != nil {
		panic(fmt.Errorf("encode connect error: %w", err))
	}
}

func connectCodeToHTTP(code connect.Code) int {
	// Return literals rather than named constants from the HTTP package to make
	// it easier to compare this function to the Connect specification.
	switch code {
	case connect.CodeCanceled:
		return 499
	case connect.CodeUnknown:
		return 500
	case connect.CodeInvalidArgument:
		return 400
	case connect.CodeDeadlineExceeded:
		return 504
	case connect.CodeNotFound:
		return 404
	case connect.CodeAlreadyExists:
		return 409
	case connect.CodePermissionDenied:
		return 403
	case connect.CodeResourceExhausted:
		return 429
	case connect.CodeFailedPrecondition:
		return 400
	case connect.CodeAborted:
		return 409
	case connect.CodeOutOfRange:
		return 400
	case connect.CodeUnimplemented:
		return 501
	case connect.CodeInternal:
		return 500
	case connect.CodeUnavailable:
		return 503
	case connect.CodeDataLoss:
		return 500
	case connect.CodeUnauthenticated:
		return 401
	default:
		return 500 // same as CodeUnknown
	}
}
