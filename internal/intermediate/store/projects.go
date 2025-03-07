package store

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go-v2/service/kms/types"
	"github.com/google/uuid"
	"github.com/tesseral-labs/tesseral/internal/common/apierror"
	"github.com/tesseral-labs/tesseral/internal/intermediate/authn"
	intermediatev1 "github.com/tesseral-labs/tesseral/internal/intermediate/gen/tesseral/intermediate/v1"
	"github.com/tesseral-labs/tesseral/internal/intermediate/store/queries"
	"github.com/tesseral-labs/tesseral/internal/store/idformat"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Store) CreateProject(ctx context.Context, req *intermediatev1.CreateProjectRequest) (*intermediatev1.CreateProjectResponse, error) {
	if authn.ProjectID(ctx) != *s.dogfoodProjectID {
		return nil, apierror.NewPermissionDeniedError("cannot create a project", fmt.Errorf("create project attempted on non-dogfood project: %s", authn.ProjectID(ctx)))
	}

	_, q, commit, rollback, err := s.tx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	intermediateSession := authn.IntermediateSession(ctx)
	if !intermediateSession.EmailVerified {
		return nil, apierror.NewPermissionDeniedError("email not verified", nil)
	}

	// create this ahead of time so we can use it in the display name and auth domain
	newProjectID := uuid.New()
	formattedNewProjectID := idformat.Project.Format(newProjectID)
	newProjectVaultDomain := fmt.Sprintf("%s.%s", strings.ReplaceAll(formattedNewProjectID, "_", "-"), s.authAppsRootDomain)

	// create a new organization under the dogfood project, accepting the same
	// primary login method used to get to this point
	qOrganization, err := q.CreateOrganization(ctx, queries.CreateOrganizationParams{
		ID:                 uuid.New(),
		DisplayName:        fmt.Sprintf("%s Backing Organization", formattedNewProjectID),
		ProjectID:          *s.dogfoodProjectID,
		LogInWithEmail:     intermediateSession.PrimaryLoginFactor == string(queries.PrimaryLoginFactorEmail),
		LogInWithGoogle:    intermediateSession.PrimaryLoginFactor == string(queries.PrimaryLoginFactorGoogleOauth),
		LogInWithMicrosoft: intermediateSession.PrimaryLoginFactor == string(queries.PrimaryLoginFactorMicrosoftOauth),
	})
	if err != nil {
		return nil, fmt.Errorf("create organization: %w", err)
	}

	// reflect the google hosted domain from the intermediate session if it exists
	if intermediateSession.GoogleHostedDomain != "" {
		if _, err := q.CreateOrganizationGoogleHostedDomain(ctx, queries.CreateOrganizationGoogleHostedDomainParams{
			ID:                 uuid.New(),
			OrganizationID:     qOrganization.ID,
			GoogleHostedDomain: intermediateSession.GoogleHostedDomain,
		}); err != nil {
			return nil, fmt.Errorf("create organization google hosted domain: %w", err)
		}
	}

	// reflect the microsoft tenant id from the intermediate session if it exists
	if intermediateSession.MicrosoftTenantId != "" {
		if _, err := q.CreateOrganizationMicrosoftTenantID(ctx, queries.CreateOrganizationMicrosoftTenantIDParams{
			ID:                uuid.New(),
			OrganizationID:    qOrganization.ID,
			MicrosoftTenantID: intermediateSession.MicrosoftTenantId,
		}); err != nil {
			return nil, fmt.Errorf("create organization microsoft tenant id: %w", err)
		}
	}

	// create a new user invite for the intermediate session user
	if _, err := q.CreateUserInvite(ctx, queries.CreateUserInviteParams{
		ID:             uuid.New(),
		OrganizationID: qOrganization.ID,
		Email:          intermediateSession.Email,
		IsOwner:        true,
	}); err != nil {
		return nil, fmt.Errorf("create user invite: %w", err)
	}

	// create a new project backed by the new organization; the login methods
	// here are only those that can work out of the box, without further
	// configuration by the user
	qProject, err := q.CreateProject(ctx, queries.CreateProjectParams{
		ID:                  newProjectID,
		RedirectUri:         req.RedirectUri,
		OrganizationID:      &qOrganization.ID,
		VaultDomain:         newProjectVaultDomain,
		EmailSendFromDomain: fmt.Sprintf("mail.%s", s.authAppsRootDomain),
		DisplayName:         req.DisplayName,
		LogInWithEmail:      true,
		LogInWithGoogle:     false,
		LogInWithMicrosoft:  false,
		LogInWithPassword:   false,
		LogInWithSaml:       false,
	})
	if err != nil {
		return nil, fmt.Errorf("create project: %w", err)
	}

	if _, err := q.CreateProjectUISettings(ctx, qProject.ID); err != nil {
		return nil, fmt.Errorf("create project ui settings: %w", err)
	}

	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, err
	}

	privateKeyBytes, err := x509.MarshalECPrivateKey(privateKey)
	if err != nil {
		return nil, err
	}

	// Encrypt the symmetric key with the KMS
	sskEncryptOutput, err := s.kms.Encrypt(ctx, &kms.EncryptInput{
		EncryptionAlgorithm: types.EncryptionAlgorithmSpecRsaesOaepSha256,
		KeyId:               &s.sessionSigningKeyKmsKeyID,
		Plaintext:           privateKeyBytes,
	})
	if err != nil {
		return nil, fmt.Errorf("encrypt session signing key: %w", err)
	}

	publicKeyBytes, err := x509.MarshalPKIXPublicKey(privateKey.Public())
	if err != nil {
		return nil, fmt.Errorf("marshal public key: %w", err)
	}

	createTime := time.Now()
	expireTime := createTime.Add(time.Hour * 24 * 7)

	if _, err := q.CreateSessionSigningKey(ctx, queries.CreateSessionSigningKeyParams{
		ID:                   uuid.New(),
		ProjectID:            qProject.ID,
		PublicKey:            publicKeyBytes,
		PrivateKeyCipherText: sskEncryptOutput.CiphertextBlob,
		CreateTime:           &createTime,
		ExpireTime:           &expireTime,
	}); err != nil {
		return nil, fmt.Errorf("create session signing key: %w", err)
	}

	if err = commit(); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}

	return &intermediatev1.CreateProjectResponse{
		Project: parseProject(qProject),
	}, nil
}

func parseProject(qProject queries.Project) *intermediatev1.Project {
	return &intermediatev1.Project{
		Id:             qProject.ID.String(),
		OrganizationId: idformat.Organization.Format(*qProject.OrganizationID),
		CreateTime:     timestamppb.New(*qProject.CreateTime),
		UpdateTime:     timestamppb.New(*qProject.UpdateTime),
		DisplayName:    qProject.DisplayName,
		VaultDomain:    qProject.VaultDomain,
	}
}
