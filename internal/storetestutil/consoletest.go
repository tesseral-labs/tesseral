package storetestutil

import (
	"context"
	"fmt"
	"log"
	"math/rand/v2"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	bkauthn "github.com/tesseral-labs/tesseral/internal/backend/authn"
	backendv1 "github.com/tesseral-labs/tesseral/internal/backend/gen/tesseral/backend/v1"
	bkstore "github.com/tesseral-labs/tesseral/internal/backend/store"
	intauthn "github.com/tesseral-labs/tesseral/internal/intermediate/authn"
	intermediatev1 "github.com/tesseral-labs/tesseral/internal/intermediate/gen/tesseral/intermediate/v1"
	intstore "github.com/tesseral-labs/tesseral/internal/intermediate/store"
	"github.com/tesseral-labs/tesseral/internal/store/idformat"
)

type Console struct {
	pool *pgxpool.Pool
	KMS  *KMS

	DogfoodProjectID *uuid.UUID
	DogfoodUserID    string
	ConsoleDomain    string
}

func NewConsole(pool *pgxpool.Pool, kms *KMS) *Console {
	console := &Console{pool: pool, KMS: kms}
	console.seed()
	return console
}

func (c *Console) seed() {
	var (
		dogfoodProjectID = uuid.MustParse("252491cc-76e3-4957-ab23-47d83c34f240")
		dogfoodUserID    = uuid.MustParse("e071bbfe-6f27-4526-ab37-0ad251742836")
	)

	c.DogfoodProjectID = &dogfoodProjectID
	c.DogfoodUserID = idformat.User.Format(dogfoodUserID)
	c.ConsoleDomain = "console.tesseral.example.com"

	const sql = `
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

INSERT INTO projects (id, display_name, log_in_with_google, log_in_with_microsoft, log_in_with_email, log_in_with_password, log_in_with_saml, log_in_with_authenticator_app, log_in_with_passkey, vault_domain, email_send_from_domain, redirect_uri, cookie_domain)
	VALUES ('252491cc-76e3-4957-ab23-47d83c34f240'::uuid, 'Tesseral Test', true, true, true, true, true, true, true, 'vault.console.tesseral.example.com', 'vault.console.tesseral.example.com', 'https://console.tesseral.example.com', 'console.tesseral.example.com');

INSERT INTO project_trusted_domains (id, project_id, domain)
VALUES
    (gen_random_uuid(), '252491cc-76e3-4957-ab23-47d83c34f240', 'vault.console.tesseral.example.com'),
    (gen_random_uuid(), '252491cc-76e3-4957-ab23-47d83c34f240', 'console.tesseral.example.com');

-- Create the Dogfood Project's backing organization
INSERT INTO organizations (id, display_name, project_id, log_in_with_google, log_in_with_microsoft, log_in_with_email, log_in_with_password, log_in_with_saml, log_in_with_authenticator_app, log_in_with_passkey, scim_enabled)
  VALUES ('7a76decb-6d79-49ce-9449-34fcc53151df'::uuid, 'project_54vwf0clhh0caqe20eujxgpeq Backing Organization', '252491cc-76e3-4957-ab23-47d83c34f240', true, false, true, true, true, true, true, true);

UPDATE projects SET organization_id = '7a76decb-6d79-49ce-9449-34fcc53151df'::uuid where id = '252491cc-76e3-4957-ab23-47d83c34f240'::uuid;

-- Create project UI settings for the dogfood project
INSERT INTO project_ui_settings (id, project_id)
  VALUES (gen_random_uuid(), '252491cc-76e3-4957-ab23-47d83c34f240'::uuid);

-- Create a user in the dogfood project
INSERT INTO users (id, email, password_bcrypt, organization_id, is_owner)
  VALUES ('e071bbfe-6f27-4526-ab37-0ad251742836'::uuid, 'root@app.tesseral.example.com', crypt('password', gen_salt('bf', 14)), '7a76decb-6d79-49ce-9449-34fcc53151df', true);
`

	_, err := c.pool.Exec(context.Background(), sql)
	if err != nil {
		log.Panicf("failed to seed database: %v", err)
	}
}

type Project struct {
	ProjectID string
	UserID    string
}

func (c *Console) NewProject(t *testing.T) Project {
	projectName := fmt.Sprintf("test-%d", rand.IntN(1<<20))

	intermediateSession := &intermediatev1.IntermediateSession{
		Id:                idformat.IntermediateSession.Format(uuid.New()),
		ProjectId:         idformat.Project.Format(*c.DogfoodProjectID),
		Email:             "root@app.tesseral.example.com",
		EmailVerified:     true,
		OrganizationId:    idformat.Organization.Format(uuid.MustParse("7a76decb-6d79-49ce-9449-34fcc53151df")),
		PrimaryAuthFactor: intermediatev1.PrimaryAuthFactor_PRIMARY_AUTH_FACTOR_EMAIL,
	}
	intstore := intstore.New(intstore.NewStoreParams{
		DB:                        c.pool,
		KMS:                       c.KMS.Client,
		SessionSigningKeyKmsKeyID: c.KMS.SessionSigningKeyID,
		S3:                        s3.NewFromConfig(*aws.NewConfig()),
		DogfoodProjectID:          c.DogfoodProjectID,
	})
	ctx := intauthn.NewContext(t.Context(), intermediateSession, idformat.Project.Format(*c.DogfoodProjectID))
	project, err := intstore.CreateProject(ctx, &intermediatev1.CreateProjectRequest{
		DisplayName: projectName,
		RedirectUri: fmt.Sprintf("https://%s.tesseral.example.com", projectName),
	})
	if err != nil {
		t.Fatalf("failed to create test project: %v", err)
	}

	_, err = c.pool.Exec(t.Context(), `
UPDATE projects SET
	log_in_with_google = true,
	log_in_with_microsoft = true,
	log_in_with_github = true,
	log_in_with_email = true,
	log_in_with_password = true,
	log_in_with_saml = true,
	log_in_with_authenticator_app = true,
	log_in_with_passkey = true
WHERE id = $1::uuid;
`,
		project.Project.Id,
	)
	if err != nil {
		t.Fatalf("failed to update test project: %v", err)
	}

	userID := uuid.New()
	userEmail := fmt.Sprintf("user-%d@%s.tesseral.example.com", rand.IntN(1<<20), projectName)

	_, err = c.pool.Exec(t.Context(), `
INSERT INTO users (id, email, password_bcrypt, organization_id, is_owner)
  VALUES ($1::uuid, $2, crypt('password', gen_salt('bf', 14)), (SELECT organization_id FROM projects WHERE id=$3::uuid), true);
`,
		userID.String(),
		userEmail,
		project.Project.Id,
	)
	if err != nil {
		t.Fatalf("failed to create test user: %v", err)
	}

	return Project{
		ProjectID: idformat.Project.Format(uuid.MustParse(project.Project.Id)),
		UserID:    idformat.User.Format(userID),
	}
}

type OrganizationParams struct {
	Project
	*backendv1.Organization
}

type Organization struct {
	ProjectID      string
	OrganizationID string
	UserID         string
}

func (c *Console) NewOrganization(t *testing.T, params OrganizationParams) Organization {
	bkstore := bkstore.New(bkstore.NewStoreParams{
		DB:                        c.pool,
		KMS:                       c.KMS.Client,
		SessionSigningKeyKmsKeyID: c.KMS.SessionSigningKeyID,
		S3:                        s3.NewFromConfig(*aws.NewConfig()),
		DogfoodProjectID:          c.DogfoodProjectID,
	})
	ctx := bkauthn.NewDogfoodSessionContext(t.Context(), bkauthn.DogfoodSessionContextData{
		ProjectID: params.ProjectID,
		UserID:    params.UserID,
		SessionID: idformat.Session.Format(uuid.New()),
	})

	organization, err := bkstore.CreateOrganization(ctx, &backendv1.CreateOrganizationRequest{
		Organization: params.Organization,
	})
	if err != nil {
		t.Fatalf("failed to create test organization: %v", err)
	}

	user, err := bkstore.CreateUser(ctx, &backendv1.CreateUserRequest{
		User: &backendv1.User{
			OrganizationId: organization.Organization.Id,
			Email:          fmt.Sprintf("user-%d@example.com", rand.IntN(1<<20)),
		},
	})
	if err != nil {
		t.Fatalf("failed to create user in organization: %v", err)
	}

	return Organization{
		ProjectID:      params.ProjectID,
		OrganizationID: organization.Organization.Id,
		UserID:         user.User.Id,
	}
}
