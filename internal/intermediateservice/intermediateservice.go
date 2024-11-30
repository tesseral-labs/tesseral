package intermediateservice

import (
	"github.com/openauth/openauth/internal/gen/intermediate/v1/intermediatev1connect"
	"github.com/openauth/openauth/internal/store"
)

type IntermediateService struct {
	Store *store.Store
	intermediatev1connect.UnimplementedIntermediateServiceHandler
}
