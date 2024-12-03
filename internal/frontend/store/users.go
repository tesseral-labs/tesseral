package store

import (
	frontendv1 "github.com/openauth/openauth/internal/frontend/gen/openauth/frontend/v1"
	"github.com/openauth/openauth/internal/frontend/store/queries"
	"github.com/openauth/openauth/internal/store/idformat"
)

func parseUser(qUser *queries.User) *frontendv1.User {
	return &frontendv1.User{
		Id:             idformat.User.Format(qUser.ID),
		OrganizationId: idformat.Organization.Format(qUser.OrganizationID),
	}
}
