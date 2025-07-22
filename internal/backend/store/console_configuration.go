package store

import (
	"context"
	"fmt"

	backendv1 "github.com/tesseral-labs/tesseral/internal/backend/gen/tesseral/backend/v1"
	"github.com/tesseral-labs/tesseral/internal/common/apierror"
	"github.com/tesseral-labs/tesseral/internal/store/idformat"
)

func (s *Store) ConsoleGetConfiguration(ctx context.Context, req *backendv1.ConsoleGetConfigurationRequest) (*backendv1.ConsoleGetConfigurationResponse, error) {
	if err := validateIsConsoleSession(ctx); err != nil {
		return nil, fmt.Errorf("validate is console session: %w", err)
	}

	if s.consoleProjectID == nil {
		return nil, apierror.NewFailedPreconditionError("console project id is not set", fmt.Errorf("console project id is nil"))
	}

	// Implementation of fetching console configuration from the store
	// This is a placeholder for actual logic to retrieve the configuration
	return &backendv1.ConsoleGetConfigurationResponse{
		Configuration: &backendv1.ConsoleConfiguration{
			ConsoleProjectId: idformat.Project.Format(derefOrEmpty(s.consoleProjectID)),
		},
	}, nil
}
