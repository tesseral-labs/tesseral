package store

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/openauth/openauth/internal/emailaddr"
	"github.com/openauth/openauth/internal/scim/authn"
	"github.com/openauth/openauth/internal/scim/internal/scimpatch"
	"github.com/openauth/openauth/internal/scim/store/queries"
	"github.com/openauth/openauth/internal/store/idformat"
)

type ListUsersRequest struct {
	Count      int
	StartIndex int
	UserName   string
}

type ListUsersResponse struct {
	TotalResults int    `json:"totalResults"`
	Users        []User `json:"Resources"`
}

// User is a SCIM representation of a user. It is suitable for JSON
// serialization.
type User any

// parsedUser is our preferred representation of SCIM users.
//
// Most IDPs will sometimes use different representations of users. Entra, in
// particular, sends "active" as a string instead of a boolean.
type parsedUser struct {
	ID       string `json:"id"`
	UserName string `json:"userName"`
	Active   bool   `json:"active"`
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
					TotalResults: 0,
					Users:        []User{},
				}, nil
			}
			return nil, fmt.Errorf("get user by email: %w", err)
		}

		return &ListUsersResponse{
			TotalResults: 1,
			Users:        []User{formatUser(qUser)},
		}, nil
	}

	count, err := q.CountUsers(ctx, authn.OrganizationID(ctx))
	if err != nil {
		return nil, fmt.Errorf("count users: %w", err)
	}

	limit := int32(10)
	if req.Count != 0 {
		limit = int32(req.Count)
	}

	offset := int32(0)
	if req.StartIndex != 0 {
		offset = int32(req.StartIndex-1) * limit
	}

	qUsers, err := q.ListUsers(ctx, queries.ListUsersParams{
		OrganizationID: authn.OrganizationID(ctx),
		Limit:          limit,
		Offset:         offset,
	})
	if err != nil {
		return nil, fmt.Errorf("list users: %w", err)
	}

	users := []User{} // intentionally not initialized as nil to avoid a JSON `null`
	for _, qUser := range qUsers {
		users = append(users, formatUser(qUser))
	}

	return &ListUsersResponse{
		TotalResults: int(count),
		Users:        users,
	}, nil
}

func (s *Store) GetUser(ctx context.Context, id string) (User, error) {
	_, q, _, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	userID, err := idformat.User.Parse(id)
	if err != nil {
		return nil, &SCIMError{
			Status: http.StatusBadRequest,
			Detail: fmt.Sprintf("invalid user id: %v", err),
		}
	}

	qUser, err := q.GetUserByID(ctx, queries.GetUserByIDParams{
		OrganizationID: authn.OrganizationID(ctx),
		ID:             userID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, &SCIMError{
				Status: http.StatusNotFound,
				Detail: "user not found",
			}
		}

		return nil, fmt.Errorf("get user by id: %w", err)
	}

	return formatUser(qUser), nil
}

func (s *Store) CreateUser(ctx context.Context, user User) (User, error) {
	_, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	parsed, err := parseUser(user)
	if err != nil {
		return nil, fmt.Errorf("parse user: %w", err)
	}

	//if err := s.validateEmailDomain(ctx, q, parsed.UserName); err != nil {
	//	return nil, fmt.Errorf("validate email domain: %w", err)
	//}

	qUser, err := q.CreateUser(ctx, queries.CreateUserParams{
		ID:             uuid.New(),
		OrganizationID: authn.OrganizationID(ctx),
		Email:          parsed.UserName,
	})
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}

	if err := commit(); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}

	return formatUser(qUser), nil
}

func (s *Store) UpdateUser(ctx context.Context, id string, user User) (User, error) {
	_, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	userID, err := idformat.User.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("parse user id: %w", err)
	}

	parsed, err := parseUser(user)
	if err != nil {
		return nil, fmt.Errorf("parse user: %w", err)
	}

	//if err := s.validateEmailDomain(ctx, q, parsed.UserName); err != nil {
	//	return nil, fmt.Errorf("validate email domain: %w", err)
	//}

	// todo do we care about this bumping deactivate_time any time an update happens to the user?
	var deactivateTime *time.Time
	if !parsed.Active {
		now := time.Now()
		deactivateTime = &now
	}

	qUser, err := q.UpdateUser(ctx, queries.UpdateUserParams{
		OrganizationID: authn.OrganizationID(ctx),
		ID:             userID,
		DeactivateTime: deactivateTime,
		Email:          parsed.UserName,
	})
	if err != nil {
		var pgxErr *pgconn.PgError
		if errors.As(err, &pgxErr) {
			if pgxErr.Code == "23505" && pgxErr.ConstraintName == "users_organization_id_email_key" {
				return nil, &SCIMError{
					Status: http.StatusBadRequest,
					Detail: "a user with that email already exists",
				}
			}
		}

		return nil, fmt.Errorf("update user: %w", err)
	}

	if err := commit(); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}

	return formatUser(qUser), nil
}

type PatchOperations struct {
	Operations []scimpatch.Operation `json:"Operations"`
}

func (s *Store) PatchUser(ctx context.Context, id string, operations PatchOperations) (User, error) {
	_, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	// load our representation of the current state
	userID, err := idformat.User.Parse(id)
	if err != nil {
		return nil, &SCIMError{
			Status: http.StatusBadRequest,
			Detail: "invalid user id",
		}
	}

	qUser, err := q.GetUserByID(ctx, queries.GetUserByIDParams{
		OrganizationID: authn.OrganizationID(ctx),
		ID:             userID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, &SCIMError{
				Status: http.StatusNotFound,
				Detail: "user not found",
			}
		}

		return nil, fmt.Errorf("get user by id: %w", err)
	}

	// load current state in SCIM representation
	scimUser := jsonify(formatUser(qUser))

	// apply patches to that representation
	if err := scimpatch.Patch(operations.Operations, &scimUser); err != nil {
		return nil, fmt.Errorf("patch user: %w", err)
	}

	// convert back to preferred representation
	parsed, err := parseUser(scimUser)
	if err != nil {
		fmt.Printf("fail to parse %v %v\n", operations, scimUser)

		return nil, fmt.Errorf("parse patched user: %w", err)
	}

	// IDPs may deprovision a user by PATCHing away everything except "active".
	// We always require email. So if a PATCHed scimUser lacks a userName (i.e.
	// email), restore it.
	if parsed.UserName == "" {
		parsed.UserName = qUser.Email
	}

	// save that new state
	//if err := s.validateEmailDomain(ctx, q, parsed.UserName); err != nil {
	//	return nil, fmt.Errorf("validate email domain: %w", err)
	//}

	// todo do we care about this bumping deactivate_time any time an update happens to the user?
	var deactivateTime *time.Time
	if !parsed.Active {
		now := time.Now()
		deactivateTime = &now
	}

	qUser, err = q.UpdateUser(ctx, queries.UpdateUserParams{
		OrganizationID: authn.OrganizationID(ctx),
		ID:             userID,
		DeactivateTime: deactivateTime,
		Email:          parsed.UserName,
	})
	if err != nil {
		var pgxErr *pgconn.PgError
		if errors.As(err, &pgxErr) {
			if pgxErr.Code == "23505" && pgxErr.ConstraintName == "users_organization_id_email_key" {
				return nil, &SCIMError{
					Status: http.StatusBadRequest,
					Detail: "a user with that email already exists",
				}
			}
		}

		return nil, fmt.Errorf("update user: %w", err)
	}

	if err := commit(); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}

	return formatUser(qUser), nil
}

func (s *Store) DeleteUser(ctx context.Context, id string) error {
	_, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return err
	}
	defer rollback()

	userID, err := idformat.User.Parse(id)
	if err != nil {
		return &SCIMError{
			Status: http.StatusBadRequest,
			Detail: "invalid user id",
		}
	}

	now := time.Now()
	if _, err := q.DeactivateUser(ctx, queries.DeactivateUserParams{
		OrganizationID: authn.OrganizationID(ctx),
		ID:             userID,
		DeactivateTime: &now,
	}); err != nil {
		return fmt.Errorf("deactivate user: %w", err)
	}

	if err := commit(); err != nil {
		return fmt.Errorf("commit: %w", err)
	}

	return nil
}

func parseUser(user User) (*parsedUser, error) {
	m, ok := user.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("users must be objects")
	}

	userName, _ := m["userName"].(string)

	var active bool
	if _, ok := m["active"]; !ok {
		active = true
	} else if a, ok := m["active"].(bool); ok {
		active = a
	} else if a, ok := m["active"].(string); ok {
		if a == "true" {
			active = true
		} else if a == "false" {
			active = false
		} else {
			return nil, fmt.Errorf("active must be a boolean")
		}
	} else {
		return nil, fmt.Errorf("active must be a boolean")
	}

	return &parsedUser{
		UserName: userName,
		Active:   active,
	}, nil
}

func formatUser(qUser queries.User) User {
	return parsedUser{
		ID:       idformat.User.Format(qUser.ID),
		UserName: qUser.Email,
		Active:   qUser.DeactivateTime == nil,
	}
}

// SCIMError is a JSON-serializable SCIM error.
type SCIMError struct {
	Status int    `json:"status"`
	Detail string `json:"detail"`
}

func (e *SCIMError) Error() string {
	return fmt.Sprintf("scim error: %d: %s", e.Status, e.Detail)
}

func (s *Store) validateEmailDomain(ctx context.Context, q *queries.Queries, email string) error {
	domain, err := emailaddr.Parse(email)
	if err != nil {
		return &SCIMError{
			Status: http.StatusBadRequest,
			Detail: "userName must be an email address",
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
		return &SCIMError{
			Status: http.StatusBadRequest,
			Detail: fmt.Sprintf("userName is not from the list of allowed domains: %s", strings.Join(qOrganizationDomains, ", ")),
		}
	}

	return nil
}

func jsonify(t any) map[string]any {
	b, err := json.Marshal(t)
	if err != nil {
		panic(err)
	}

	var v map[string]any
	if err := json.Unmarshal(b, &v); err != nil {
		panic(err)
	}
	return v
}
