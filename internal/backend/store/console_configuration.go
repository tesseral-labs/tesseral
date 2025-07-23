package store

import (
	"context"

	backendv1 "github.com/tesseral-labs/tesseral/internal/backend/gen/tesseral/backend/v1"
	"github.com/tesseral-labs/tesseral/internal/store/idformat"
)

func (s *Store) ConsoleGetConfiguration(ctx context.Context, req *backendv1.ConsoleGetConfigurationRequest) (*backendv1.ConsoleGetConfigurationResponse, error) {
	return &backendv1.ConsoleGetConfigurationResponse{
		Configuration: &backendv1.ConsoleConfiguration{
			ConsoleProjectId: idformat.Project.Format(*s.consoleProjectID),
		},
	}, nil
}
