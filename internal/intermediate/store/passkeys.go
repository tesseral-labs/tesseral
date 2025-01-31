package store

import (
	"context"
	"crypto/x509"
	"fmt"

	"github.com/openauth/openauth/internal/common/apierror"
	"github.com/openauth/openauth/internal/intermediate/authn"
	intermediatev1 "github.com/openauth/openauth/internal/intermediate/gen/openauth/intermediate/v1"
	"github.com/openauth/openauth/internal/intermediate/store/queries"
	"github.com/openauth/openauth/internal/webauthn"
)

func (s *Store) GetPasskeyOptions(ctx context.Context, req *intermediatev1.GetPasskeyOptionsRequest) (*intermediatev1.GetPasskeyOptionsResponse, error) {
	_, q, _, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	if err := s.checkShouldRegisterPasskey(ctx, q); err != nil {
		return nil, fmt.Errorf("check should register passkey: %w", err)
	}

	qProject, err := q.GetProjectByID(ctx, authn.ProjectID(ctx))
	if err != nil {
		return nil, fmt.Errorf("get project by id: %w", err)
	}

	qIntermediateSession, err := q.GetIntermediateSessionByID(ctx, authn.IntermediateSessionID(ctx))
	if err != nil {
		return nil, fmt.Errorf("get intermediate session by id: %w", err)
	}

	return &intermediatev1.GetPasskeyOptionsResponse{
		RpId:            *qProject.AuthDomain,
		RpName:          qProject.DisplayName,
		UserId:          fmt.Sprintf("%s|%s", *qIntermediateSession.Email, *qIntermediateSession.OrganizationID),
		UserDisplayName: *qIntermediateSession.Email,
	}, nil
}

func (s *Store) RegisterPasskey(ctx context.Context, req *intermediatev1.RegisterPasskeyRequest) (*intermediatev1.RegisterPasskeyResponse, error) {
	_, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	if err := s.checkShouldRegisterPasskey(ctx, q); err != nil {
		return nil, fmt.Errorf("check should register passkey: %w", err)
	}

	qProject, err := q.GetProjectByID(ctx, authn.ProjectID(ctx))
	if err != nil {
		return nil, fmt.Errorf("get project by id: %w", err)
	}

	cred, err := webauthn.Parse(&webauthn.ParseRequest{
		RPID:              *qProject.AuthDomain,
		AttestationObject: req.AttestationObject,
	})
	if err != nil {
		return nil, fmt.Errorf("parse webauthn credential: %w", err)
	}

	publicKey, err := x509.MarshalPKIXPublicKey(cred.PublicKey)
	if err != nil {
		return nil, fmt.Errorf("marshal public key: %w", err)
	}

	if _, err := q.RegisterPasskey(ctx, queries.RegisterPasskeyParams{
		ID:                  authn.IntermediateSessionID(ctx),
		PasskeyCredentialID: cred.ID,
		PasskeyPublicKey:    publicKey,
		PasskeyAaguid:       &cred.AAGUID,
	}); err != nil {
		return nil, fmt.Errorf("register passkey: %w", err)
	}

	if err := commit(); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}

	return &intermediatev1.RegisterPasskeyResponse{}, nil
}

func (s *Store) checkShouldRegisterPasskey(ctx context.Context, q *queries.Queries) error {
	// don't register passkeys if you're already matching a user, and that user
	// has at least one passkey

	qIntermediateSession, err := q.GetIntermediateSessionByID(ctx, authn.IntermediateSessionID(ctx))
	if err != nil {
		return fmt.Errorf("get intermediate session by id: %w", err)
	}

	qOrg, err := q.GetProjectOrganizationByID(ctx, queries.GetProjectOrganizationByIDParams{
		ProjectID: authn.ProjectID(ctx),
		ID:        *qIntermediateSession.OrganizationID,
	})
	if err != nil {
		return fmt.Errorf("get organization by id: %w", err)
	}

	qUser, err := s.matchUser(ctx, q, qOrg, qIntermediateSession)
	if err != nil {
		return fmt.Errorf("match user: %w", err)
	}

	// no matching user; it's ok to register passkeys
	if qUser == nil {
		return nil
	}

	// does the matching user have any passkeys?
	hasPasskeys, err := q.GetUserHasPasskey(ctx, qUser.ID)
	if err != nil {
		return fmt.Errorf("get user has passkey: %w", err)
	}

	if hasPasskeys {
		return apierror.NewFailedPreconditionError("user already has passkeys", nil)
	}
	return nil
}
