package store

import (
	"bytes"
	"context"
	"crypto/rand"
	"errors"
	"math/big"
	"time"

	"github.com/google/uuid"
	"github.com/openauth-dev/openauth/internal/crypto"
	"github.com/openauth-dev/openauth/internal/store/idformat"
	"github.com/openauth-dev/openauth/internal/store/queries"
)

type MethodVerificationChallenge struct {
	ID 									uuid.UUID
	ProjectID 					uuid.UUID
	Email 							string
	AuthMethod 					queries.AuthMethod
	ExpireTime 					time.Time
	CompleteTime 				time.Time
	SecretTokenSha256 	[]byte
}

var ErrMethodVerificationChallengeExpired = errors.New("method verification challenge has expired")
var ErrMethodVerificationChallengeSecretTokenMismatch = errors.New("method verification challenge secret token mismatch")

type CreateMethodVerificationChallengeRequest struct {
	ProjectID 				string
	UnverifiedEmail 	string
	AuthMethod 				queries.AuthMethod
}

func (s *Store) CreateMethodVerificationChallenge(
	ctx *context.Context, 
	req *CreateMethodVerificationChallengeRequest,
) (*MethodVerificationChallenge, error) {
	_, q, commit, rollback, err := s.tx(*ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	projectId, err := idformat.Project.Parse(req.ProjectID)
	if err != nil {
		return nil, err
	}

	// Allow 15 minutes for the user to verify their email before expiring the intermediate session
	expiresAt := time.Now().Add(time.Minute * 15)

	// Create a new secret token for the challenge
	secretToken, err := generateSecretToken()
	if err != nil {
		return nil, err
	}
	secretTokenSha256 := crypto.StringToSha256(secretToken)

	createdMethodVerificationChallenge, err := q.CreateMethodVerificationChallenge(*ctx, queries.CreateMethodVerificationChallengeParams{
		ID: uuid.New(),
		ProjectID: projectId,
		Email: req.UnverifiedEmail,
		AuthMethod: req.AuthMethod,
		ExpireTime: &expiresAt,
		SecretTokenSha256: secretTokenSha256,
	})
	if err != nil {
		return nil, err
	}

	if err := commit(); err != nil {
		return nil, err
	}

	return transformMethodVerificationChallenge(createdMethodVerificationChallenge), nil
}

type CompleteMethodVerificationChallengeRequest struct {
	ID 					string
	SecretToken string
}

func (s *Store) CompleteMethodVerificationChallenge(
	ctx *context.Context, 
	req *CompleteMethodVerificationChallengeRequest,
) error {
	_, q, commit, rollback, err := s.tx(*ctx)
	if err != nil {
		return nil
	}
	defer rollback()

	mvcID, err := idformat.MethodVerificationChallenge.Parse(req.ID)
	if err != nil {
		return err
	}

	existingMVC, err := q.GetMethodVerificationChallengeByID(*ctx, mvcID)
	if err != nil {
		return err
	}

	// Check if the challenge has expired
	if existingMVC.ExpireTime.Before(time.Now()) {
		return ErrMethodVerificationChallengeExpired
	}

	secretTokenSha256 := crypto.StringToSha256(req.SecretToken)

	// Check if the secret token matches
	if !bytes.Equal(existingMVC.SecretTokenSha256, secretTokenSha256) {
		return ErrMethodVerificationChallengeSecretTokenMismatch
	}

	// Update the challenge to mark it as complete
	completeTime := time.Now()
	_, err = q.CompleteMethodVerificationChallenge(*ctx, queries.CompleteMethodVerificationChallengeParams{
		ID: mvcID,
		CompleteTime: &completeTime,
	})
	if err != nil {
		return err
	}

	if err := commit(); err != nil {
		return err
	}

	return nil
}

func generateSecretToken() (string, error) {
	// Define the range for a 6-digit number: [100000, 999999]
	min := 100000
	max := 999999
	rangeSize := max - min + 1

	// Generate a secure random number
	randomBigInt, err := rand.Int(rand.Reader, big.NewInt(int64(rangeSize)))
	if err != nil {
		return "", err
	}

	// Convert to the desired range
	randomNumber := min + int(randomBigInt.Int64())

	return string(randomNumber), nil
}

func transformMethodVerificationChallenge(mvc queries.MethodVerificationChallenge) *MethodVerificationChallenge {
	return &MethodVerificationChallenge{
		ID: mvc.ID,
		ProjectID: mvc.ProjectID,
		Email: mvc.Email,
		AuthMethod: mvc.AuthMethod,
		ExpireTime: *mvc.ExpireTime,
		CompleteTime: *mvc.CompleteTime,
		SecretTokenSha256: mvc.SecretTokenSha256,
	}
}