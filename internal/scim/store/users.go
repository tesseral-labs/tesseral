package store

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/openauth/openauth/internal/emailaddr"
	"github.com/openauth/openauth/internal/scim/authn"
	"github.com/openauth/openauth/internal/scim/store/queries"
	"github.com/openauth/openauth/internal/store/idformat"
)

type ListUsersRequest struct {
	UserName string
}

type ListUsersResponse struct {
	Schemas      []string `json:"schemas,omitempty"`
	TotalResults int      `json:"totalResults"`
	Users        []*User  `json:"Resources"`
}

type User struct {
	Schemas  []string `json:"schemas,omitempty"`
	ID       string   `json:"id"`
	Active   bool     `json:"active"`
	UserName string   `json:"userName"`
}

func (s *Store) ListUsers(ctx context.Context, req *ListUsersRequest) (*ListUsersResponse, error) {
	_, q, _, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	if req.UserName != "" {
		qUser, err := q.GetUserByEmail(ctx, queries.GetUserByEmailParams{
			OrganizationID: authn.OrganizationID(ctx),
			Email:          req.UserName,
		})
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return &ListUsersResponse{
					Schemas:      []string{"urn:ietf:params:scim:schemas:core:2.0:User"},
					TotalResults: 0,
					Users:        []*User{},
				}, nil
			}
			return nil, fmt.Errorf("get user by email: %w", err)
		}

		return &ListUsersResponse{
			Schemas:      []string{"urn:ietf:params:scim:schemas:core:2.0:User"},
			TotalResults: 1,
			Users:        []*User{parseUser(false, qUser)},
		}, nil
	}

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

	var users []*User
	for _, qUser := range qUsers {
		users = append(users, parseUser(false, qUser))
	}

	return &ListUsersResponse{
		Schemas:      []string{"urn:ietf:params:scim:schemas:core:2.0:User"},
		TotalResults: int(count),
		Users:        users,
	}, nil
}

func (s *Store) GetUser(ctx context.Context, id string) (*User, error) {
	_, q, _, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	userID, err := idformat.User.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("parse user id: %w", err)
	}

	qUser, err := q.GetUserByID(ctx, queries.GetUserByIDParams{
		OrganizationID: authn.OrganizationID(ctx),
		ID:             userID,
	})
	if err != nil {
		return nil, fmt.Errorf("get user by id: %w", err)
	}

	return parseUser(true, qUser), nil
}

func (s *Store) CreateUser(ctx context.Context, req *User) (*User, error) {
	_, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	if err := s.validateEmailDomain(ctx, q, req.UserName); err != nil {
		return nil, fmt.Errorf("validate email domain: %w", err)
	}

	qUser, err := q.CreateUser(ctx, queries.CreateUserParams{
		ID:             uuid.New(),
		OrganizationID: authn.OrganizationID(ctx),
		Email:          req.UserName,
	})
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}

	if err := commit(); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}

	return parseUser(true, qUser), nil
}

func (s *Store) UpdateUser(ctx context.Context, req *User) (*User, error) {
	_, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	userID, err := idformat.User.Parse(req.ID)
	if err != nil {
		return nil, fmt.Errorf("parse user id: %w", err)
	}

	if err := s.validateEmailDomain(ctx, q, req.UserName); err != nil {
		return nil, fmt.Errorf("validate email domain: %w", err)
	}

	// todo do we care about this bumping deactivate_time any time an update happens to the user?
	var deactivateTime *time.Time
	if !req.Active {
		now := time.Now()
		deactivateTime = &now
	}

	qUser, err := q.UpdateUser(ctx, queries.UpdateUserParams{
		OrganizationID: authn.OrganizationID(ctx),
		ID:             userID,
		DeactivateTime: deactivateTime,
		Email:          req.UserName,
	})
	if err != nil {
		return nil, fmt.Errorf("update user: %w", err)
	}

	if err := commit(); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}

	return parseUser(true, qUser), nil
}

func parseUser(withSchema bool, qUser queries.User) *User {
	var schemas []string
	if withSchema {
		schemas = []string{"urn:ietf:params:scim:schemas:core:2.0:User"}
	}

	return &User{
		Schemas:  schemas,
		ID:       idformat.User.Format(qUser.ID),
		Active:   qUser.DeactivateTime == nil,
		UserName: qUser.Email,
	}
}

// BadEmailDomainError indicates that a email address is not within an
// organization's list of domains.
//
// Instances of BadEmailDomainError are JSON-serializable SCIM errors.
type BadEmailDomainError struct {
	Status int    `json:"status"`
	Detail string `json:"detail"`
}

func (e *BadEmailDomainError) Error() string {
	return e.Detail
}

func (s *Store) validateEmailDomain(ctx context.Context, q *queries.Queries, email string) error {
	domain, err := emailaddr.Parse(email)
	if err != nil {
		return &BadEmailDomainError{
			Status: http.StatusBadRequest,
			Detail: fmt.Sprintf("userName must be an email address"),
		}
	}

	qOrganizationDomains, err := q.GetOrganizationDomains(ctx, authn.OrganizationID(ctx))
	if err != nil {
		return fmt.Errorf("get organization domains: %w", err)
	}

	var domainOk bool
	for _, orgDomain := range qOrganizationDomains {
		if orgDomain == domain {
			domainOk = true
			break
		}
	}

	if !domainOk {
		return &BadEmailDomainError{
			Status: http.StatusBadRequest,
			Detail: fmt.Sprintf("userName is not from the list of allowed domains: %s", strings.Join(qOrganizationDomains, ", ")),
		}
	}

	return nil
}
