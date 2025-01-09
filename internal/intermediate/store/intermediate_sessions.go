package store

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/openauth/openauth/internal/intermediate/authn"
	intermediatev1 "github.com/openauth/openauth/internal/intermediate/gen/openauth/intermediate/v1"
	"github.com/openauth/openauth/internal/intermediate/store/queries"
	"github.com/openauth/openauth/internal/store/idformat"
)

func (s *Store) GetIntermediateSessionByToken(ctx context.Context, token string) (*intermediatev1.IntermediateSession, error) {
	tokenUUID, err := idformat.IntermediateSessionToken.Parse(token)
	if err != nil {
		return nil, fmt.Errorf("parse token: %w", err)
	}

	tokenSHA256 := sha256.Sum256(tokenUUID[:])
	qIntermediateSession, err := s.q.GetIntermediateSessionByTokenSHA256(ctx, tokenSHA256[:])
	if err != nil {
		return nil, fmt.Errorf("get intermediate session by token sha256: %w", err)
	}

	// todo this is what parseIntermediateSession should do, but already token
	// by another function returning a hand-written type
	intermediateSession := &intermediatev1.IntermediateSession{
		Id:        idformat.IntermediateSession.Format(qIntermediateSession.ID),
		ProjectId: idformat.Project.Format(qIntermediateSession.ProjectID),
	}

	if qIntermediateSession.Email != nil {
		intermediateSession.Email = *qIntermediateSession.Email
	}

	if qIntermediateSession.GoogleUserID != nil {
		intermediateSession.GoogleUserId = *qIntermediateSession.GoogleUserID
	}

	if qIntermediateSession.MicrosoftUserID != nil {
		intermediateSession.MicrosoftUserId = *qIntermediateSession.MicrosoftUserID
	}

	return intermediateSession, nil
}

func (s *Store) Whoami(ctx context.Context, req *intermediatev1.WhoamiRequest) (*intermediatev1.WhoamiResponse, error) {
	_, q, _, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	intermediateSession := authn.IntermediateSession(ctx)
	var isEmailVerified bool

	if intermediateSession.GoogleUserId != "" {
		// Check if the google user id is verified
		isGoogleEmailVerified, err := q.IsGoogleEmailVerified(ctx, queries.IsGoogleEmailVerifiedParams{
			Email:        intermediateSession.Email,
			GoogleUserID: &intermediateSession.GoogleUserId,
			ProjectID:    authn.ProjectID(ctx),
		})
		if err != nil {
			return nil, err
		}

		isEmailVerified = isGoogleEmailVerified
	}

	return &intermediatev1.WhoamiResponse{
		Email:           intermediateSession.Email,
		GoogleUserId:    intermediateSession.GoogleUserId,
		IsEmailVerified: isEmailVerified,
		MicrosoftUserId: intermediateSession.MicrosoftUserId,
	}, nil
}

type IntermediateSession struct {
	ID                           uuid.UUID
	CreateTime                   time.Time
	Email                        string
	EmailVerificationChallengeID uuid.UUID
	ExpireTime                   time.Time
	GoogleUserID                 string
	MicrosoftUserID              string
	ProjectID                    uuid.UUID
	TokenSha256                  []byte
	Revoked                      bool
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
		ID:         uuid.Must(uuid.NewV7()),
		ProjectID:  projectId,
		Email:      &req.Email,
		ExpireTime: &expiresAt,
	})
	if err != nil {
		return nil, err
	}

	if err := commit(); err != nil {
		return nil, err
	}

	return parseIntermediateSession(&createdIntermediateSession), nil
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

	panic("unimplemented")

	// TODO what do we do here?

	//// Check if the email in the request matches the email in the intermediate session
	//if existingIntermediateSession.UnverifiedEmail != &req.Email {
	//	return nil, ErrIntermediateSessionEmailMismatch
	//}
	//
	//
	//session, err := q.VerifyIntermediateSessionEmail(*ctx, queries.VerifyIntermediateSessionEmailParams{
	//	ID:            sessionId,
	//	VerifiedEmail: &req.Email,
	//})
	//if err != nil {
	//	return nil, err
	//}

	//return parseIntermediateSession(&session), nil
}

func parseIntermediateSession(i *queries.IntermediateSession) *IntermediateSession {
	return &IntermediateSession{
		ID:              i.ID,
		CreateTime:      *i.CreateTime,
		Email:           *i.Email,
		ExpireTime:      *i.ExpireTime,
		GoogleUserID:    *i.GoogleUserID,
		MicrosoftUserID: *i.MicrosoftUserID,
		ProjectID:       i.ProjectID,
		TokenSha256:     i.TokenSha256,
		Revoked:         i.Revoked,
	}
}
