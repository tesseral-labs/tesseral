package service

import (
	"github.com/tesseral-labs/tesseral/internal/backend/gen/tesseral/backend/v1/backendv1connect"
	"github.com/tesseral-labs/tesseral/internal/backend/store"
)

type Service struct {
	Store *store.Store
	backendv1connect.UnimplementedBackendServiceHandler
}
