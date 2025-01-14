package store

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/openauth/openauth/internal/backend/authn"
	backendv1 "github.com/openauth/openauth/internal/backend/gen/openauth/backend/v1"
	"github.com/openauth/openauth/internal/backend/store/queries"
	"github.com/openauth/openauth/internal/store/idformat"
)

func (s *Store) CreateProjectRedirectURI(ctx context.Context, req *backendv1.CreateProjectRedirectURIRequest) (*backendv1.CreateProjectRedirectURIResponse, error) {
	_, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	projectID := authn.ProjectID(ctx)

	qProjectRedirectURI, err := q.CreateProjectRedirectURI(ctx, queries.CreateProjectRedirectURIParams{
		ID:        uuid.New(),
		ProjectID: projectID,
		Uri:       req.ProjectRedirectUri.Uri,
	})
	if err != nil {
		return nil, fmt.Errorf("create project redirect uri: %w", err)
	}

	if err := commit(); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}

	return &backendv1.CreateProjectRedirectURIResponse{
		ProjectRedirectUri: parseProjectRedirectURI(qProjectRedirectURI),
	}, nil
}

func (s *Store) DeleteProjectRedirectURI(ctx context.Context, req *backendv1.DeleteProjectRedirectURIRequest) (*backendv1.DeleteProjectRedirectURIResponse, error) {
	_, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	projectRedirectUriID, err := idformat.ProjectRedirectURI.Parse(req.Id)
	if err != nil {
		return nil, fmt.Errorf("parse project redirect uri id: %w", err)
	}

	err = q.DeleteProjectRedirectURI(ctx, queries.DeleteProjectRedirectURIParams{
		ProjectID: authn.ProjectID(ctx),
		ID:        projectRedirectUriID,
	})
	if err != nil {
		return nil, fmt.Errorf("delete project redirect uri: %w", err)
	}

	if err := commit(); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}

	return &backendv1.DeleteProjectRedirectURIResponse{}, nil
}

func (s *Store) GetProjectRedirectURI(ctx context.Context, req *backendv1.GetProjectRedirectURIRequest) (*backendv1.GetProjectRedirectURIResponse, error) {
	_, q, _, _, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}

	projectRedirectUriID, err := idformat.ProjectRedirectURI.Parse(req.Id)
	if err != nil {
		return nil, fmt.Errorf("parse project redirect uri id: %w", err)
	}

	qProjectRedirectURI, err := q.GetProjectRedirectURI(ctx, queries.GetProjectRedirectURIParams{
		ProjectID: authn.ProjectID(ctx),
		ID:        projectRedirectUriID,
	})
	if err != nil {
		return nil, fmt.Errorf("get project redirect uri: %w", err)
	}

	return &backendv1.GetProjectRedirectURIResponse{
		ProjectRedirectUri: parseProjectRedirectURI(qProjectRedirectURI),
	}, nil
}

func (s *Store) ListProjectRedirectURIs(ctx context.Context, req *backendv1.ListProjectRedirectURIsRequest) (*backendv1.ListProjectRedirectURIsResponse, error) {
	_, q, _, _, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}

	qProjectRedirectURIs, err := q.ListProjectRedirectURIs(ctx, authn.ProjectID(ctx))
	if err != nil {
		return nil, fmt.Errorf("list project redirect uris: %w", err)
	}

	projectRedirectURIs := make([]*backendv1.ProjectRedirectURI, len(qProjectRedirectURIs))
	for i, pru := range qProjectRedirectURIs {
		projectRedirectURIs[i] = parseProjectRedirectURI(pru)
	}

	return &backendv1.ListProjectRedirectURIsResponse{
		ProjectRedirectUris: projectRedirectURIs,
	}, nil
}

func (s *Store) UpdateProjectRedirectURI(ctx context.Context, req *backendv1.UpdateProjectRedirectURIRequest) (*backendv1.UpdateProjectRedirectURIResponse, error) {
	_, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	projectRedirectUriID, err := idformat.ProjectRedirectURI.Parse(req.Id)
	if err != nil {
		return nil, fmt.Errorf("parse project redirect uri id: %w", err)
	}

	qProjectRedirectURI, err := q.UpdateProjectRedirectURI(ctx, queries.UpdateProjectRedirectURIParams{
		ProjectID: authn.ProjectID(ctx),
		ID:        projectRedirectUriID,
		Uri:       req.ProjectRedirectUri.Uri,
		IsPrimary: req.ProjectRedirectUri.IsPrimary,
	})
	if err != nil {
		return nil, fmt.Errorf("update project redirect uri: %w", err)
	}

	if err := commit(); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}

	return &backendv1.UpdateProjectRedirectURIResponse{
		ProjectRedirectUri: parseProjectRedirectURI(qProjectRedirectURI),
	}, nil
}

func parseProjectRedirectURI(pru queries.ProjectRedirectUri) *backendv1.ProjectRedirectURI {
	return &backendv1.ProjectRedirectURI{
		Id:        pru.ID.String(),
		ProjectId: idformat.Project.Format(pru.ProjectID),
		IsPrimary: pru.IsPrimary,
		Uri:       pru.Uri,
	}
}
