package intermediateservice

import (
	"github.com/openauth-dev/openauth/internal/gen/intermediate/v1/intermediatev1connect"
	"github.com/openauth-dev/openauth/internal/store"
)

type IntermediateService struct {
	Store *store.Store
	intermediatev1connect.UnimplementedIntermediateServiceHandler
}
