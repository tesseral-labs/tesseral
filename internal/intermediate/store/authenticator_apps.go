package store

import (
	"context"
	"crypto/rand"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go-v2/service/kms/types"
	"github.com/tesseral-labs/tesseral/internal/bcryptcost"
	"github.com/tesseral-labs/tesseral/internal/common/apierror"
	"github.com/tesseral-labs/tesseral/internal/intermediate/authn"
	intermediatev1 "github.com/tesseral-labs/tesseral/internal/intermediate/gen/tesseral/intermediate/v1"
	"github.com/tesseral-labs/tesseral/internal/intermediate/store/queries"
	"github.com/tesseral-labs/tesseral/internal/totp"
	"golang.org/x/crypto/bcrypt"
)

const (
	// how many recovery codes to generate for an authenticator app
	recoveryCodeCount = 10

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

	qProject, err := q.GetProjectByID(ctx, authn.ProjectID(ctx))
	if err != nil {
		return nil, fmt.Errorf("get project by id: %w", err)
	}

	qOrg, err := q.GetProjectOrganizationByID(ctx, queries.GetProjectOrganizationByIDParams{
		ProjectID: authn.ProjectID(ctx),
		ID:        *qProject.OrganizationID,
	})
	if err != nil {
		return nil, fmt.Errorf("get organization by id: %w", err)
	}

	qIntermediateSession, err := q.GetIntermediateSessionByID(ctx, authn.IntermediateSessionID(ctx))
	if err != nil {
		return nil, fmt.Errorf("get intermediate session by id: %w", err)
	}

	qMatchingUser, err := s.matchUser(ctx, q, qOrg, qIntermediateSession)
	if err != nil {
		return nil, fmt.Errorf("match user: %w", err)
	}

	if _, err := q.UpdateIntermediateSessionAuthenticatorAppSecretCiphertext(ctx, queries.UpdateIntermediateSessionAuthenticatorAppSecretCiphertextParams{
		ID:                               authn.IntermediateSessionID(ctx),
		AuthenticatorAppSecretCiphertext: encryptRes.CiphertextBlob,
	}); err != nil {
		return nil, fmt.Errorf("update intermediate session authenticator app secret ciphertext: %w", err)
	}

	if err := commit(); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}

	key := totp.Key{Secret: secret[:]}
	return &intermediatev1.GetAuthenticatorAppOptionsResponse{
		OtpauthUri: key.OTPAuthURI(qProject.DisplayName, qMatchingUser.Email),
	}, nil
}

func (s *Store) RegisterAuthenticatorApp(ctx context.Context, req *intermediatev1.RegisterAuthenticatorAppRequest) (*intermediatev1.RegisterAuthenticatorAppResponse, error) {
	if err := s.checkShouldRegisterAuthenticatorApp(ctx); err != nil {
		return nil, fmt.Errorf("check should register authenticator app: %w", err)
	}

	secret, err := s.getIntermediateSessionPendingAuthenticatorAppSecret(ctx)
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
	var recoveryCodeBcrypts [][]byte
	for i := 0; i < recoveryCodeCount; i++ {
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
		recoveryCodeBcrypts = append(recoveryCodeBcrypts, codeBcrypt)
	}

	if _, err := q.UpdateIntermediateSessionAuthenticatorAppVerified(ctx, authn.IntermediateSessionID(ctx)); err != nil {
		return nil, fmt.Errorf("update intermediate session authenticator app verified: %w", err)
	}

	if _, err := q.UpdateIntermediateSessionAuthenticatorAppBackupCodeBcrypts(ctx, queries.UpdateIntermediateSessionAuthenticatorAppBackupCodeBcryptsParams{
		ID:                                  authn.IntermediateSessionID(ctx),
		AuthenticatorAppRecoveryCodeBcrypts: recoveryCodeBcrypts,
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
	if err := s.checkAuthenticatorAppLockedOut(ctx); err != nil {
		return nil, fmt.Errorf("check authenticator app not locked out: %w", err)
	}

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
		if err := s.updateUserAuthenticatorAppLockoutState(ctx, false); err != nil {
			return fmt.Errorf("update user authenticator app lockout state: %w", err)
		}

		return apierror.NewInvalidArgumentError("invalid totp code", err)
	}

	if err := s.updateUserAuthenticatorAppLockoutState(ctx, true); err != nil {
		return fmt.Errorf("update user authenticator app lockout state: %w", err)
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
	if err != nil {
		return fmt.Errorf("get organization by id: %w", err)
	}

	qMatchingUser, err := s.matchUser(ctx, q, qOrg, qIntermediateSession)
	if err != nil {
		return fmt.Errorf("match user: %w", err)
	}

	var ok bool
	var recoveryCodeBcrypts [][]byte
	for _, b := range qMatchingUser.AuthenticatorAppRecoveryCodeBcrypts {
		if bcrypt.CompareHashAndPassword(b, []byte(backupCode)) == nil {
			ok = true
			continue // do not keep this bcrypt around; the code is used
		}

		recoveryCodeBcrypts = append(recoveryCodeBcrypts, b)
	}

	// write back the remaining backup codes to the user
	if _, err := q.UpdateUserAuthenticatorAppRecoveryCodeBcrypts(ctx, queries.UpdateUserAuthenticatorAppRecoveryCodeBcryptsParams{
		ID:                                  qMatchingUser.ID,
		AuthenticatorAppRecoveryCodeBcrypts: recoveryCodeBcrypts,
	}); err != nil {
		return fmt.Errorf("update intermediate session authenticator app backup code bcrypts: %w", err)
	}

	// commit; our writes conflict with those from
	// updateUserAuthenticatorAppLockoutState
	if err := commit(); err != nil {
		return fmt.Errorf("commit: %w", err)
	}

	if !ok {
		if err := s.updateUserAuthenticatorAppLockoutState(ctx, false); err != nil {
			return fmt.Errorf("update user authenticator app lockout state: %w", err)
		}

		return apierror.NewInvalidArgumentError("invalid backup code", nil)
	}

	if err := s.updateUserAuthenticatorAppLockoutState(ctx, true); err != nil {
		return fmt.Errorf("update user authenticator app lockout state: %w", err)
	}

	return nil
}

func (s *Store) checkAuthenticatorAppLockedOut(ctx context.Context) error {
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

	qMatchingUser, err := s.matchUser(ctx, q, qOrg, qIntermediateSession)
	if err != nil {
		return fmt.Errorf("match user: %w", err)
	}

	if qMatchingUser.AuthenticatorAppLockoutExpireTime != nil && qMatchingUser.AuthenticatorAppLockoutExpireTime.After(time.Now()) {
		return apierror.NewFailedPreconditionError("too many authenticator app attempts; user is temporarily locked out", nil)
	}
	return nil
}

func (s *Store) updateUserAuthenticatorAppLockoutState(ctx context.Context, lastAttemptSuccessful bool) error {
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
	if err != nil {
		return fmt.Errorf("get organization by id: %w", err)
	}

	qMatchingUser, err := s.matchUser(ctx, q, qOrg, qIntermediateSession)
	if err != nil {
		return fmt.Errorf("match user: %w", err)
	}

	if lastAttemptSuccessful {
		if _, err := q.UpdateUserFailedAuthenticatorAppAttempts(ctx, queries.UpdateUserFailedAuthenticatorAppAttemptsParams{
			ID:                             qMatchingUser.ID,
			FailedAuthenticatorAppAttempts: 0,
		}); err != nil {
			return fmt.Errorf("update user failed authenticator app attempts: %w", err)
		}

		if err := commit(); err != nil {
			return fmt.Errorf("commit: %w", err)
		}

		return nil
	}

	attempts := qMatchingUser.FailedAuthenticatorAppAttempts + 1
	if attempts >= backupCodeLockoutAttempts {
		// lock the user out
		expireTime := time.Now().Add(backupCodeLockoutDuration)
		if _, err := q.UpdateUserAuthenticatorAppLockoutExpireTime(ctx, queries.UpdateUserAuthenticatorAppLockoutExpireTimeParams{
			ID:                                qMatchingUser.ID,
			AuthenticatorAppLockoutExpireTime: &expireTime,
		}); err != nil {
			return fmt.Errorf("update user authenticator app lockout expire time: %w", err)
		}

		// reset fail count
		if _, err := q.UpdateUserFailedAuthenticatorAppAttempts(ctx, queries.UpdateUserFailedAuthenticatorAppAttemptsParams{
			ID:                             qMatchingUser.ID,
			FailedAuthenticatorAppAttempts: 0,
		}); err != nil {
			return fmt.Errorf("update user failed authenticator app attempts: %w", err)
		}

		if err := commit(); err != nil {
			return fmt.Errorf("commit: %w", err)
		}

		return apierror.NewFailedPreconditionError("too many authenticator app attempts; user is temporarily locked out", nil)
	}

	// bump attempt count
	if _, err := q.UpdateUserFailedAuthenticatorAppAttempts(ctx, queries.UpdateUserFailedAuthenticatorAppAttemptsParams{
		ID:                             qMatchingUser.ID,
		FailedAuthenticatorAppAttempts: attempts,
	}); err != nil {
		return fmt.Errorf("update user failed authenticator app attempts: %w", err)
	}
	return nil
}

func (s *Store) getIntermediateSessionPendingAuthenticatorAppSecret(ctx context.Context) ([]byte, error) {
	_, q, _, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	qIntermediateSession, err := q.GetIntermediateSessionByID(ctx, authn.IntermediateSessionID(ctx))
	if err != nil {
		return nil, fmt.Errorf("get intermediate session by id: %w", err)
	}

	// close tx before calling kms
	if err := rollback(); err != nil {
		return nil, fmt.Errorf("rollback: %w", err)
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

	// close tx before calling kms
	if err := rollback(); err != nil {
		return nil, fmt.Errorf("rollback: %w", err)
	}

	decryptRes, err := s.kms.Decrypt(ctx, &kms.DecryptInput{
		KeyId:               &s.authenticatorAppSecretsKMSKeyID,
		EncryptionAlgorithm: types.EncryptionAlgorithmSpecRsaesOaepSha256,
		CiphertextBlob:      qMatchingUser.AuthenticatorAppSecretCiphertext,
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
