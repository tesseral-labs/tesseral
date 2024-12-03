package store

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/openauth/openauth/internal/store/idformat"
	"github.com/openauth/openauth/internal/store/queries"
)

type IntermediateSession struct {
	ID              uuid.UUID
	ProjectID       uuid.UUID
	UnverifiedEmail string
	VerifiedEmail   string
	CreateTime      time.Time
	ExpireTime      time.Time
	Token           string
	TokenSha256     []byte
	Revoked         bool
}

var ErrIntermediateSessionRevoked = errors.New("intermediate session has been revoked")
var ErrIntermediateSessionExpired = errors.New("intermediate session has expired")
var ErrIntermediateSessionEmailMismatch = errors.New("intermediate session email mismatch")

type CreateIntermediateSessionRequest struct {
	ProjectID string
	Email     string
}

func (s *Store) CreateIntermediateSession(ctx *context.Context, req *CreateIntermediateSessionRequest) (*IntermediateSession, error) {
	_, q, commit, rollback, err := s.tx(*ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	projectId, err := idformat.Project.Parse(req.ProjectID)
	if err != nil {
		return nil, err
	}

	// Allow 15 minutes for the user to verify their email before expiring the intermediate session
	expiresAt := time.Now().Add(time.Minute * 15)

	createdIntermediateSession, err := q.CreateIntermediateSession(*ctx, queries.CreateIntermediateSessionParams{
		ID:              uuid.New(),
		ProjectID:       projectId,
		UnverifiedEmail: &req.Email,
		ExpireTime:      &expiresAt,
	})
	if err != nil {
		return nil, err
	}

	if err := commit(); err != nil {
		return nil, err
	}

	return parseIntermediateSession(&createdIntermediateSession), nil
}

func (s *Store) RevokeIntermediateSession(ctx *context.Context, ID string) error {
	_, q, commit, rollback, err := s.tx(*ctx)
	if err != nil {
		return err
	}
	defer rollback()

	sessionId, err := idformat.IntermediateSession.Parse(ID)
	if err != nil {
		return err
	}

	if _, err := q.RevokeIntermediateSession(*ctx, sessionId); err != nil {
		return err
	}

	if err := commit(); err != nil {
		return err
	}

	return nil
}

func (s *Store) GetIntermediateSession(ctx *context.Context, id string) (*IntermediateSession, error) {
	_, q, _, rollback, err := s.tx(*ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	sessionId, err := idformat.IntermediateSession.Parse(id)
	if err != nil {
		return nil, err
	}

	session, err := q.GetIntermediateSessionByID(*ctx, sessionId)
	if err != nil {
		return nil, err
	}

	return parseIntermediateSession(&session), nil
}

type VerifyIntermediateSessionEmailRequest struct {
	ID    string
	Email string
}

func (s *Store) VerifyIntermediateSessionEmail(
	ctx *context.Context,
	req *VerifyIntermediateSessionEmailRequest,
) (*IntermediateSession, error) {
	_, q, _, rollback, err := s.tx(*ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	sessionId, err := idformat.IntermediateSession.Parse(req.ID)
	if err != nil {
		return nil, err
	}

	// Get the intermediate session so we can perform some checks
	existingIntermediateSession, err := q.GetIntermediateSessionByID(*ctx, sessionId)
	if err != nil {
		return nil, err
	}

	// Check if the intermediate session has been revoked
	if existingIntermediateSession.Revoked {
		return nil, ErrIntermediateSessionRevoked
	}

	// Check if the intermediate session has expired
	if existingIntermediateSession.ExpireTime.Before(time.Now()) {
		return nil, ErrIntermediateSessionExpired
	}

	// Check if the email in the request matches the email in the intermediate session
	if existingIntermediateSession.UnverifiedEmail != &req.Email {
		return nil, ErrIntermediateSessionEmailMismatch
	}

	session, err := q.VerifyIntermediateSessionEmail(*ctx, queries.VerifyIntermediateSessionEmailParams{
		ID:            sessionId,
		VerifiedEmail: &req.Email,
	})
	if err != nil {
		return nil, err
	}

	return parseIntermediateSession(&session), nil
}

func parseIntermediateSession(i *queries.IntermediateSession) *IntermediateSession {
	return &IntermediateSession{
		ID:              i.ID,
		ProjectID:       i.ProjectID,
		UnverifiedEmail: *i.UnverifiedEmail,
		VerifiedEmail:   *i.VerifiedEmail,
		CreateTime:      *i.CreateTime,
		ExpireTime:      *i.ExpireTime,
		Token:           i.Token,
		TokenSha256:     i.TokenSha256,
		Revoked:         i.Revoked,
	}
}
