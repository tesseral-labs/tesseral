package store

import (
	"context"

	intermediatev1 "github.com/openauth/openauth/internal/intermediate/gen/openauth/intermediate/v1"
)

func (s *Store) ExchangeIntermediateSessionForSession(ctx context.Context, req *intermediatev1.ExchangeIntermediateSessionForSessionRequest) (*intermediatev1.ExchangeIntermediateSessionForSessionResponse, error) {
	// _, q, commit, rollback, err := s.tx(ctx)
	// if err != nil {
	// 	return nil, err
	// }
	// defer rollback()

	// projectID := projectid.ProjectID(ctx)

	// organizationID, err := idformat.Organization.Parse(req.OrganizationId)
	// if err != nil {
	// 	return nil, err
	// }

	// organization, err := q.GetProjectOrganizationByID(ctx, queries.GetProjectOrganizationByIDParams{
	// 	ID:        organizationID,
	// 	ProjectID: projectID,
	// })
	// if err != nil {
	// 	return nil, err
	// }

	// user, err := q.GetOrganizationUser

	return nil, nil
}
