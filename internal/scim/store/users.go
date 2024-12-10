package store

import (
	"context"
	"fmt"

	"github.com/openauth/openauth/internal/scim/authn"
	"github.com/openauth/openauth/internal/scim/store/queries"
	"github.com/openauth/openauth/internal/store/idformat"
)

type ListUsersRequest struct {
}

type ListUsersResponse struct {
	Schemas      []string `json:"schemas,omitempty"`
	TotalResults int      `json:"totalResults"`
	Users        []User   `json:"Resources"`
}

type User struct {
	Schemas  []string `json:"schemas,omitempty"`
	ID       string   `json:"id"`
	UserName string   `json:"username"`
}

func (s *Store) ListUsers(ctx context.Context, req *ListUsersRequest) (*ListUsersResponse, error) {
	_, q, _, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	count, err := q.CountUsers(ctx, authn.OrganizationID(ctx))
	if err != nil {
		return nil, fmt.Errorf("count users: %w", err)
	}

	qUsers, err := q.ListUsers(ctx, queries.ListUsersParams{
		OrganizationID: authn.OrganizationID(ctx),
		Limit:          10,
		Offset:         0,
	})
	if err != nil {
		return nil, fmt.Errorf("list users: %w", err)
	}

	var users []User
	for _, qUser := range qUsers {
		users = append(users, User{
			ID:       idformat.User.Format(qUser.ID),
			UserName: qUser.Email,
		})
	}

	return &ListUsersResponse{
		Schemas:      []string{"urn:ietf:params:scim:schemas:core:2.0:User"},
		TotalResults: int(count),
		Users:        users,
	}, nil
}
