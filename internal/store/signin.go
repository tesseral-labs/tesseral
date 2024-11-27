package store

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	intermediatev1 "github.com/openauth-dev/openauth/internal/gen/intermediate/v1"
	"github.com/openauth-dev/openauth/internal/store/idformat"
	"github.com/openauth-dev/openauth/internal/store/queries"
	"github.com/openauth-dev/openauth/internal/ujwt"
)

func (s *Store) SignInWithEmail(
	ctx *context.Context,
	req *intermediatev1.SignInWithEmailRequest,
) (*intermediatev1.SignInWithEmailResponse, error) {
	_, q, commit, rollback, err := s.tx(*ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	projectId, err := idformat.Project.Parse(req.ProjectId)
	if err != nil {
		return nil, err
	}

	users, err := q.ListUsersByEmail(*ctx, &req.Email)
	if err != nil {
		return nil, err
	}

	if users != nil {
		// TODO: Implement factor checking before issuing a session
		return nil, errors.New("not implemented")
	} else {
		// Send a verification email then issue an intermediate session, 
		// so the user can verify their email address and create an organization

		expiresAt := time.Now().Add(15 * time.Minute)

		signingKey, err := s.GetIntermediateSessionSigningKeyByProjectID(*ctx, req.ProjectId)
		if err != nil {
			return nil, err
		}

		signingKeyId := idformat.IntermediateSessionSigningKey.Format(signingKey.ID)

		sessionToken :=  ujwt.Sign(string(signingKeyId), signingKey.PrivateKey, &intermediatev1.IntermediateSessionClaims{
			Email: req.Email,
			ExpiresAt: expiresAt.Unix(),
			IssuedAt: time.Now().Unix(),
			ProjectId: req.ProjectId,
		})

		intermediateSession, err := q.CreateIntermediateSession(*ctx, queries.CreateIntermediateSessionParams{
			ID: uuid.New(),
			ProjectID: projectId,
			UnverifiedEmail: &req.Email,
			ExpireTime: &expiresAt,
			Token: sessionToken,
		})
		if err != nil {
			return nil, err
		}

		if err := commit(); err != nil {
			return nil, err
		}

		return &intermediatev1.SignInWithEmailResponse{
			SessionToken: intermediateSession.Token,
		}, nil
	}
}