package storetesting

import (
	"context"
	"crypto/sha256"
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	backendv1 "github.com/tesseral-labs/tesseral/internal/backend/gen/tesseral/backend/v1"
	"github.com/tesseral-labs/tesseral/internal/store/idformat"
)

type Environment struct {
	DB  *pgxpool.Pool
	KMS *testKms
	S3  *testS3

	DogfoodProjectID   *uuid.UUID
	DogfoodUserID      string
	ConsoleDomain      string
	AuthAppsRootDomain string
}

// NewEnvironment initializes and returns test dependencies for testing store layers.
func NewEnvironment() (*Environment, func()) {
	db, cleanupDB := newDB()
	kms, cleanupKms := newKMS()
	s3, cleanupS3 := newS3()

	env := &Environment{
		DB:  db,
		S3:  s3,
		KMS: kms,
	}
	env.seed()

	return env, func() {
		cleanupS3()
		cleanupKms()
		cleanupDB()
	}
}

func (e *Environment) seed() {
	var (
		dogfoodProjectID = uuid.MustParse("252491cc-76e3-4957-ab23-47d83c34f240")
		dogfoodUserID    = uuid.MustParse("e071bbfe-6f27-4526-ab37-0ad251742836")
	)

	e.DogfoodProjectID = &dogfoodProjectID
	e.DogfoodUserID = idformat.User.Format(dogfoodUserID)
	e.ConsoleDomain = "console.tesseral.example.com"
	e.AuthAppsRootDomain = "tesseral.example.app"

	const sql = `
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

INSERT INTO projects (id, display_name, log_in_with_google, log_in_with_microsoft, log_in_with_email, log_in_with_password, log_in_with_saml, log_in_with_oidc, log_in_with_authenticator_app, log_in_with_passkey, vault_domain, email_send_from_domain, redirect_uri, cookie_domain)
	VALUES ('252491cc-76e3-4957-ab23-47d83c34f240'::uuid, 'Tesseral Test', true, true, true, true, true, true, true, true, 'vault.console.tesseral.example.com', 'vault.console.tesseral.example.com', 'https://console.tesseral.example.com', 'console.tesseral.example.com');

INSERT INTO project_trusted_domains (id, project_id, domain)
VALUES
    (gen_random_uuid(), '252491cc-76e3-4957-ab23-47d83c34f240', 'vault.console.tesseral.example.com'),
    (gen_random_uuid(), '252491cc-76e3-4957-ab23-47d83c34f240', 'console.tesseral.example.com');

-- Create the Dogfood Project's backing organization
INSERT INTO organizations (id, display_name, project_id, log_in_with_google, log_in_with_microsoft, log_in_with_email, log_in_with_password, log_in_with_saml, log_in_with_authenticator_app, log_in_with_passkey, scim_enabled)
  VALUES ('7a76decb-6d79-49ce-9449-34fcc53151df'::uuid, 'project_54vwf0clhh0caqe20eujxgpeq Backing Organization', '252491cc-76e3-4957-ab23-47d83c34f240', false, false, false, false, false, false, false, false);

UPDATE projects SET organization_id = '7a76decb-6d79-49ce-9449-34fcc53151df'::uuid where id = '252491cc-76e3-4957-ab23-47d83c34f240'::uuid;

-- Create project UI settings for the dogfood project
INSERT INTO project_ui_settings (id, project_id)
  VALUES (gen_random_uuid(), '252491cc-76e3-4957-ab23-47d83c34f240'::uuid);

-- Create a user in the dogfood project
INSERT INTO users (id, email, password_bcrypt, organization_id, is_owner)
  VALUES ('e071bbfe-6f27-4526-ab37-0ad251742836'::uuid, 'root@app.tesseral.example.com', crypt('password', gen_salt('bf', 14)), '7a76decb-6d79-49ce-9449-34fcc53151df', true);
`

	_, err := e.DB.Exec(context.Background(), sql)
	if err != nil {
		log.Panicf("failed to seed database: %v", err)
	}
}

func (e *Environment) NewProject(t *testing.T) (string, string) {
	projectID := uuid.New()
	formattedProjectID := idformat.Project.Format(projectID)
	projectVaultDomain := fmt.Sprintf("%s.%s", strings.ReplaceAll(formattedProjectID, "_", "-"), "example.com")

	// Create the backing organization for the new project
	organizationID := uuid.New()
	_, err := e.DB.Exec(t.Context(), `
INSERT INTO organizations (id, display_name, project_id, log_in_with_google, log_in_with_microsoft, log_in_with_email, log_in_with_password, log_in_with_saml, log_in_with_oidc, log_in_with_authenticator_app, log_in_with_passkey, scim_enabled)
  VALUES ($1::uuid, $2, $3::uuid, true, true, true, true, true, true, true, true, true);
`,
		organizationID.String(),
		fmt.Sprintf("%s Backing Organization", formattedProjectID),
		*e.DogfoodProjectID,
	)
	if err != nil {
		t.Fatalf("failed to create backing organization for test project: %v", err)
	}

	// Create the project with the new vault domain
	_, err = e.DB.Exec(t.Context(), `
INSERT INTO projects (id, organization_id, display_name, log_in_with_google, log_in_with_microsoft, log_in_with_github, log_in_with_email, log_in_with_password, log_in_with_saml, log_in_with_oidc, log_in_with_authenticator_app, log_in_with_passkey, vault_domain, email_send_from_domain, redirect_uri, cookie_domain, api_keys_enabled, api_key_secret_token_prefix, entitled_backend_api_keys, entitled_custom_vault_domains, audit_logs_enabled)
  VALUES ($1::uuid, $2::uuid, $3, true, true, true, true, true, true, true, true, true, $4, $4, $4, $4, true, 'test_sk_', true, true, true);
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
	_, err = e.DB.Exec(t.Context(), `
INSERT INTO project_ui_settings (id, project_id)
  VALUES (gen_random_uuid(), $1::uuid);
`,
		projectID.String(),
	)
	if err != nil {
		t.Fatalf("failed to create project UI settings for test project: %v", err)
	}

	// Create the project trusted domains
	_, err = e.DB.Exec(t.Context(), `
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

	// Create the project session signing key
	_, err = e.DB.Exec(t.Context(), `
INSERT INTO session_signing_keys (id, project_id, public_key, private_key_cipher_text, expire_time) 
  VALUES (
    gen_random_uuid(), 
    $1::uuid, 
    decode('3059301306072a8648ce3d020106082a8648ce3d03010703420004a82072a20d2217055f0c5f9f9283e128d5bc26334b19024c93f6ad50619bbe83bc565a2fbdc05e02dc3f1452ff273d7ec2534e2cbe7fe395443d887b128dd7b8', 'hex'), 
    decode('a1931242e0770f54e2e8365053ff4b72dc72faba0830cff2099655d78aa188f750b9b1557e70566f00449fed97a5b8a94a113e8049a6ea71436a08e135f35a7b86863f47f36e3e0b62dad8da491f28aba812a93e7a2a44913c6b2377c7ea4d89991eba682d9cfb17d5bcfa3f608e973dd61aa9910453e8d48058ea80ccbd0d5961de3fd25dcfe893dbdd84a43112d1533b4ebae65e35b0e8eca25b1af53eec97304899cb542ac850e59a6c5521ecbee5549329a451c8c948d82f1d6858a6d2680d987e72945ad5b4166c3529b70ce1106573874fb68847ed823567a9edfeac712d464ac5b339f80365be985ab69703d7100c65c872765b04a9ee575002edadef', 'hex'), 
    (SELECT NOW() + INTERVAL '1 year')
  );
`,
		projectID.String(),
	)
	if err != nil {
		t.Fatalf("failed to create project session signing key for test project: %v", err)
	}

	userID := uuid.New()
	formattedUserID := idformat.User.Format(userID)
	userEmail := fmt.Sprintf("%s@%s", formattedUserID, projectVaultDomain)

	_, err = e.DB.Exec(t.Context(), `
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

	return formattedProjectID, formattedUserID
}

func (e *Environment) NewOrganization(t *testing.T, projectID string, organization *backendv1.Organization) string {
	projectUUID, err := idformat.Project.Parse(projectID)
	if err != nil {
		t.Fatalf("failed to parse project ID: %v", err)
	}

	organizationID := uuid.New()
	formattedOrganizationID := idformat.Organization.Format(organizationID)

	// Create the organization
	_, err = e.DB.Exec(t.Context(), `
INSERT INTO organizations (id, display_name, project_id, log_in_with_google, log_in_with_microsoft, log_in_with_email, log_in_with_password, log_in_with_saml, log_in_with_oidc, log_in_with_authenticator_app, log_in_with_passkey, scim_enabled, api_keys_enabled, custom_roles_enabled)
  VALUES ($1::uuid, $2, $3::uuid, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14);
`,
		organizationID.String(),
		organization.DisplayName,
		uuid.UUID(projectUUID).String(),
		organization.GetLogInWithGoogle(),
		organization.GetLogInWithMicrosoft(),
		organization.GetLogInWithEmail(),
		organization.GetLogInWithPassword(),
		organization.GetLogInWithSaml(),
		organization.GetLogInWithOidc(),
		organization.GetLogInWithAuthenticatorApp(),
		organization.GetLogInWithPasskey(),
		organization.GetScimEnabled(),
		organization.GetApiKeysEnabled(),
		organization.GetCustomRolesEnabled(),
	)
	if err != nil {
		t.Fatalf("failed to create organization: %v", err)
	}

	return formattedOrganizationID
}

func (e *Environment) NewUser(t *testing.T, organizationID string, user *backendv1.User) string {
	organizationUUID, err := idformat.Organization.Parse(organizationID)
	if err != nil {
		t.Fatalf("failed to parse organization ID: %v", err)
	}

	userID := uuid.New()
	formattedUserID := idformat.User.Format(userID)

	if user.Email == "" {
		user.Email = fmt.Sprintf("%s@%s", formattedUserID, e.ConsoleDomain)
	}

	// Create the user
	_, err = e.DB.Exec(t.Context(), `
INSERT INTO users (id, email, password_bcrypt, organization_id, is_owner, display_name, google_user_id, microsoft_user_id, github_user_id)
  VALUES ($1::uuid, $2, crypt('password', gen_salt('bf', 14)), $3::uuid, $4, $5, $6, $7, $8);
`,
		userID.String(),
		user.Email,
		uuid.UUID(organizationUUID).String(),
		user.GetOwner(),
		user.DisplayName,
		user.GoogleUserId,
		user.MicrosoftUserId,
		user.GithubUserId,
	)
	if err != nil {
		t.Fatalf("failed to create user: %v", err)
	}
	return formattedUserID
}

func (e *Environment) NewSession(t *testing.T, userID string) (string, string) {
	userUUID, err := idformat.User.Parse(userID)
	if err != nil {
		t.Fatalf("failed to parse user ID: %v", err)
	}

	sessionID := uuid.New()
	formattedSessionID := idformat.Session.Format(sessionID)

	// Create the session
	refreshTokenID := uuid.New()
	refreshTokenSha256 := sha256.Sum256(refreshTokenID[:])
	refreshToken := idformat.SessionRefreshToken.Format(refreshTokenID)
	_, err = e.DB.Exec(t.Context(), `
INSERT INTO sessions (id, user_id, expire_time, refresh_token_sha256, primary_auth_factor)
  VALUES ($1::uuid, $2::uuid, $3, $4, $5);
`,
		sessionID.String(),
		uuid.UUID(userUUID).String(),
		time.Now().Add(24*time.Hour),
		refreshTokenSha256[:],
		"email",
	)
	if err != nil {
		t.Fatalf("failed to create session: %v", err)
	}
	return formattedSessionID, refreshToken
}

func (e *Environment) NewIntermediateSession(t *testing.T, projectID string) string {
	intermediateSessionID := uuid.New()
	projectUUID, err := idformat.Project.Parse(projectID)
	if err != nil {
		t.Fatalf("failed to parse project ID: %v", err)
	}

	const intermediateSessionDuration = time.Minute * 15
	expireTime := time.Now().Add(intermediateSessionDuration)

	secretToken := uuid.New()
	secretTokenSHA256 := sha256.Sum256(secretToken[:])
	_, err = e.DB.Exec(t.Context(), `
INSERT INTO intermediate_sessions (id, project_id, expire_time, secret_token_sha256)
  VALUES ($1::uuid, $2::uuid, $3, $4);
`,
		intermediateSessionID.String(),
		uuid.UUID(projectUUID).String(),
		expireTime,
		secretTokenSHA256[:],
	)
	if err != nil {
		t.Fatalf("failed to create intermediate session: %v", err)
	}

	return idformat.IntermediateSessionSecretToken.Format(secretToken)
}
