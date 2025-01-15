package store

import (
	"bytes"
	"context"
	"crypto/sha256"
	"fmt"
	"time"

	"connectrpc.com/connect"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go-v2/service/kms/types"
	"github.com/google/uuid"
	"github.com/openauth/openauth/internal/errorcodes"
	"github.com/openauth/openauth/internal/intermediate/authn"
	intermediatev1 "github.com/openauth/openauth/internal/intermediate/gen/openauth/intermediate/v1"
	"github.com/openauth/openauth/internal/intermediate/store/queries"
	"github.com/openauth/openauth/internal/microsoftoauth"
	"github.com/openauth/openauth/internal/store/idformat"
)

func (s *Store) GetMicrosoftOAuthRedirectURL(ctx context.Context, req *intermediatev1.GetMicrosoftOAuthRedirectURLRequest) (*intermediatev1.GetMicrosoftOAuthRedirectURLResponse, error) {
	_, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	qProject, err := q.GetProjectByID(ctx, authn.ProjectID(ctx))
	if err != nil {
		return nil, fmt.Errorf("get project by id: %v", err)
	}

	if qProject.MicrosoftOauthClientID == nil {
		return nil, connect.NewError(connect.CodeFailedPrecondition, errorcodes.NewFailedPreconditionError())
	}

	token := uuid.New()
	tokenSha256 := sha256.Sum256(token[:])
	expiresAt := time.Now().Add(155 * time.Minute)

	// Since this is the entrypoint for the google oauth flow, we create the intermediate session here
	intermediateSession, err := q.CreateIntermediateSession(ctx, queries.CreateIntermediateSessionParams{
		ID:          uuid.Must(uuid.NewV7()),
		ProjectID:   authn.ProjectID(ctx),
		ExpireTime:  &expiresAt,
		TokenSha256: tokenSha256[:],
	})
	if err != nil {
		return nil, fmt.Errorf("create intermediate session: %v", err)
	}

	state := uuid.NewString()
	stateSHA := sha256.Sum256([]byte(state))
	if _, err := q.UpdateIntermediateSessionMicrosoftOAuthStateSHA256(ctx, queries.UpdateIntermediateSessionMicrosoftOAuthStateSHA256Params{
		ID:                        intermediateSession.ID,
		MicrosoftOauthStateSha256: stateSHA[:],
	}); err != nil {
		return nil, fmt.Errorf("update intermediate session microsoft oauth state: %v", err)
	}

	url := microsoftoauth.GetAuthorizeURL(&microsoftoauth.GetAuthorizeURLRequest{
		MicrosoftOAuthClientID: *qProject.MicrosoftOauthClientID,
		RedirectURI:            req.RedirectUrl,
		State:                  state,
	})

	if err := commit(); err != nil {
		return nil, fmt.Errorf("commit: %v", err)
	}

	return &intermediatev1.GetMicrosoftOAuthRedirectURLResponse{
		IntermediateSessionToken: idformat.IntermediateSessionToken.Format(token),
		Url:                      url,
	}, nil
}

func (s *Store) RedeemMicrosoftOAuthCode(ctx context.Context, req *intermediatev1.RedeemMicrosoftOAuthCodeRequest) (*intermediatev1.RedeemMicrosoftOAuthCodeResponse, error) {
	qProject, qIntermediateSession, err := s.getProjectAndIntermediateSession(ctx)
	if err != nil {
		return nil, fmt.Errorf("get project and intermediate session: %v", err)
	}

	if qProject.MicrosoftOauthClientID == nil || qProject.MicrosoftOauthClientSecretCiphertext == nil {
		return nil, connect.NewError(connect.CodeFailedPrecondition, errorcodes.NewFailedPreconditionError())
	}

	stateSHA := sha256.Sum256([]byte(req.State))
	if !bytes.Equal(qIntermediateSession.MicrosoftOauthStateSha256, stateSHA[:]) {
		return nil, connect.NewError(connect.CodeInvalidArgument, errorcodes.NewFailedPreconditionError())
	}

	decryptRes, err := s.kms.Decrypt(ctx, &kms.DecryptInput{
		CiphertextBlob:      qProject.MicrosoftOauthClientSecretCiphertext,
		EncryptionAlgorithm: types.EncryptionAlgorithmSpecRsaesOaepSha256,
		KeyId:               &s.microsoftOAuthClientSecretsKMSKeyID,
	})
	if err != nil {
		return nil, fmt.Errorf("decrypt microsoft oauth client secret: %v", err)
	}

	// start new tx now that all kms and oauth http i/o is done
	_, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	redeemRes, err := s.microsoftOAuthClient.RedeemCode(ctx, &microsoftoauth.RedeemCodeRequest{
		MicrosoftOAuthClientID:     *qProject.MicrosoftOauthClientID,
		MicrosoftOAuthClientSecret: string(decryptRes.Plaintext),
		RedirectURI:                req.RedirectUrl,
		Code:                       req.Code,
	})
	if err != nil {
		return nil, fmt.Errorf("redeem microsoft oauth code: %v", err)
	}

	// todo what if the intermediate session already has an email
	// todo what if redeem comes back with the "public" microsoft tenant ID?
	if _, err := q.UpdateIntermediateSessionMicrosoftDetails(ctx, queries.UpdateIntermediateSessionMicrosoftDetailsParams{
		ID:                authn.IntermediateSessionID(ctx),
		Email:             &redeemRes.Email,
		MicrosoftUserID:   &redeemRes.MicrosoftUserID,
		MicrosoftTenantID: &redeemRes.MicrosoftTenantID,
	}); err != nil {
		return nil, fmt.Errorf("update intermediate session microsoft details: %v", err)
	}

	if err := commit(); err != nil {
		return nil, err
	}

	return &intermediatev1.RedeemMicrosoftOAuthCodeResponse{}, nil
}
