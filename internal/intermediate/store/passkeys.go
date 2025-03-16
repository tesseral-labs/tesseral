package store

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"fmt"

	"github.com/tesseral-labs/tesseral/internal/common/apierror"
	"github.com/tesseral-labs/tesseral/internal/intermediate/authn"
	intermediatev1 "github.com/tesseral-labs/tesseral/internal/intermediate/gen/tesseral/intermediate/v1"
	"github.com/tesseral-labs/tesseral/internal/intermediate/store/queries"
	"github.com/tesseral-labs/tesseral/internal/store/idformat"
	"github.com/tesseral-labs/tesseral/internal/webauthn"
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
		RpId:            qProject.VaultDomain,
		RpName:          qProject.DisplayName,
		UserId:          fmt.Sprintf("%s|%s", *qIntermediateSession.Email, idformat.Organization.Format(*qIntermediateSession.OrganizationID)),
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
		RPID:              qProject.VaultDomain,
		AttestationObject: req.AttestationObject,
	})
	if err != nil {
		return nil, fmt.Errorf("parse webauthn credential: %w", err)
	}

	publicKey, err := x509.MarshalPKIXPublicKey(cred.PublicKey)
	if err != nil {
		return nil, fmt.Errorf("marshal public key: %w", err)
	}

	if _, err := q.UpdateIntermediateSessionRegisterPasskey(ctx, queries.UpdateIntermediateSessionRegisterPasskeyParams{
		ID:                  authn.IntermediateSessionID(ctx),
		PasskeyCredentialID: cred.ID,
		PasskeyPublicKey:    publicKey,
		PasskeyAaguid:       &cred.AAGUID,
		PasskeyRpID:         &req.RpId,
	}); err != nil {
		return nil, fmt.Errorf("register passkey: %w", err)
	}

	if err := commit(); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}

	return &intermediatev1.RegisterPasskeyResponse{}, nil
}

func (s *Store) IssuePasskeyChallenge(ctx context.Context, req *intermediatev1.IssuePasskeyChallengeRequest) (*intermediatev1.IssuePasskeyChallengeResponse, error) {
	_, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	qProject, err := q.GetProjectByID(ctx, authn.ProjectID(ctx))
	if err != nil {
		return nil, fmt.Errorf("get project by id: %w", err)
	}

	qIntermediateSession, err := q.GetIntermediateSessionByID(ctx, authn.IntermediateSessionID(ctx))
	if err != nil {
		return nil, fmt.Errorf("get intermediate session by id: %w", err)
	}

	qOrg, err := q.GetProjectOrganizationByID(ctx, queries.GetProjectOrganizationByIDParams{
		ProjectID: authn.ProjectID(ctx),
		ID:        *qIntermediateSession.OrganizationID,
	})
	if err != nil {
		return nil, fmt.Errorf("get organization by id: %w", err)
	}

	qMatchingUser, err := s.matchUser(ctx, q, qOrg, qIntermediateSession)
	if err != nil {
		return nil, fmt.Errorf("match user: %w", err)
	}

	credentialIDs, err := q.GetUserPasskeyCredentialIDs(ctx, qMatchingUser.ID)
	if err != nil {
		return nil, fmt.Errorf("get user passkey credential ids: %w", err)
	}

	var challenge [32]byte
	if _, err := rand.Read(challenge[:]); err != nil {
		return nil, fmt.Errorf("read random bytes: %w", err)
	}

	challengeSHA256 := sha256.Sum256(challenge[:])
	if _, err := q.UpdateIntermediateSessionPasskeyVerifyChallengeSHA256(ctx, queries.UpdateIntermediateSessionPasskeyVerifyChallengeSHA256Params{
		ID:                           authn.IntermediateSessionID(ctx),
		PasskeyVerifyChallengeSha256: challengeSHA256[:],
	}); err != nil {
		return nil, fmt.Errorf("update intermediate session passkey verify challenge: %w", err)
	}

	if err := commit(); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}

	return &intermediatev1.IssuePasskeyChallengeResponse{
		RpId:          qProject.VaultDomain,
		CredentialIds: credentialIDs,
		Challenge:     challenge[:],
	}, nil
}

func (s *Store) VerifyPasskey(ctx context.Context, req *intermediatev1.VerifyPasskeyRequest) (*intermediatev1.VerifyPasskeyResponse, error) {
	_, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	qIntermediateSession, err := q.GetIntermediateSessionByID(ctx, authn.IntermediateSessionID(ctx))
	if err != nil {
		return nil, fmt.Errorf("get intermediate session by id: %w", err)
	}

	qOrg, err := q.GetProjectOrganizationByID(ctx, queries.GetProjectOrganizationByIDParams{
		ProjectID: authn.ProjectID(ctx),
		ID:        *qIntermediateSession.OrganizationID,
	})
	if err != nil {
		return nil, fmt.Errorf("get organization by id: %w", err)
	}

	qMatchingUser, err := s.matchUser(ctx, q, qOrg, qIntermediateSession)
	if err != nil {
		return nil, fmt.Errorf("match user: %w", err)
	}

	qPasskey, err := q.GetPasskeyByCredentialID(ctx, queries.GetPasskeyByCredentialIDParams{
		CredentialID: req.CredentialId,
		UserID:       qMatchingUser.ID,
	})
	if err != nil {
		return nil, fmt.Errorf("get passkey by credential id: %w", err)
	}

	qTrustedDomains, err := q.GetProjectTrustedDomains(ctx, authn.ProjectID(ctx))
	if err != nil {
		return nil, fmt.Errorf("get project trusted domains: %w", err)
	}

	// the set of origins we expect passkeys to be using
	var origins []string
	for _, qProjectTrustedDomain := range qTrustedDomains {
		origins = append(origins, fmt.Sprintf("https://%s", qProjectTrustedDomain.Domain))
	}

	publicKey, err := x509.ParsePKIXPublicKey(qPasskey.PublicKey)
	if err != nil {
		return nil, fmt.Errorf("marshal public key: %w", err)
	}

	credential := webauthn.Credential{
		PublicKey: publicKey,
	}

	if err := credential.Verify(&webauthn.VerifyRequest{
		RPID:              qPasskey.RpID,
		ChallengeSHA256:   qIntermediateSession.PasskeyVerifyChallengeSha256,
		ClientDataJSON:    req.ClientDataJson,
		AuthenticatorData: req.AuthenticatorData,
		Signature:         req.Signature,
		Origins:           origins,
	}); err != nil {
		return nil, apierror.NewInvalidArgumentError("invalid passkey verification", fmt.Errorf("verify passkey: %w", err))
	}

	if _, err := q.UpdateIntermediateSessionPasskeyVerified(ctx, authn.IntermediateSessionID(ctx)); err != nil {
		return nil, fmt.Errorf("update intermediate session passkey verified: %w", err)
	}

	if err := commit(); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}

	return &intermediatev1.VerifyPasskeyResponse{}, nil
}

func (s *Store) checkShouldRegisterPasskey(ctx context.Context, q *queries.Queries) error {
	// don't register passkeys if you're already matching a user, and that user
	// has at least one active passkey

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

	// does the matching user have any active passkeys?
	hasPasskeys, err := q.GetUserHasActivePasskey(ctx, qUser.ID)
	if err != nil {
		return fmt.Errorf("get user has passkey: %w", err)
	}

	if hasPasskeys {
		return apierror.NewFailedPreconditionError("user already has active passkeys", nil)
	}
	return nil
}
