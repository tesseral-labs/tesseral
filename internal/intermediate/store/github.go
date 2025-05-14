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
	"github.com/tesseral-labs/tesseral/internal/githuboauth"
	"github.com/tesseral-labs/tesseral/internal/intermediate/authn"
	intermediatev1 "github.com/tesseral-labs/tesseral/internal/intermediate/gen/tesseral/intermediate/v1"
	"github.com/tesseral-labs/tesseral/internal/intermediate/store/queries"
)

func (s *Store) GetGithubOAuthRedirectURL(ctx context.Context, req *intermediatev1.GetGithubOAuthRedirectURLRequest) (*intermediatev1.GetGithubOAuthRedirectURLResponse, error) {
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

	if qProject.GithubOauthClientID == nil {
		return nil, apierror.NewFailedPreconditionError("github oauth client id not set", fmt.Errorf("github oauth client id not set"))
	}

	state := uuid.NewString()
	stateSHA := sha256.Sum256([]byte(state))
	if _, err := q.UpdateIntermediateSessionGithubOAuthStateSHA256(ctx, queries.UpdateIntermediateSessionGithubOAuthStateSHA256Params{
		ID:                     authn.IntermediateSessionID(ctx),
		GithubOauthStateSha256: stateSHA[:],
	}); err != nil {
		return nil, fmt.Errorf("update intermediate session github oauth state: %v", err)
	}

	url := githuboauth.GetAuthorizeURL(&githuboauth.GetAuthorizeURLRequest{
		GithubOAuthClientID: *qProject.GithubOauthClientID,
		RedirectURI:         req.RedirectUrl,
		State:               state,
	})

	if err := commit(); err != nil {
		return nil, fmt.Errorf("commit: %v", err)
	}

	return &intermediatev1.GetGithubOAuthRedirectURLResponse{
		Url: url,
	}, nil
}

func (s *Store) RedeemGithubOAuthCode(ctx context.Context, req *intermediatev1.RedeemGithubOAuthCodeRequest) (*intermediatev1.RedeemGithubOAuthCodeResponse, error) {
	qProject, qIntermediateSession, err := s.getProjectAndIntermediateSession(ctx)
	if err != nil {
		return nil, fmt.Errorf("get project and intermediate session: %v", err)
	}

	if err := enforceProjectLoginEnabled(*qProject); err != nil {
		return nil, fmt.Errorf("enforce project login enabled: %w", err)
	}

	if qProject.GithubOauthClientID == nil || qProject.GithubOauthClientSecretCiphertext == nil {
		return nil, apierror.NewFailedPreconditionError("github oauth client id or secret not set", fmt.Errorf("github oauth client id or secret not set"))
	}

	stateSHA := sha256.Sum256([]byte(req.State))
	if !bytes.Equal(qIntermediateSession.GithubOauthStateSha256, stateSHA[:]) {
		return nil, apierror.NewInvalidArgumentError("invalid state", fmt.Errorf("invalid state"))
	}

	decryptRes, err := s.kms.Decrypt(ctx, &kms.DecryptInput{
		CiphertextBlob:      qProject.GithubOauthClientSecretCiphertext,
		EncryptionAlgorithm: types.EncryptionAlgorithmSpecRsaesOaepSha256,
		KeyId:               &s.githubOAuthClientSecretsKMSKeyID,
	})
	if err != nil {
		return nil, fmt.Errorf("decrypt github oauth client secret: %v", err)
	}

	redeemRes, err := s.githubOAuthClient.RedeemCode(ctx, &githuboauth.RedeemCodeRequest{
		GithubOAuthClientID:     *qProject.GithubOauthClientID,
		GithubOAuthClientSecret: string(decryptRes.Plaintext),
		RedirectURI:             req.RedirectUrl,
		Code:                    req.Code,
	})
	if err != nil {
		return nil, apierror.NewInvalidArgumentError("failed to redeem github oauth code", fmt.Errorf("redeem github oauth code: %v", err))
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

	if _, err := q.UpdateIntermediateSessionGithubDetails(ctx, queries.UpdateIntermediateSessionGithubDetailsParams{
		ID:                authn.IntermediateSessionID(ctx),
		Email:             &redeemRes.Email,
		GithubUserID:      &redeemRes.GithubUserID,
		UserDisplayName:   refOrNil(redeemRes.DisplayName),
		ProfilePictureUrl: refOrNil(redeemRes.ProfilePictureURL),
	}); err != nil {
		return nil, fmt.Errorf("update intermediate session github details: %v", err)
	}

	if redeemRes.EmailVerified {
		if _, err := q.CreateVerifiedEmail(ctx, queries.CreateVerifiedEmailParams{
			ID:           uuid.New(),
			ProjectID:    authn.ProjectID(ctx),
			Email:        redeemRes.Email,
			GithubUserID: &redeemRes.GithubUserID,
		}); err != nil {
			return nil, fmt.Errorf("create verified email: %v", err)
		}
	}

	primaryAuthFactor := queries.PrimaryAuthFactorGithub
	if _, err := q.UpdateIntermediateSessionPrimaryAuthFactor(ctx, queries.UpdateIntermediateSessionPrimaryAuthFactorParams{
		ID:                authn.IntermediateSessionID(ctx),
		PrimaryAuthFactor: &primaryAuthFactor,
	}); err != nil {
		return nil, fmt.Errorf("update intermediate session primary auth factor: %v", err)
	}

	if err := commit(); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}

	return &intermediatev1.RedeemGithubOAuthCodeResponse{}, nil
}
