package store

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go-v2/service/kms/types"
	"github.com/google/uuid"
	"github.com/tesseral-labs/tesseral/internal/common/apierror"
	"github.com/tesseral-labs/tesseral/internal/frontend/authn"
	frontendv1 "github.com/tesseral-labs/tesseral/internal/frontend/gen/tesseral/frontend/v1"
	"github.com/tesseral-labs/tesseral/internal/frontend/store/queries"
	"github.com/tesseral-labs/tesseral/internal/store/idformat"
	"github.com/tesseral-labs/tesseral/internal/totp"
)

// how many recovery codes to generate for an authenticator app
//
// keep in sync with intermediate/store/authenticator_apps.go
const recoveryCodeCount = 10

func (s *Store) GetAuthenticatorAppOptions(ctx context.Context) (*frontendv1.GetAuthenticatorAppOptionsResponse, error) {
	var secret [32]byte
	if _, err := rand.Read(secret[:]); err != nil {
		return nil, fmt.Errorf("read random bytes: %w", err)
	}

	encryptRes, err := s.kms.Encrypt(ctx, &kms.EncryptInput{
		KeyId:               &s.authenticatorAppSecretsKMSKeyID,
		Plaintext:           secret[:],
		EncryptionAlgorithm: types.EncryptionAlgorithmSpecRsaesOaepSha256,
	})
	if err != nil {
		return nil, fmt.Errorf("encrypt authenticator app secret: %w", err)
	}

	_, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	if _, err := q.CreateUserAuthenticatorAppChallenge(ctx, queries.CreateUserAuthenticatorAppChallengeParams{
		UserID:                           authn.UserID(ctx),
		AuthenticatorAppSecretCiphertext: encryptRes.CiphertextBlob,
	}); err != nil {
		return nil, fmt.Errorf("create user authenticator app challenge: %w", err)
	}

	qProject, err := q.GetProjectByID(ctx, authn.ProjectID(ctx))
	if err != nil {
		return nil, fmt.Errorf("get project by id: %w", err)
	}

	qUser, err := q.GetUserByID(ctx, authn.UserID(ctx))
	if err != nil {
		return nil, fmt.Errorf("get user by id: %w", err)
	}

	if err := commit(); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}

	key := totp.Key{Secret: secret[:]}
	return &frontendv1.GetAuthenticatorAppOptionsResponse{
		OtpauthUri: key.OTPAuthURI(qProject.DisplayName, qUser.Email),
	}, nil
}

func (s *Store) RegisterAuthenticatorApp(ctx context.Context, req *frontendv1.RegisterAuthenticatorAppRequest) (*frontendv1.RegisterAuthenticatorAppResponse, error) {
	secret, err := s.getUserAuthenticatorAppChallengeSecret(ctx)
	if err != nil {
		return nil, fmt.Errorf("get authenticator app secret: %w", err)
	}

	key := totp.Key{Secret: secret}
	if err := key.Validate(time.Now(), req.TotpCode); err != nil {
		return nil, apierror.NewInvalidTOTPCodeError("incorrect totp code", fmt.Errorf("validate totp code: %w", err))
	}

	var recoveryCodes []string
	var recoveryCodeSHA256s [][]byte
	for i := 0; i < recoveryCodeCount; i++ {
		recoveryCode := uuid.New()
		recoveryCodeSHA := sha256.Sum256(recoveryCode[:])

		recoveryCodes = append(recoveryCodes, idformat.AuthenticatorAppRecoveryCode.Format(recoveryCode))
		recoveryCodeSHA256s = append(recoveryCodeSHA256s, recoveryCodeSHA[:])
	}

	_, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	qUserAuthenticatorAppChallenge, err := q.GetUserAuthenticatorAppChallenge(ctx, authn.UserID(ctx))
	if err != nil {
		return nil, fmt.Errorf("get user authenticator app challenge: %w", err)
	}

	if err := q.DeleteUserAuthenticatorAppChallenge(ctx, authn.UserID(ctx)); err != nil {
		return nil, fmt.Errorf("delete user authenticator app challenge: %w", err)
	}

	qUser, err := q.UpdateUserAuthenticatorApp(ctx, queries.UpdateUserAuthenticatorAppParams{
		ID:                                  authn.UserID(ctx),
		AuthenticatorAppSecretCiphertext:    qUserAuthenticatorAppChallenge.AuthenticatorAppSecretCiphertext,
		AuthenticatorAppRecoveryCodeSha256s: recoveryCodeSHA256s,
	})
	if err != nil {
		return nil, fmt.Errorf("update user authenticator app: %w", err)
	}

	if _, err := s.logAuditEvent(ctx, q, logAuditEventParams{
		EventName: "tesseral.users.register_authenticator_app",
		EventDetails: &frontendv1.UserAuthenticatorAppRegistered{
			User: parseUser(qUser),
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

	return &frontendv1.RegisterAuthenticatorAppResponse{
		RecoveryCodes: recoveryCodes,
	}, nil
}

func (s *Store) getUserAuthenticatorAppChallengeSecret(ctx context.Context) ([]byte, error) {
	_, q, _, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	qUserAuthenticatorAppChallenge, err := q.GetUserAuthenticatorAppChallenge(ctx, authn.UserID(ctx))
	if err != nil {
		return nil, fmt.Errorf("get user authenticator app challenge: %w", err)
	}

	decryptRes, err := s.kms.Decrypt(ctx, &kms.DecryptInput{
		KeyId:               &s.authenticatorAppSecretsKMSKeyID,
		EncryptionAlgorithm: types.EncryptionAlgorithmSpecRsaesOaepSha256,
		CiphertextBlob:      qUserAuthenticatorAppChallenge.AuthenticatorAppSecretCiphertext,
	})
	if err != nil {
		return nil, fmt.Errorf("decrypt authenticator app secret ciphertext: %w", err)
	}

	return decryptRes.Plaintext, nil
}
