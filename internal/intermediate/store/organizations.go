package store

import (
	"context"

	"github.com/google/uuid"
	"github.com/openauth/openauth/internal/intermediate/authn"
	intermediatev1 "github.com/openauth/openauth/internal/intermediate/gen/openauth/intermediate/v1"
	"github.com/openauth/openauth/internal/intermediate/store/queries"
	"github.com/openauth/openauth/internal/store/idformat"
)

func (s *Store) ListOrganizations(
	ctx context.Context,
	req *intermediatev1.ListOrganizationsRequest,
) (*intermediatev1.ListOrganizationsResponse, error) {
	_, q, _, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	intermediateSession := authn.IntermediateSession(ctx)

	var startID uuid.UUID
	if err := s.pageEncoder.Unmarshal(req.PageToken, &startID); err != nil {
		return nil, err
	}

	qProject, err := q.GetProjectByID(ctx, authn.ProjectID(ctx))
	if err != nil {
		return nil, err
	}

	limit := 10
	qOrganizationRecords := []queries.Organization{}

	if intermediateSession.GoogleUserId != "" {
		qGoogleOrganizationRecords, err := q.ListOrganizationsByGoogleUserID(ctx, queries.ListOrganizationsByGoogleUserIDParams{
			Email:        intermediateSession.Email,
			GoogleUserID: &intermediateSession.GoogleUserId,
			ID:           startID,
			Limit:        int32(limit + 1),
			ProjectID:    authn.ProjectID(ctx),
		})
		if err != nil {
			return nil, err
		}

		if len(qGoogleOrganizationRecords) > 0 {
			qOrganizationRecords = qGoogleOrganizationRecords
		}
	} else if intermediateSession.MicrosoftUserId != "" {
		qMicrosoftOrganizationRecords, err := q.ListOrganizationsByMicrosoftUserID(ctx, queries.ListOrganizationsByMicrosoftUserIDParams{
			Email:           intermediateSession.Email,
			ID:              startID,
			Limit:           int32(limit + 1),
			MicrosoftUserID: &intermediateSession.MicrosoftUserId,
			ProjectID:       authn.ProjectID(ctx),
		})
		if err != nil {
			return nil, err
		}

		if len(qMicrosoftOrganizationRecords) > 0 {
			qOrganizationRecords = qMicrosoftOrganizationRecords
		}
	} else {
		qEmailOrganizationRecords, err := q.ListOrganizationsByEmail(ctx, queries.ListOrganizationsByEmailParams{
			Email:     intermediateSession.Email,
			ID:        startID,
			Limit:     int32(limit + 1),
			ProjectID: authn.ProjectID(ctx),
		})
		if err != nil {
			return nil, err
		}

		if len(qEmailOrganizationRecords) > 0 {
			qOrganizationRecords = qEmailOrganizationRecords
		}
	}

	organizations := []*intermediatev1.Organization{}
	for _, organization := range qOrganizationRecords {
		organizations = append(organizations, parseOrganization(organization, qProject))
	}

	var nextPageToken string
	if len(organizations) == limit+1 {
		nextPageToken = s.pageEncoder.Marshal(organizations[limit].Id)
		organizations = organizations[:limit]
	}

	return &intermediatev1.ListOrganizationsResponse{
		Organizations: organizations,
		NextPageToken: nextPageToken,
	}, nil
}

func parseOrganization(organization queries.Organization, project queries.Project) *intermediatev1.Organization {
	logInWithGoogleEnabled := project.LogInWithGoogleEnabled && (!organization.OverrideLogInMethods || organization.OverrideLogInWithGoogleEnabled != nil && *organization.OverrideLogInWithGoogleEnabled)
	logInWithMicrosoftEnabled := project.LogInWithMicrosoftEnabled && (!organization.OverrideLogInMethods || organization.OverrideLogInWithMicrosoftEnabled != nil && *organization.OverrideLogInWithMicrosoftEnabled)
	logInWithPasswordEnabled := project.LogInWithPasswordEnabled && (!organization.OverrideLogInMethods || organization.OverrideLogInWithPasswordEnabled != nil && *organization.OverrideLogInWithPasswordEnabled)

	return &intermediatev1.Organization{
		Id:                        idformat.Organization.Format(organization.ID),
		DisplayName:               organization.DisplayName,
		LogInWithGoogleEnabled:    logInWithGoogleEnabled,
		LogInWithMicrosoftEnabled: logInWithMicrosoftEnabled,
		LogInWithPasswordEnabled:  logInWithPasswordEnabled,
	}
}
