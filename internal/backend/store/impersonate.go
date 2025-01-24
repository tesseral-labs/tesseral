package store

import (
	"context"
	"crypto/sha256"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/openauth/openauth/internal/backend/authn"
	backendv1 "github.com/openauth/openauth/internal/backend/gen/openauth/backend/v1"
	"github.com/openauth/openauth/internal/backend/store/queries"
	"github.com/openauth/openauth/internal/common/apierror"
	"github.com/openauth/openauth/internal/store/idformat"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const userImpersonationTokenDuration = time.Second * 30

func (s *Store) CreateUserImpersonationToken(ctx context.Context, req *backendv1.CreateUserImpersonationTokenRequest) (*backendv1.CreateUserImpersonationTokenResponse, error) {
	if err := validateIsDogfoodSession(ctx); err != nil {
		return nil, fmt.Errorf("validate is dogfood session: %w", err)
	}

	impersonatorID, err := idformat.User.Parse(authn.GetContextData(ctx).DogfoodSession.UserID)
	if err != nil {
		panic(fmt.Errorf("parse user id: %w", err))
	}

	_, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	qImpersonator, err := q.GetUser(ctx, queries.GetUserParams{
		ProjectID: *s.dogfoodProjectID,
		ID:        impersonatorID,
	})
	if err != nil {
		return nil, fmt.Errorf("get impersonator: %w", err)
	}

	if !qImpersonator.IsOwner {
		return nil, apierror.NewPermissionDeniedError("only owners may impersonate others", fmt.Errorf("impersonator is not an owner"))
	}

	impersonatedID, err := idformat.User.Parse(req.UserImpersonationToken.ImpersonatedId)
	if err != nil {
		return nil, fmt.Errorf("parse user id: %w", err)
	}

	// Ensure the impersonator can target the target impersonated user.
	//
	// You can only impersonate someone if the organization you belong to is the
	// backing organization for the project the target belongs to.
	//
	// One way to think about this: when you're logged into app.tesseral.com,
	// the only users you can impersonate are the ones you can see when
	// navigating around your project's organizations and users.
	//
	// That set of displayed users all belong to your project. When you're
	// logged into app.tesseral.com, you're logged into the dogfood project from
	// the backend's perspective. But the projects we let you list are (!!!)
	//
	// TODO complete this logic
	if _, err := q.GetUserForImpersonation(ctx, queries.GetUserForImpersonationParams{
		ID:                         impersonatedID,
		ImpersonatorOrganizationID: &qImpersonator.OrganizationID,
	}); err != nil {
		return nil, fmt.Errorf("get impersonated user: %w", err)
	}

	secretToken := uuid.New()
	secretTokenSHA256 := sha256.Sum256(secretToken[:])

	expireTime := time.Now().Add(userImpersonationTokenDuration)
	qUserImpersonationToken, err := q.CreateUserImpersonationToken(ctx, queries.CreateUserImpersonationTokenParams{
		ID:                uuid.New(),
		ImpersonatorID:    impersonatorID,
		ImpersonatedID:    impersonatedID,
		ExpireTime:        &expireTime,
		SecretTokenSha256: secretTokenSHA256[:],
	})
	if err != nil {
		return nil, fmt.Errorf("create user impersonation token: %w", err)
	}

	if err := commit(); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}

	userImpersonationToken := parseUserImpersonationToken(qUserImpersonationToken)
	userImpersonationToken.SecretToken = idformat.UserImpersonationSecretToken.Format(secretToken)
	return &backendv1.CreateUserImpersonationTokenResponse{
		UserImpersonationToken: userImpersonationToken,
	}, nil
}

func parseUserImpersonationToken(qUserImpersonationToken queries.UserImpersonationToken) *backendv1.UserImpersonationToken {
	var redeemTime *timestamppb.Timestamp
	if qUserImpersonationToken.RedeemTime != nil {
		redeemTime = timestamppb.New(*qUserImpersonationToken.RedeemTime)
	}

	return &backendv1.UserImpersonationToken{
		Id:             idformat.UserImpersonationToken.Format(qUserImpersonationToken.ID),
		ImpersonatorId: idformat.User.Format(qUserImpersonationToken.ImpersonatorID),
		CreateTime:     timestamppb.New(*qUserImpersonationToken.CreateTime),
		RedeemTime:     redeemTime,
		SecretToken:    "", // intentionally left blank
		ImpersonatedId: idformat.User.Format(qUserImpersonationToken.ImpersonatedID),
	}
}
