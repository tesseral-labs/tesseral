package consoletest

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tesseral-labs/tesseral/internal/store/idformat"
)

type Console struct {
	pool             *pgxpool.Pool
	DogfoodProjectID *uuid.UUID
}

func New(t *testing.T, pool *pgxpool.Pool) *Console {
	console := &Console{pool: pool}
	console.DogfoodProjectID = console.seed(t)
	return console
}

func (c *Console) seed(t *testing.T) *uuid.UUID {
	dogfoodProjectID := uuid.New()
	c.DogfoodProjectID = &dogfoodProjectID

	sql := fmt.Sprintf(`
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

INSERT INTO projects (id, display_name, log_in_with_google, log_in_with_microsoft, log_in_with_email, log_in_with_password, log_in_with_saml, log_in_with_authenticator_app, log_in_with_passkey, vault_domain, email_send_from_domain, redirect_uri, cookie_domain)
	VALUES ('%s'::uuid, 'Tesseral Test', true, true, true, true, true, true, true, 'vault.console.tesseral.example.com', 'vault.console.tesseral.example.com', 'https://console.tesseral.example.com', 'console.tesseral.example.com');

INSERT INTO project_trusted_domains (id, project_id, domain)
VALUES
    (gen_random_uuid(), '%s', 'vault.console.tesseral.example.com'),
    (gen_random_uuid(), '%s', 'console.tesseral.example.com');

-- Create the Dogfood Project's backing organization
INSERT INTO organizations (id, display_name, project_id, log_in_with_google, log_in_with_microsoft, log_in_with_email, log_in_with_password, log_in_with_saml, log_in_with_authenticator_app, log_in_with_passkey, scim_enabled)
  VALUES ('7a76decb-6d79-49ce-9449-34fcc53151df'::uuid, 'project_54vwf0clhh0caqe20eujxgpeq Backing Organization', '%s', true, false, true, true, true, true, true, true);

UPDATE projects SET organization_id = '7a76decb-6d79-49ce-9449-34fcc53151df'::uuid where id = '%s'::uuid;

-- Create project UI settings for the dogfood project
INSERT INTO project_ui_settings (id, project_id)
  VALUES (gen_random_uuid(), '%s'::uuid);
	`, dogfoodProjectID.String(), dogfoodProjectID.String(), dogfoodProjectID.String(), dogfoodProjectID.String(), dogfoodProjectID.String(), dogfoodProjectID.String())

	_, err := c.pool.Exec(context.Background(), sql)
	if err != nil {
		t.Fatalf("failed to seed database: %v", err)
	}

	return &dogfoodProjectID
}

type ProjectParams struct {
	Name          string
	LoginWithSaml bool
}

type Project struct {
	ProjectID     string
	BackendAPIKey string
}

func (c *Console) CreateProject(t *testing.T, params ProjectParams) Project {
	projectID := uuid.New()
	projectBackingOrganizationID := uuid.New()
	backendAPIKey := uuid.New()

	sql := fmt.Sprintf(`
INSERT INTO projects (id, display_name, log_in_with_google, log_in_with_microsoft, log_in_with_email, log_in_with_password, log_in_with_saml, log_in_with_authenticator_app, log_in_with_passkey, vault_domain, email_send_from_domain, redirect_uri, cookie_domain)
	VALUES ('%s'::uuid, 'Test Project', true, true, true, true, %t, true, true, 'vault.%s.tesseral.example.com', 'vault.%s.tesseral.example.com', 'https://%s.tesseral.example.com', '%s.tesseral.example.com');

INSERT INTO backend_api_keys (id, project_id, secret_token_sha256, display_name)
  VALUES (gen_random_uuid(), '%s', digest(uuid_send('%s'::uuid), 'sha256'), 'test');

-- Backing organization
INSERT INTO organizations (id, display_name, project_id, log_in_with_saml, scim_enabled, log_in_with_email)
VALUES ('%s'::uuid, '%s Backing Organization', '%s', false, false, true);

update projects set organization_id = '%s'::uuid where id = '%s'::uuid;
	`, projectID.String(), params.LoginWithSaml, params.Name, params.Name, params.Name, params.Name, projectID.String(), backendAPIKey.String(),
		projectBackingOrganizationID.String(), params.Name, projectID.String(), projectBackingOrganizationID.String(), projectID.String())

	_, err := c.pool.Exec(context.Background(), sql)
	if err != nil {
		t.Fatalf("failed to create test project: %v", err)
	}

	return Project{
		ProjectID:     idformat.Project.Format(projectID),
		BackendAPIKey: idformat.BackendAPIKey.Format(backendAPIKey),
	}
}

type OrganizationParams struct {
	ProjectID     string
	Name          string
	LoginWithSaml bool
}

type Organization struct {
	ProjectID      string
	OrganizationID string
	UserID         string
}

func (c *Console) CreateOrganization(t *testing.T, params OrganizationParams) Organization {
	projectUUID, err := idformat.Project.Parse(params.ProjectID)
	if err != nil {
		t.Fatalf("parse project ID: %v", err)
	}

	organizationID := uuid.New()
	userID := uuid.New()
	userEmail := fmt.Sprintf("%s@%s.example.com", userID.String(), params.Name)

	_, err = c.pool.Exec(context.Background(), `
INSERT INTO organizations (id, display_name, project_id, log_in_with_google, log_in_with_microsoft, log_in_with_email, log_in_with_password, log_in_with_saml, log_in_with_authenticator_app, log_in_with_passkey, scim_enabled)
  VALUES ($1::uuid, $2, $3::uuid, true, false, true, true, $4, true, true, true);
`,
		organizationID.String(),
		params.Name,
		uuid.UUID(projectUUID).String(),
		params.LoginWithSaml,
	)
	if err != nil {
		t.Fatalf("failed to create test organization: %v", err)
	}

	_, err = c.pool.Exec(context.Background(), `
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

	return Organization{
		ProjectID:      params.ProjectID,
		OrganizationID: idformat.Organization.Format(organizationID),
		UserID:         idformat.User.Format(userID),
	}
}
