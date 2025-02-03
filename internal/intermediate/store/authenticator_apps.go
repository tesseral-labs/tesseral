package store

import (
	"context"
	"crypto/rand"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go-v2/service/kms/types"
	"github.com/openauth/openauth/internal/bcryptcost"
	"github.com/openauth/openauth/internal/common/apierror"
	"github.com/openauth/openauth/internal/intermediate/authn"
	intermediatev1 "github.com/openauth/openauth/internal/intermediate/gen/openauth/intermediate/v1"
	"github.com/openauth/openauth/internal/intermediate/store/queries"
	"github.com/openauth/openauth/internal/totp"
	"golang.org/x/crypto/bcrypt"
)

const (
	// how many backup codes to generate for an authenticator app
	backupCodeCount = 10

	// after this many failed attempts, lock out a user
	backupCodeLockoutAttempts = 5

	// how long to lock users out
	backupCodeLockoutDuration = time.Minute * 10
)

func (s *Store) GetAuthenticatorAppOptions(ctx context.Context, req *intermediatev1.GetAuthenticatorAppOptionsRequest) (*intermediatev1.GetAuthenticatorAppOptionsResponse, error) {
	if err := s.checkShouldRegisterAuthenticatorApp(ctx); err != nil {
		return nil, fmt.Errorf("check should register authenticator app: %w", err)
	}

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

	if _, err := q.UpdateIntermediateSessionAuthenticatorAppSecretCiphertext(ctx, queries.UpdateIntermediateSessionAuthenticatorAppSecretCiphertextParams{
		ID:                               authn.IntermediateSessionID(ctx),
		AuthenticatorAppSecretCiphertext: encryptRes.CiphertextBlob,
	}); err != nil {
		return nil, fmt.Errorf("update intermediate session authenticator app secret ciphertext: %w", err)
	}

	if err := commit(); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}

	return &intermediatev1.GetAuthenticatorAppOptionsResponse{
		Secret: secret[:],
	}, nil
}

func (s *Store) RegisterAuthenticatorApp(ctx context.Context, req *intermediatev1.RegisterAuthenticatorAppRequest) (*intermediatev1.RegisterAuthenticatorAppResponse, error) {
	if err := s.checkShouldRegisterAuthenticatorApp(ctx); err != nil {
		return nil, fmt.Errorf("check should register authenticator app: %w", err)
	}

	secret, err := s.getAuthenticatorAppSecret(ctx)
	if err != nil {
		return nil, fmt.Errorf("get authenticator app secret: %w", err)
	}

	key := totp.Key{Secret: secret}
	if err := key.Validate(time.Now(), req.TotpCode); err != nil {
		return nil, apierror.NewInvalidArgumentError("invalid totp code", err)
	}

	_, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	var backupCodes []string
	var backupCodeBcrypts [][]byte
	for i := 0; i < backupCodeCount; i++ {
		var code [8]byte
		if _, err := rand.Read(code[:]); err != nil {
			return nil, fmt.Errorf("read random bytes: %w", err)
		}

		codeFormatted := fmt.Sprintf("%x-%x-%x-%x", code[:2], code[2:4], code[4:6], code[6:])

		codeBcrypt, err := bcrypt.GenerateFromPassword([]byte(codeFormatted), bcryptcost.Cost)
		if err != nil {
			return nil, fmt.Errorf("generate bcrypt hash: %w", err)
		}

		backupCodes = append(backupCodes, codeFormatted)
		backupCodeBcrypts = append(backupCodeBcrypts, codeBcrypt)
	}

	if _, err := q.UpdateIntermediateSessionAuthenticatorAppVerified(ctx, authn.IntermediateSessionID(ctx)); err != nil {
		return nil, fmt.Errorf("update intermediate session authenticator app verified: %w", err)
	}

	if _, err := q.UpdateIntermediateSessionAuthenticatorAppBackupCodeBcrypts(ctx, queries.UpdateIntermediateSessionAuthenticatorAppBackupCodeBcryptsParams{
		ID:                                authn.IntermediateSessionID(ctx),
		AuthenticatorAppBackupCodeBcrypts: backupCodeBcrypts,
	}); err != nil {
		return nil, fmt.Errorf("update intermediate session authenticator app backup code bcrypts: %w", err)
	}

	if err := commit(); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}

	return &intermediatev1.RegisterAuthenticatorAppResponse{
		RecoveryCodes: backupCodes,
	}, nil
}

func (s *Store) VerifyAuthenticatorApp(ctx context.Context, req *intermediatev1.VerifyAuthenticatorAppRequest) (*intermediatev1.VerifyAuthenticatorAppResponse, error) {
	if req.RecoveryCode != "" {
		if err := s.verifyAuthenticatorAppByBackupCode(ctx, req.RecoveryCode); err != nil {
			return nil, fmt.Errorf("verify authenticator app by backup code: %w", err)
		}
	} else {
		if err := s.verifyAuthenticatorAppByTOTPCode(ctx, req.TotpCode); err != nil {
			return nil, fmt.Errorf("verify authenticator app by totp code: %w", err)
		}
	}

	_, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	if _, err := q.UpdateIntermediateSessionAuthenticatorAppVerified(ctx, authn.IntermediateSessionID(ctx)); err != nil {
		return nil, fmt.Errorf("update intermediate session authenticator app verified: %w", err)
	}

	if err := commit(); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}

	return &intermediatev1.VerifyAuthenticatorAppResponse{}, nil
}

func (s *Store) verifyAuthenticatorAppByTOTPCode(ctx context.Context, totpCode string) error {
	secret, err := s.getAuthenticatorAppSecret(ctx)
	if err != nil {
		return fmt.Errorf("get authenticator app secret: %w", err)
	}

	key := totp.Key{Secret: secret}
	if err := key.Validate(time.Now(), totpCode); err != nil {
		return apierror.NewInvalidArgumentError("invalid totp code", err)
	}

	return nil
}

func (s *Store) verifyAuthenticatorAppByBackupCode(ctx context.Context, backupCode string) error {
	_, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return err
	}
	defer rollback()

	qIntermediateSession, err := q.GetIntermediateSessionByID(ctx, authn.IntermediateSessionID(ctx))
	if err != nil {
		return fmt.Errorf("get intermediate session by id: %w", err)
	}

	qOrg, err := q.GetProjectOrganizationByID(ctx, queries.GetProjectOrganizationByIDParams{
		ProjectID: authn.ProjectID(ctx),
		ID:        *qIntermediateSession.OrganizationID,
	})

	qMatchingUser, err := s.matchUser(ctx, q, qOrg, qIntermediateSession)
	if err != nil {
		return fmt.Errorf("match user: %w", err)
	}

	if qMatchingUser.AuthenticatorAppBackupCodeLockoutExpireTime != nil && qMatchingUser.AuthenticatorAppBackupCodeLockoutExpireTime.After(time.Now()) {
		return apierror.NewFailedPreconditionError("too many backup code attempts attempts; user is temporarily locked out", nil)
	}

	var ok bool
	var backupCodeBcrypts [][]byte
	for _, b := range qIntermediateSession.AuthenticatorAppBackupCodeBcrypts {
		if bcrypt.CompareHashAndPassword(b, []byte(backupCode)) == nil {
			ok = true
			continue // do not keep this bcrypt around; the code is used
		}

		backupCodeBcrypts = append(backupCodeBcrypts, b)
	}

	// update backup code attempt / lockout state
	if ok {
		// success; reset fail count
		if _, err := q.UpdateUserFailedAuthenticatorAppBackupCodeAttempts(ctx, queries.UpdateUserFailedAuthenticatorAppBackupCodeAttemptsParams{
			ID:                                       qMatchingUser.ID,
			FailedAuthenticatorAppBackupCodeAttempts: aws.Int32(0),
		}); err != nil {
			return fmt.Errorf("update user failed authenticator app backup code attempts: %w", err)
		}
	} else {
		// failure; bump fail count and check to see if we should lock out
		attempts := qMatchingUser.FailedPasswordAttempts + 1
		if attempts >= backupCodeLockoutAttempts {
			// lock the user out
			backupCodeLockoutExpireTime := time.Now().Add(backupCodeLockoutDuration)
			if _, err := q.UpdateUserAuthenticatorAppBackupCodeLockoutExpireTime(ctx, queries.UpdateUserAuthenticatorAppBackupCodeLockoutExpireTimeParams{
				ID: qMatchingUser.ID,
				AuthenticatorAppBackupCodeLockoutExpireTime: &backupCodeLockoutExpireTime,
			}); err != nil {
				return fmt.Errorf("update user authenticator app backup code lockout expire time: %w", err)
			}

			// reset fail count
			if _, err := q.UpdateUserFailedAuthenticatorAppBackupCodeAttempts(ctx, queries.UpdateUserFailedAuthenticatorAppBackupCodeAttemptsParams{
				ID:                                       qMatchingUser.ID,
				FailedAuthenticatorAppBackupCodeAttempts: aws.Int32(0),
			}); err != nil {
				return fmt.Errorf("update user failed authenticator app backup code attempts: %w", err)
			}

			if err := commit(); err != nil {
				return fmt.Errorf("commit: %w", err)
			}

			return apierror.NewFailedPreconditionError("too many backup code attempts attempts; user is temporarily locked out", nil)
		}

		// update fail count, but do not lock out
		if _, err := q.UpdateUserFailedAuthenticatorAppBackupCodeAttempts(ctx, queries.UpdateUserFailedAuthenticatorAppBackupCodeAttemptsParams{
			ID:                                       qMatchingUser.ID,
			FailedAuthenticatorAppBackupCodeAttempts: &attempts,
		}); err != nil {
			return fmt.Errorf("update user failed password attempts: %w", err)
		}

		return apierror.NewInvalidArgumentError("invalid backup code", nil)
	}

	if _, err := q.UpdateIntermediateSessionAuthenticatorAppBackupCodeBcrypts(ctx, queries.UpdateIntermediateSessionAuthenticatorAppBackupCodeBcryptsParams{
		ID:                                authn.IntermediateSessionID(ctx),
		AuthenticatorAppBackupCodeBcrypts: backupCodeBcrypts,
	}); err != nil {
		return fmt.Errorf("update intermediate session authenticator app backup code bcrypts: %w", err)
	}

	if err := commit(); err != nil {
		return fmt.Errorf("commit: %w", err)
	}

	return nil
}

func (s *Store) getAuthenticatorAppSecret(ctx context.Context) ([]byte, error) {
	_, q, _, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	qIntermediateSession, err := q.GetIntermediateSessionByID(ctx, authn.IntermediateSessionID(ctx))
	if err != nil {
		return nil, fmt.Errorf("get intermediate session by id: %w", err)
	}

	decryptRes, err := s.kms.Decrypt(ctx, &kms.DecryptInput{
		KeyId:               &s.authenticatorAppSecretsKMSKeyID,
		EncryptionAlgorithm: types.EncryptionAlgorithmSpecRsaesOaepSha256,
		CiphertextBlob:      qIntermediateSession.AuthenticatorAppSecretCiphertext,
	})
	if err != nil {
		return nil, fmt.Errorf("decrypt authenticator app secret ciphertext: %w", err)
	}

	return decryptRes.Plaintext, nil
}

func (s *Store) checkShouldRegisterAuthenticatorApp(ctx context.Context) error {
	// don't register an authenticator app if you're already matching a user,
	// and that user has one

	_, q, _, rollback, err := s.tx(ctx)
	if err != nil {
		return err
	}
	defer rollback()

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

	if qUser.AuthenticatorAppSecretCiphertext != nil {
		return apierror.NewFailedPreconditionError("user already has an authenticator app", nil)
	}
	return nil
}
