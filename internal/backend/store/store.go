package store

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/cloudflare/cloudflare-go/v4"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	stripeclient "github.com/stripe/stripe-go/v82/client"
	svix "github.com/svix/svix-webhooks/go"
	"github.com/tesseral-labs/tesseral/internal/backend/authn"
	"github.com/tesseral-labs/tesseral/internal/backend/store/queries"
	"github.com/tesseral-labs/tesseral/internal/cloudflaredoh"
	"github.com/tesseral-labs/tesseral/internal/common/apierror"
	common "github.com/tesseral-labs/tesseral/internal/common/store"
	"github.com/tesseral-labs/tesseral/internal/pagetoken"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Store struct {
	db                                    *pgxpool.Pool
	dogfoodProjectID                      *uuid.UUID
	consoleDomain                         string
	intermediateSessionSigningKeyKMSKeyID string
	kms                                   *kms.Client
	ses                                   *sesv2.Client
	cloudflare                            *cloudflare.Client
	cloudflareDOH                         *cloudflaredoh.Client
	pageEncoder                           pagetoken.Encoder
	q                                     *queries.Queries
	s3                                    *s3.Client
	s3PresignClient                       *s3.PresignClient
	s3UserContentBucketName               string
	sessionSigningKeyKmsKeyID             string
	googleOAuthClientSecretsKMSKeyID      string
	microsoftOAuthClientSecretsKMSKeyID   string
	githubOAuthClientSecretsKMSKeyID      string
	userContentBaseUrl                    string
	authAppsRootDomain                    string
	tesseralDNSCloudflareZoneID           string
	tesseralDNSVaultCNAMEValue            string
	sesSPFMXRecordValue                   string
	stripe                                *stripeclient.API
	stripePriceIDGrowthTier               string
	svixClient                            *svix.Svix
	common                                *common.Store
}

type NewStoreParams struct {
	DB                                    *pgxpool.Pool
	DogfoodProjectID                      *uuid.UUID
	ConsoleDomain                         string
	IntermediateSessionSigningKeyKMSKeyID string
	KMS                                   *kms.Client
	SES                                   *sesv2.Client
	Cloudflare                            *cloudflare.Client
	CloudflareDOH                         *cloudflaredoh.Client
	PageEncoder                           pagetoken.Encoder
	S3                                    *s3.Client
	S3UserContentBucketName               string
	SessionSigningKeyKmsKeyID             string
	GoogleOAuthClientSecretsKMSKeyID      string
	MicrosoftOAuthClientSecretsKMSKeyID   string
	GithubOAuthClientSecretsKMSKeyID      string
	UserContentBaseUrl                    string
	AuthAppsRootDomain                    string
	TesseralDNSCloudflareZoneID           string
	TesseralDNSVaultCNAMEValue            string
	SESSPFMXRecordValue                   string
	Stripe                                *stripeclient.API
	StripePriceIDGrowthTier               string
	SvixClient                            *svix.Svix
	CommonStore                           *common.Store
}

func New(p NewStoreParams) *Store {
	store := &Store{
		db:                                    p.DB,
		dogfoodProjectID:                      p.DogfoodProjectID,
		consoleDomain:                         p.ConsoleDomain,
		intermediateSessionSigningKeyKMSKeyID: p.IntermediateSessionSigningKeyKMSKeyID,
		kms:                                   p.KMS,
		ses:                                   p.SES,
		cloudflare:                            p.Cloudflare,
		cloudflareDOH:                         p.CloudflareDOH,
		pageEncoder:                           p.PageEncoder,
		q:                                     queries.New(p.DB),
		s3:                                    p.S3,
		s3PresignClient:                       s3.NewPresignClient(p.S3),
		s3UserContentBucketName:               p.S3UserContentBucketName,
		sessionSigningKeyKmsKeyID:             p.SessionSigningKeyKmsKeyID,
		googleOAuthClientSecretsKMSKeyID:      p.GoogleOAuthClientSecretsKMSKeyID,
		microsoftOAuthClientSecretsKMSKeyID:   p.MicrosoftOAuthClientSecretsKMSKeyID,
		githubOAuthClientSecretsKMSKeyID:      p.GithubOAuthClientSecretsKMSKeyID,
		userContentBaseUrl:                    p.UserContentBaseUrl,
		authAppsRootDomain:                    p.AuthAppsRootDomain,
		tesseralDNSCloudflareZoneID:           p.TesseralDNSCloudflareZoneID,
		tesseralDNSVaultCNAMEValue:            p.TesseralDNSVaultCNAMEValue,
		sesSPFMXRecordValue:                   p.SESSPFMXRecordValue,
		stripe:                                p.Stripe,
		stripePriceIDGrowthTier:               p.StripePriceIDGrowthTier,
		svixClient:                            p.SvixClient,
		common:                                p.CommonStore,
	}

	return store
}

func (s *Store) tx(ctx context.Context) (tx pgx.Tx, q *queries.Queries, commit func() error, rollback func() error, err error) {
	tx, err = s.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("begin tx: %w", err)
	}

	commit = func() error { return tx.Commit(ctx) }
	rollback = func() error { return tx.Rollback(ctx) }
	return tx, queries.New(tx), commit, rollback, nil
}

func derefOrEmpty[T any](t *T) T {
	var z T
	if t == nil {
		return z
	}
	return *t
}

func refOrNil[T comparable](t T) *T {
	var z T
	if t == z {
		return nil
	}
	return &t
}

func timestampOrNil(t *time.Time) *timestamppb.Timestamp {
	if t == nil || t.IsZero() {
		return nil
	}
	return timestamppb.New(*t)
}

// validateIsDogfoodSession returns an error if the caller isn't a dogfood
// session.
//
// The intention of this method is to allow endpoints to prevent themselves from
// being called by project API keys.
func validateIsDogfoodSession(ctx context.Context) error {
	data := authn.GetContextData(ctx)
	if data.DogfoodSession == nil {
		return apierror.NewUnauthenticatedError("this endpoint cannot be invoked by project API keys", fmt.Errorf("non-dogfood session request"))
	}
	return nil
}
