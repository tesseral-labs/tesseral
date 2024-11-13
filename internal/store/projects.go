package store

import (
	"context"

	"github.com/google/uuid"
	backendv1 "github.com/openauth-dev/openauth/internal/gen/backend/v1"
	openauthv1 "github.com/openauth-dev/openauth/internal/gen/openauth/v1"
	"github.com/openauth-dev/openauth/internal/store/idformat"
	"github.com/openauth-dev/openauth/internal/store/queries"
)

func (s *Store) CreateProject(ctx context.Context, req *openauthv1.CreateProjectRequest) (*openauthv1.Project, error) {
	_, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	// Create a new project
	createdProject, err := q.CreateProject(ctx, queries.CreateProjectParams{
		ID: uuid.New(),
	})
	if err != nil {
		return nil, err
	}

	// Create the managing organization for the project
	// - this is required to create a relationship between the project 
	//   and the dogfooding project
	_, err = q.CreateOrganization(ctx, queries.CreateOrganizationParams{
		ID: uuid.New(),
		ProjectID: createdProject.ID,
		DisplayName: req.DisplayName,
	})
	if err != nil {
		return nil, err
	}

	// Update the project with the dogfooding project ID
	updatedProject, err := q.UpdateProjectOrganizationID(ctx, queries.UpdateProjectOrganizationIDParams{
		ID: createdProject.ID,
		OrganizationID: s.dogfoodProjectID,
	})
	if err != nil {
		return nil, err
	}

	// Commit all changes
	if err := commit(); err != nil {
		return nil, err
	}

	// Return the updated project
	return transformProject(&updatedProject), nil
}

func (s *Store) GetProject(ctx context.Context, req *openauthv1.ResourceIdRequest) (*openauthv1.Project, error) {
	id, err := idformat.Project.Parse(req.Id)
	if err != nil {
		return nil, err
	}

	_, q, _, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	project, err := q.GetProjectByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return transformProject(&project), nil
}

// TODO: Ensure that this function can only be called via a backend service reuqest
func (s *Store) ListProjects(ctx context.Context, req *backendv1.ListProjectsRequest) (*backendv1.ListProjectsResponse, error) {
	_, q, _, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	var startID uuid.UUID
	if err := s.pageEncoder.Unmarshal(req.PageToken, &startID); err != nil {
		return nil, err
	}

	limit := 10
	projectRecords, err := q.ListProjects(ctx, int32(limit + 1))
	if err != nil {
		return nil, err
	}

	projects := []*openauthv1.Project{}
	for _, project := range projectRecords {
		projects = append(projects, transformProject(&project))
	}

	var nextPageToken string
	if len(projects) == limit+1 {
		nextPageToken = s.pageEncoder.Marshal(projectRecords[limit].ID)
		projects = projects[:limit]
	}

	return &backendv1.ListProjectsResponse{
		Projects: projects,
		NextPageToken: nextPageToken,
	}, nil
}

func (s *Store) UpdateProject(ctx context.Context, req *openauthv1.Project) (*openauthv1.Project, error) {
	_, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	id, err := idformat.Project.Parse(req.Id)
	if err != nil {
		return nil, err
	}

	project, err := q.GetProjectByID(ctx, id)
	if err != nil {
		return nil, err
	}

	updates := queries.UpdateProjectParams{
		ID: project.ID,
	}

	// Conditionally configure Google OAuth
	if req.GoogleOauthClientId != "" {
		updates.GoogleOauthClientID = &req.GoogleOauthClientId
	}
	if req.GoogleOauthClientSecret != "" {
		updates.GoogleOauthClientSecret = &req.GoogleOauthClientSecret
	}

	// Conditionally configure Microsoft OAuth
	if req.MicrosoftOauthClientId != "" {
		updates.MicrosoftOauthClientID = &req.MicrosoftOauthClientId
	}
	if req.MicrosoftOauthClientSecret != "" {
		updates.MicrosoftOauthClientSecret = &req.MicrosoftOauthClientSecret
	}

	// Conditionally enable/disable login methods
	if req.LogInWithGoogleEnabled != project.LogInWithGoogleEnabled {
		updates.LogInWithGoogleEnabled = req.LogInWithGoogleEnabled
	}
	if req.LogInWithMicrosoftEnabled != project.LogInWithMicrosoftEnabled {
		updates.LogInWithMicrosoftEnabled = req.LogInWithMicrosoftEnabled
	}
	if req.LogInWithPasswordEnabled != project.LogInWithPasswordEnabled {
		updates.LogInWithPasswordEnabled = req.LogInWithPasswordEnabled
	}

	updatedProject, err := q.UpdateProject(ctx, updates)
	if err != nil {
		return nil, err
	}

	if err := commit(); err != nil {
		return nil, err
	}

	return transformProject(&updatedProject), nil
}

func transformProject(project *queries.Project) *openauthv1.Project {
	return &openauthv1.Project{
		Id: project.ID.String(),
		OrganizationId: project.OrganizationID.String(),
		LogInWithPasswordEnabled: project.LogInWithPasswordEnabled,
		LogInWithGoogleEnabled: project.LogInWithGoogleEnabled,
		LogInWithMicrosoftEnabled: project.LogInWithMicrosoftEnabled,
		GoogleOauthClientId: *project.GoogleOauthClientID,
		GoogleOauthClientSecret: *project.GoogleOauthClientSecret,
		MicrosoftOauthClientId: *project.MicrosoftOauthClientID,
		MicrosoftOauthClientSecret: *project.MicrosoftOauthClientSecret,
	}
}