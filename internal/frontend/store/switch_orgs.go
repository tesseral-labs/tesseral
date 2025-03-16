package store

import (
	"context"
	"fmt"

	"github.com/tesseral-labs/tesseral/internal/frontend/authn"
	frontendv1 "github.com/tesseral-labs/tesseral/internal/frontend/gen/tesseral/frontend/v1"
	"github.com/tesseral-labs/tesseral/internal/frontend/store/queries"
	"github.com/tesseral-labs/tesseral/internal/store/idformat"
)

func (s *Store) ListSwitchableOrganizations(ctx context.Context, req *frontendv1.ListSwitchableOrganizationsRequest) (*frontendv1.ListSwitchableOrganizationsResponse, error) {
	_, q, _, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	qUser, err := q.GetUserByID(ctx, authn.UserID(ctx))
	if err != nil {
		return nil, fmt.Errorf("get user by id: %w", err)
	}

	fmt.Println("list switchable organizations", qUser.Email)

	qOrgs, err := q.ListSwitchableOrganizations(ctx, queries.ListSwitchableOrganizationsParams{
		ProjectID: authn.ProjectID(ctx),
		Email:     qUser.Email,
	})
	if err != nil {
		return nil, fmt.Errorf("list switchable organizations: %w", err)
	}

	fmt.Println("list switchable organizations", qOrgs)

	var orgs []*frontendv1.SwitchableOrganization
	for _, qOrg := range qOrgs {
		displayName := qOrg.DisplayName
		if authn.ProjectID(ctx) == *s.dogfoodProjectID {
			// for the dogfood project, use the display name of the project this
			// org backs
			qProject, err := q.GetProjectByBackingOrganizationID(ctx, &qOrg.ID)
			if err != nil {
				return nil, fmt.Errorf("get project by backing organization id: %w", err)
			}

			displayName = qProject.DisplayName
		}

		orgs = append(orgs, &frontendv1.SwitchableOrganization{
			Id:          idformat.Organization.Format(qOrg.ID),
			DisplayName: displayName,
		})
	}

	return &frontendv1.ListSwitchableOrganizationsResponse{SwitchableOrganizations: orgs}, nil
}
