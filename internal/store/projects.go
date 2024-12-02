package store

import (
	"context"

	"github.com/google/uuid"
	backendv1 "github.com/openauth/openauth/internal/gen/backend/v1"
	openauthv1 "github.com/openauth/openauth/internal/gen/openauth/v1"
	"github.com/openauth/openauth/internal/store/idformat"
	"github.com/openauth/openauth/internal/store/queries"
)

func (s *Store) CreateProject(ctx context.Context, req *backendv1.CreateProjectRequest) (*openauthv1.Project, error) {
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
		ID:          uuid.New(),
		ProjectID:   createdProject.ID,
		DisplayName: req.DisplayName,
	})
	if err != nil {
		return nil, err
	}

	// Update the project with the dogfooding project ID
	updatedProject, err := q.UpdateProjectOrganizationID(ctx, queries.UpdateProjectOrganizationIDParams{
		ID:             createdProject.ID,
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
	return parseProject(&updatedProject), nil
}

func (s *Store) GetProject(ctx context.Context, req *backendv1.GetProjectRequest) (*openauthv1.Project, error) {
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

	return parseProject(&project), nil
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
	projectRecords, err := q.ListProjects(ctx, int32(limit+1))
	if err != nil {
		return nil, err
	}

	projects := []*openauthv1.Project{}
	for _, project := range projectRecords {
		projects = append(projects, parseProject(&project))
	}

	var nextPageToken string
	if len(projects) == limit+1 {
		nextPageToken = s.pageEncoder.Marshal(projectRecords[limit].ID)
		projects = projects[:limit]
	}

	return &backendv1.ListProjectsResponse{
		Projects:      projects,
		NextPageToken: nextPageToken,
	}, nil
}

func (s *Store) UpdateProject(ctx context.Context, req *backendv1.UpdateProjectRequest) (*openauthv1.Project, error) {
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
	if req.Project.GoogleOauthClientId != "" {
		updates.GoogleOauthClientID = &req.Project.GoogleOauthClientId
	}
	if req.Project.GoogleOauthClientSecret != "" {
		updates.GoogleOauthClientSecret = &req.Project.GoogleOauthClientSecret
	}

	// Conditionally configure Microsoft OAuth
	if req.Project.MicrosoftOauthClientId != "" {
		updates.MicrosoftOauthClientID = &req.Project.MicrosoftOauthClientId
	}
	if req.Project.MicrosoftOauthClientSecret != "" {
		updates.MicrosoftOauthClientSecret = &req.Project.MicrosoftOauthClientSecret
	}

	// Conditionally enable/disable login methods
	if req.Project.LogInWithGoogleEnabled != project.LogInWithGoogleEnabled {
		updates.LogInWithGoogleEnabled = req.Project.LogInWithGoogleEnabled
	}
	if req.Project.LogInWithMicrosoftEnabled != project.LogInWithMicrosoftEnabled {
		updates.LogInWithMicrosoftEnabled = req.Project.LogInWithMicrosoftEnabled
	}
	if req.Project.LogInWithPasswordEnabled != project.LogInWithPasswordEnabled {
		updates.LogInWithPasswordEnabled = req.Project.LogInWithPasswordEnabled
	}

	updatedProject, err := q.UpdateProject(ctx, updates)
	if err != nil {
		return nil, err
	}

	if err := commit(); err != nil {
		return nil, err
	}

	return parseProject(&updatedProject), nil
}

func parseProject(project *queries.Project) *openauthv1.Project {
	return &openauthv1.Project{
		Id:                         idformat.Project.Format(project.ID),
		OrganizationId:             idformat.Organization.Format(*project.OrganizationID),
		LogInWithPasswordEnabled:   project.LogInWithPasswordEnabled,
		LogInWithGoogleEnabled:     project.LogInWithGoogleEnabled,
		LogInWithMicrosoftEnabled:  project.LogInWithMicrosoftEnabled,
		GoogleOauthClientId:        derefOrEmpty(project.GoogleOauthClientID),
		GoogleOauthClientSecret:    derefOrEmpty(project.GoogleOauthClientSecret),
		MicrosoftOauthClientId:     derefOrEmpty(project.MicrosoftOauthClientID),
		MicrosoftOauthClientSecret: derefOrEmpty(project.MicrosoftOauthClientSecret),
	}
}
