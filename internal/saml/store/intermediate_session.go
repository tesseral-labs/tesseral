package store

import (
	"context"
	"crypto/sha256"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/tesseral-labs/tesseral/internal/common/apierror"
	intermediatev1 "github.com/tesseral-labs/tesseral/internal/intermediate/gen/tesseral/intermediate/v1"
	"github.com/tesseral-labs/tesseral/internal/saml/authn"
	"github.com/tesseral-labs/tesseral/internal/saml/store/queries"
	"github.com/tesseral-labs/tesseral/internal/store/idformat"
)

type IntermediateSession struct {
	IntermediateSession *intermediatev1.IntermediateSession
	SecretToken         string
}

func (s *Store) CreateIntermediateSession(ctx context.Context) (*IntermediateSession, error) {
	_, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	const intermediateSessionDuration = time.Minute * 15
	expireTime := time.Now().Add(intermediateSessionDuration)

	secretToken := uuid.New()
	secretTokenSHA256 := sha256.Sum256(secretToken[:])
	qIntermediateSession, err := q.CreateIntermediateSession(ctx, queries.CreateIntermediateSessionParams{
		ID:                uuid.Must(uuid.NewV7()),
		ProjectID:         authn.ProjectID(ctx),
		ExpireTime:        &expireTime,
		SecretTokenSha256: secretTokenSHA256[:],
	})
	if err != nil {
		return nil, fmt.Errorf("create intermediate session: %w", err)
	}

	if err := commit(); err != nil {
		return nil, fmt.Errorf("commit transaction: %w", err)
	}

	return &IntermediateSession{
		IntermediateSession: parseIntermediateSession(qIntermediateSession),
		SecretToken:         idformat.IntermediateSessionSecretToken.Format(secretToken),
	}, nil
}

func (s *Store) AuthenticateIntermediateSession(ctx context.Context, projectUUID uuid.UUID, secretToken string) (*intermediatev1.IntermediateSession, error) {
	_, q, _, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	secretTokenUUID, err := idformat.IntermediateSessionSecretToken.Parse(secretToken)
	if err != nil {
		return nil, apierror.NewInvalidArgumentError("invalid intermediate session secret token", fmt.Errorf("parse intermediate session secret token: %w", err))
	}

	secretTokenSHA := sha256.Sum256(secretTokenUUID[:])
	qIntermediateSession, err := q.GetIntermediateSessionByTokenSHA256AndProjectID(ctx, queries.GetIntermediateSessionByTokenSHA256AndProjectIDParams{
		ProjectID:         projectUUID,
		SecretTokenSha256: secretTokenSHA[:],
	})
	if err != nil {
		return nil, fmt.Errorf("get intermediate session by token sha256 and project id: %w", err)
	}

	return parseIntermediateSession(qIntermediateSession), nil
}

func parseIntermediateSession(qIntermediateSession queries.IntermediateSession) *intermediatev1.IntermediateSession {
	return &intermediatev1.IntermediateSession{
		Id: idformat.IntermediateSession.Format(qIntermediateSession.ID),
	}
}
