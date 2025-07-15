package store

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"slices"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	auditlogv1 "github.com/tesseral-labs/tesseral/internal/auditlog/gen/tesseral/auditlog/v1"
	"github.com/tesseral-labs/tesseral/internal/emailaddr"
	"github.com/tesseral-labs/tesseral/internal/scim/authn"
	"github.com/tesseral-labs/tesseral/internal/scim/internal/scimpatch"
	"github.com/tesseral-labs/tesseral/internal/scim/store/queries"
	"github.com/tesseral-labs/tesseral/internal/store/idformat"
)

type ListUsersRequest struct {
	Count      int
	StartIndex int
	UserName   string
}

type ListUsersResponse struct {
	Schemas      []string `json:"schemas,omitempty"`
	TotalResults int      `json:"totalResults"`
	Users        []User   `json:"Resources"`
}

// User is a SCIM representation of a user. It is suitable for JSON
// serialization.
type User any

// parsedUser is our preferred representation of SCIM users.
//
// Most IDPs will sometimes use different representations of users. Entra, in
// particular, sends "active" as a string instead of a boolean.
type parsedUser struct {
	Schemas  []string `json:"schemas,omitempty"`
	ID       string   `json:"id"`
	UserName string   `json:"userName"`
	Active   bool     `json:"active"`
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
					Schemas:      []string{"urn:ietf:params:scim:api:messages:2.0:ListResponse"},
					TotalResults: 0,
					Users:        []User{},
				}, nil
			}
			return nil, fmt.Errorf("get user by email: %w", err)
		}

		return &ListUsersResponse{
			Schemas:      []string{"urn:ietf:params:scim:api:messages:2.0:ListResponse"},
			TotalResults: 1,
			Users:        []User{formatUser(false, qUser, true)},
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
		users = append(users, formatUser(false, qUser, true))
	}

	return &ListUsersResponse{
		Schemas:      []string{"urn:ietf:params:scim:api:messages:2.0:ListResponse"},
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

	return formatUser(true, qUser, true), nil
}

func (s *Store) CreateUser(ctx context.Context, user User) (User, error) {
	tx, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	parsed, err := parseUser(user)
	if err != nil {
		return nil, fmt.Errorf("parse user: %w", err)
	}

	if err := s.validateEmailDomain(ctx, q, parsed.UserName); err != nil {
		return nil, fmt.Errorf("validate email domain: %w", err)
	}

	qUser, err := q.CreateUser(ctx, queries.CreateUserParams{
		ID:             uuid.New(),
		OrganizationID: authn.OrganizationID(ctx),
		Email:          parsed.UserName,
	})
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}

	auditUser, err := s.auditlogStore.GetUser(ctx, tx, qUser.ID)
	if err != nil {
		return nil, fmt.Errorf("get user for audit log: %w", err)
	}

	if _, err := s.logAuditEvent(ctx, q, logAuditEventParams{
		EventName: "tesseral.users.create",
		EventDetails: &auditlogv1.CreateUser{
			User: auditUser,
		},
		OrganizationID: &qUser.OrganizationID,
		ResourceType:   queries.AuditLogEventResourceTypeUser,
		ResourceID:     &qUser.ID,
	}); err != nil {
		return nil, fmt.Errorf("log audit event: %w", err)
	}

	if err := commit(); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}

	return formatUser(true, qUser, true), nil
}

func (s *Store) UpdateUser(ctx context.Context, id string, user User) (User, error) {
	tx, q, commit, rollback, err := s.tx(ctx)
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

	if err := s.validateEmailDomain(ctx, q, parsed.UserName); err != nil {
		return nil, fmt.Errorf("validate email domain: %w", err)
	}

	if !parsed.Active {
		return s.DeleteUser(ctx, id)
	}

	auditPreviousUser, err := s.auditlogStore.GetUser(ctx, tx, userID)
	if err != nil {
		return nil, fmt.Errorf("get user for audit log: %w", err)
	}

	qUser, err := q.UpdateUser(ctx, queries.UpdateUserParams{
		OrganizationID: authn.OrganizationID(ctx),
		ID:             userID,
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

	auditUser, err := s.auditlogStore.GetUser(ctx, tx, qUser.ID)
	if err != nil {
		return nil, fmt.Errorf("get user for audit log: %w", err)
	}

	if _, err := s.logAuditEvent(ctx, q, logAuditEventParams{
		EventName: "tesseral.users.update",
		EventDetails: &auditlogv1.UpdateUser{
			PreviousUser: auditPreviousUser,
			User:         auditUser,
		},
		OrganizationID: &qUser.OrganizationID,
		ResourceType:   queries.AuditLogEventResourceTypeUser,
		ResourceID:     &qUser.ID,
	}); err != nil {
		return nil, fmt.Errorf("log audit event: %w", err)
	}

	if err := commit(); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}

	return formatUser(true, qUser, true), nil
}

type PatchOperations struct {
	Operations []scimpatch.Operation `json:"Operations"`
}

func (s *Store) PatchUser(ctx context.Context, id string, operations PatchOperations) (User, error) {
	tx, q, commit, rollback, err := s.tx(ctx)
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
	scimUser := jsonify(formatUser(false, qUser, true))

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

	if err := s.validateEmailDomain(ctx, q, parsed.UserName); err != nil {
		return nil, fmt.Errorf("validate email domain: %w", err)
	}

	if !parsed.Active {
		return s.DeleteUser(ctx, id)
	}

	auditPreviousUser, err := s.auditlogStore.GetUser(ctx, tx, userID)
	if err != nil {
		return nil, fmt.Errorf("get user for audit log: %w", err)
	}

	qUser, err = q.UpdateUser(ctx, queries.UpdateUserParams{
		OrganizationID: authn.OrganizationID(ctx),
		ID:             userID,
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

	auditUser, err := s.auditlogStore.GetUser(ctx, tx, qUser.ID)
	if err != nil {
		return nil, fmt.Errorf("get user for audit log: %w", err)
	}

	if _, err := s.logAuditEvent(ctx, q, logAuditEventParams{
		EventName: "tesseral.users.update",
		EventDetails: &auditlogv1.UpdateUser{
			PreviousUser: auditPreviousUser,
			User:         auditUser,
		},
		OrganizationID: &qUser.OrganizationID,
		ResourceType:   queries.AuditLogEventResourceTypeUser,
		ResourceID:     &qUser.ID,
	}); err != nil {
		return nil, fmt.Errorf("log audit event: %w", err)
	}

	if err := commit(); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}

	return formatUser(true, qUser, true), nil
}

func (s *Store) DeleteUser(ctx context.Context, id string) (User, error) {
	tx, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

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

	auditUser, err := s.auditlogStore.GetUser(ctx, tx, userID)
	if err != nil {
		return nil, fmt.Errorf("get user for audit log: %w", err)
	}

	if _, err := q.DeleteUser(ctx, queries.DeleteUserParams{
		ID:             userID,
		OrganizationID: authn.OrganizationID(ctx),
	}); err != nil {
		return nil, fmt.Errorf("delete user: %w", err)
	}

	if _, err := s.logAuditEvent(ctx, q, logAuditEventParams{
		EventName: "tesseral.users.delete",
		EventDetails: &auditlogv1.DeleteUser{
			User: auditUser,
		},
		OrganizationID: &qUser.OrganizationID,
		ResourceType:   queries.AuditLogEventResourceTypeUser,
		ResourceID:     &qUser.ID,
	}); err != nil {
		return nil, fmt.Errorf("log audit event: %w", err)
	}

	if err := commit(); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}

	return formatUser(true, qUser, false), nil
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
		switch a {
		case "True":
			active = true
		case "False":
			active = false
		default:
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

func formatUser(withSchema bool, qUser queries.User, active bool) User {
	var schemas []string
	if withSchema {
		schemas = []string{"urn:ietf:params:scim:schemas:core:2.0:User"}
	}

	return parsedUser{
		Schemas:  schemas,
		ID:       idformat.User.Format(qUser.ID),
		UserName: qUser.Email,
		Active:   active,
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

	domainOk := slices.Contains(qOrganizationDomains, domain)
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
