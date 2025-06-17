package storetestutil

import (
	"context"
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	backendv1 "github.com/tesseral-labs/tesseral/internal/backend/gen/tesseral/backend/v1"
	"github.com/tesseral-labs/tesseral/internal/store/idformat"
)

type Console struct {
	db  *pgxpool.Pool
	kms *KMS

	DogfoodProjectID *uuid.UUID
	DogfoodUserID    string
	ConsoleDomain    string
}

func NewConsole(db *pgxpool.Pool, kms *KMS) *Console {
	console := &Console{db: db, kms: kms}
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

	_, err := c.db.Exec(context.Background(), sql)
	if err != nil {
		log.Panicf("failed to seed database: %v", err)
	}
}

type Project struct {
	ProjectID string
	UserID    string
}

func (c *Console) NewProject(t *testing.T) Project {
	const rootDomain = "tesseral.example.app"

	projectID := uuid.New()
	formattedProjectID := idformat.Project.Format(projectID)
	projectVaultDomain := fmt.Sprintf("%s.%s", strings.ReplaceAll(formattedProjectID, "_", "-"), rootDomain)

	// Create the backing organization for the new project
	organizationID := uuid.New()
	_, err := c.db.Exec(t.Context(), `
INSERT INTO organizations (id, display_name, project_id, log_in_with_google, log_in_with_microsoft, log_in_with_email, log_in_with_password, log_in_with_saml, log_in_with_authenticator_app, log_in_with_passkey, scim_enabled)
  VALUES ($1::uuid, $2, $3::uuid, true, true, true, true, true, true, true, true);
`,
		organizationID.String(),
		fmt.Sprintf("%s Backing Organization", formattedProjectID),
		*c.DogfoodProjectID,
	)
	if err != nil {
		t.Fatalf("failed to create backing organization for test project: %v", err)
	}

	// Create the project with the new vault domain
	_, err = c.db.Exec(t.Context(), `
INSERT INTO projects (id, organization_id, display_name, log_in_with_google, log_in_with_microsoft, log_in_with_email, log_in_with_password, log_in_with_saml, log_in_with_authenticator_app, log_in_with_passkey, vault_domain, email_send_from_domain, redirect_uri, cookie_domain)
  VALUES ($1::uuid, $2::uuid, $3, true, true, true, true, true, true, true, $4, $4, $4, $4);
`,
		projectID.String(),
		organizationID.String(),
		formattedProjectID,
		projectVaultDomain,
	)
	if err != nil {
		t.Fatalf("failed to create test project: %v", err)
	}

	// Create the project UI settings
	_, err = c.db.Exec(t.Context(), `
INSERT INTO project_ui_settings (id, project_id)
  VALUES (gen_random_uuid(), $1::uuid);
`,
		projectID.String(),
	)
	if err != nil {
		t.Fatalf("failed to create project UI settings for test project: %v", err)
	}

	// Create the project trusted domains
	_, err = c.db.Exec(t.Context(), `
INSERT INTO project_trusted_domains (id, project_id, domain)
  VALUES
	(gen_random_uuid(), $1::uuid, $2);
`,
		projectID.String(),
		projectVaultDomain,
	)
	if err != nil {
		t.Fatalf("failed to create project trusted domains for test project: %v", err)
	}

	userID := uuid.New()
	formattedUserID := idformat.User.Format(userID)
	userEmail := fmt.Sprintf("%s@%s", formattedUserID, projectVaultDomain)

	_, err = c.db.Exec(t.Context(), `
INSERT INTO users (id, email, password_bcrypt, organization_id, is_owner)
  VALUES ($1::uuid, $2, crypt('password', gen_salt('bf', 14)), $3::uuid, true);
`,
		userID.String(),
		userEmail,
		organizationID.String(),
	)
	if err != nil {
		t.Fatalf("failed to create test user: %v", err)
	}

	return Project{
		ProjectID: formattedProjectID,
		UserID:    formattedUserID,
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
	projectID, err := idformat.Project.Parse(params.Project.ProjectID)
	if err != nil {
		t.Fatalf("failed to parse project ID: %v", err)
	}

	organizationID := uuid.New()
	formattedOrganizationID := idformat.Organization.Format(organizationID)

	// Create the organization
	_, err = c.db.Exec(t.Context(), `
INSERT INTO organizations (id, display_name, project_id, log_in_with_google, log_in_with_microsoft, log_in_with_email, log_in_with_password, log_in_with_saml, log_in_with_authenticator_app, log_in_with_passkey, scim_enabled)
  VALUES ($1::uuid, $2, $3::uuid, $4, $5, $6, $7, $8, $9, $10, $11);
`,
		organizationID.String(),
		params.Organization.DisplayName,
		uuid.UUID(projectID).String(),
		params.Organization.GetLogInWithGoogle(),
		params.Organization.GetLogInWithMicrosoft(),
		params.Organization.GetLogInWithEmail(),
		params.Organization.GetLogInWithPassword(),
		params.Organization.GetLogInWithSaml(),
		params.Organization.GetLogInWithAuthenticatorApp(),
		params.Organization.GetLogInWithPasskey(),
		params.Organization.GetScimEnabled(),
	)
	if err != nil {
		t.Fatalf("failed to create organization: %v", err)
	}

	// Create an owner for the organization
	userID := uuid.New()
	formattedUserID := idformat.User.Format(userID)
	userEmail := fmt.Sprintf("%s@%s", formattedUserID, params.Project.ProjectID)

	_, err = c.db.Exec(t.Context(), `
INSERT INTO users (id, email, password_bcrypt, organization_id, is_owner)
  VALUES ($1::uuid, $2, crypt('password', gen_salt('bf', 14)), $3::uuid, true);
`,
		userID.String(),
		userEmail,
		organizationID.String(),
	)
	if err != nil {
		t.Fatalf("failed to create user in organization: %v", err)
	}

	return Organization{
		ProjectID:      params.ProjectID,
		OrganizationID: formattedOrganizationID,
		UserID:         formattedUserID,
	}
}
