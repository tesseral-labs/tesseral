package store

import (
	"bytes"
	"context"
	"crypto/sha256"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go-v2/service/kms/types"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/tesseral-labs/tesseral/internal/common/apierror"
	"github.com/tesseral-labs/tesseral/internal/intermediate/authn"
	intermediatev1 "github.com/tesseral-labs/tesseral/internal/intermediate/gen/tesseral/intermediate/v1"
	"github.com/tesseral-labs/tesseral/internal/intermediate/store/queries"
	"github.com/tesseral-labs/tesseral/internal/microsoftoauth"
)

func (s *Store) GetMicrosoftOAuthRedirectURL(ctx context.Context, req *intermediatev1.GetMicrosoftOAuthRedirectURLRequest) (*intermediatev1.GetMicrosoftOAuthRedirectURLResponse, error) {
	_, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	qProject, err := q.GetProjectByID(ctx, authn.ProjectID(ctx))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apierror.NewNotFoundError("project not found", fmt.Errorf("get project by id: %v", err))
		}

		return nil, fmt.Errorf("get project by id: %v", err)
	}

	if err := enforceProjectLoginEnabled(qProject); err != nil {
		return nil, fmt.Errorf("enforce project login enabled: %w", err)
	}

	clientID := s.defaultMicrosoftOAuthClientID
	redirectURI := s.defaultMicrosoftOAuthRedirectURI

	if derefOrEmpty(qProject.MicrosoftOauthClientID) != "" {
		clientID = *qProject.MicrosoftOauthClientID
		redirectURI = req.RedirectUrl
	}

	state := uuid.NewString()
	stateSHA := sha256.Sum256([]byte(state))
	if _, err := q.UpdateIntermediateSessionMicrosoftOAuthStateSHA256(ctx, queries.UpdateIntermediateSessionMicrosoftOAuthStateSHA256Params{
		ID:                        authn.IntermediateSessionID(ctx),
		MicrosoftOauthStateSha256: stateSHA[:],
	}); err != nil {
		return nil, fmt.Errorf("update intermediate session microsoft oauth state: %v", err)
	}

	url := microsoftoauth.GetAuthorizeURL(&microsoftoauth.GetAuthorizeURLRequest{
		MicrosoftOAuthClientID: clientID,
		RedirectURI:            redirectURI,
		State:                  state,
	})

	if err := commit(); err != nil {
		return nil, fmt.Errorf("commit: %v", err)
	}

	return &intermediatev1.GetMicrosoftOAuthRedirectURLResponse{
		Url: url,
	}, nil
}

func (s *Store) RedeemMicrosoftOAuthCode(ctx context.Context, req *intermediatev1.RedeemMicrosoftOAuthCodeRequest) (*intermediatev1.RedeemMicrosoftOAuthCodeResponse, error) {
	qProject, qIntermediateSession, err := s.getProjectAndIntermediateSession(ctx)
	if err != nil {
		return nil, fmt.Errorf("get project and intermediate session: %v", err)
	}

	if err := enforceProjectLoginEnabled(*qProject); err != nil {
		return nil, fmt.Errorf("enforce project login enabled: %w", err)
	}

	clientID := s.defaultMicrosoftOAuthClientID
	clientSecret := s.defaultMicrosoftOAuthClientSecret
	redirectURI := s.defaultMicrosoftOAuthRedirectURI

	if derefOrEmpty(qProject.MicrosoftOauthClientID) != "" {
		clientID = *qProject.MicrosoftOauthClientID
		redirectURI = req.RedirectUrl

		decryptRes, err := s.kms.Decrypt(ctx, &kms.DecryptInput{
			CiphertextBlob:      qProject.MicrosoftOauthClientSecretCiphertext,
			EncryptionAlgorithm: types.EncryptionAlgorithmSpecRsaesOaepSha256,
			KeyId:               &s.microsoftOAuthClientSecretsKMSKeyID,
		})
		if err != nil {
			return nil, fmt.Errorf("decrypt microsoft oauth client secret: %v", err)
		}

		clientSecret = string(decryptRes.Plaintext)
	}

	stateSHA := sha256.Sum256([]byte(req.State))
	if !bytes.Equal(qIntermediateSession.MicrosoftOauthStateSha256, stateSHA[:]) {
		return nil, apierror.NewInvalidArgumentError("invalid state", fmt.Errorf("invalid state"))
	}

	redeemRes, err := s.microsoftOAuthClient.RedeemCode(ctx, &microsoftoauth.RedeemCodeRequest{
		MicrosoftOAuthClientID:     clientID,
		MicrosoftOAuthClientSecret: clientSecret,
		RedirectURI:                redirectURI,
		Code:                       req.Code,
	})
	if err != nil {
		return nil, apierror.NewInvalidArgumentError("failed to redeem microsoft oauth code", fmt.Errorf("redeem microsoft oauth code: %v", err))
	}

	if qIntermediateSession.Email != nil && redeemRes.Email != *qIntermediateSession.Email {
		return nil, apierror.NewInvalidArgumentError("Email mismatch", fmt.Errorf("email mismatch"))
	}

	// start new tx now that all kms and oauth http i/o is done
	_, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	if _, err := q.UpdateIntermediateSessionMicrosoftDetails(ctx, queries.UpdateIntermediateSessionMicrosoftDetailsParams{
		ID:                authn.IntermediateSessionID(ctx),
		Email:             &redeemRes.Email,
		MicrosoftUserID:   &redeemRes.MicrosoftUserID,
		MicrosoftTenantID: refOrNil(redeemRes.MicrosoftTenantID),
	}); err != nil {
		return nil, fmt.Errorf("update intermediate session microsoft details: %v", err)
	}

	primaryAuthFactor := queries.PrimaryAuthFactorMicrosoft
	if _, err := q.UpdateIntermediateSessionPrimaryAuthFactor(ctx, queries.UpdateIntermediateSessionPrimaryAuthFactorParams{
		ID:                authn.IntermediateSessionID(ctx),
		PrimaryAuthFactor: &primaryAuthFactor,
	}); err != nil {
		return nil, fmt.Errorf("update intermediate session primary auth factor: %v", err)
	}

	if err := commit(); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}

	return &intermediatev1.RedeemMicrosoftOAuthCodeResponse{}, nil
}
