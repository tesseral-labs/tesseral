-- name: CreateOrganization :one
INSERT INTO organizations (id, project_id, display_name, log_in_with_google, log_in_with_microsoft, log_in_with_email, log_in_with_password, log_in_with_saml, log_in_with_authenticator_app, log_in_with_passkey, scim_enabled)
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
RETURNING
    *;

-- name: GetOrganizationByProjectIDAndID :one
SELECT
    *
FROM
    organizations
WHERE
    id = $1
    AND project_id = $2;

-- name: GetProjectByID :one
SELECT
    *
FROM
    projects
WHERE
    id = $1;

-- name: GetProjectAPIKeyBySecretTokenSHA256 :one
SELECT
    *
FROM
    project_api_keys
WHERE
    secret_token_sha256 = $1;

-- name: GetSessionSigningKeysByProjectID :many
SELECT
    *
FROM
    session_signing_keys
WHERE
    project_id = $1;

-- name: ListOrganizationsByProjectId :many
SELECT
    *
FROM
    organizations
WHERE
    project_id = $1
ORDER BY
    id
LIMIT $2;

-- name: GetProjectIDOrganizationBacks :one
SELECT
    projects.id
FROM
    organizations
    JOIN projects ON projects.organization_id = organizations.id
WHERE
    organization_id = $1;

-- name: ListProjects :many
SELECT
    *
FROM
    projects
ORDER BY
    id
LIMIT $1;

-- name: UpdateOrganization :one
UPDATE
    organizations
SET
    update_time = now(),
    display_name = $2,
    log_in_with_google = $3,
    log_in_with_microsoft = $4,
    log_in_with_email = $5,
    log_in_with_password = $6,
    log_in_with_authenticator_app = $7,
    log_in_with_passkey = $8,
    log_in_with_saml = $9,
    scim_enabled = $10
WHERE
    id = $1
RETURNING
    *;

-- name: DeleteOrganization :exec
DELETE FROM organizations
WHERE id = $1;

-- name: UpdateProject :one
UPDATE
    projects
SET
    update_time = now(),
    display_name = $2,
    log_in_with_google = $3,
    log_in_with_microsoft = $4,
    log_in_with_email = $5,
    log_in_with_password = $6,
    log_in_with_saml = $7,
    log_in_with_authenticator_app = $8,
    log_in_with_passkey = $9,
    google_oauth_client_id = $10,
    google_oauth_client_secret_ciphertext = $11,
    microsoft_oauth_client_id = $12,
    microsoft_oauth_client_secret_ciphertext = $13,
    redirect_uri = $14,
    after_login_redirect_uri = $15,
    after_signup_redirect_uri = $16
WHERE
    id = $1
RETURNING
    *;

-- name: GetProjectTrustedDomains :many
SELECT
    *
FROM
    project_trusted_domains
WHERE
    project_id = $1;

-- name: DeleteProjectTrustedDomainsByProjectID :exec
DELETE FROM project_trusted_domains
WHERE project_id = $1;

-- name: CreateProjectTrustedDomain :one
INSERT INTO project_trusted_domains (id, project_id, DOMAIN)
    VALUES ($1, $2, $3)
RETURNING
    *;

-- name: DisableProjectOrganizationsLogInWithGoogle :exec
UPDATE
    organizations
SET
    log_in_with_google = FALSE
WHERE
    project_id = $1;

-- name: DisableProjectOrganizationsLogInWithMicrosoft :exec
UPDATE
    organizations
SET
    log_in_with_microsoft = FALSE
WHERE
    project_id = $1;

-- name: DisableProjectOrganizationsLogInWithEmail :exec
UPDATE
    organizations
SET
    log_in_with_email = FALSE
WHERE
    project_id = $1;

-- name: DisableProjectOrganizationsLogInWithPassword :exec
UPDATE
    organizations
SET
    log_in_with_password = FALSE
WHERE
    project_id = $1;

-- name: DisableProjectOrganizationsLogInWithSAML :exec
UPDATE
    organizations
SET
    log_in_with_saml = FALSE
WHERE
    project_id = $1;

-- name: DisableProjectOrganizationsLogInWithAuthenticatorApp :exec
UPDATE
    organizations
SET
    log_in_with_authenticator_app = FALSE
WHERE
    project_id = $1;

-- name: DisableProjectOrganizationsLogInWithPasskey :exec
UPDATE
    organizations
SET
    log_in_with_passkey = FALSE
WHERE
    project_id = $1;

-- name: UpdateProjectOrganizationID :one
UPDATE
    projects
SET
    organization_id = $2
WHERE
    id = $1
RETURNING
    *;

-- name: UpdateUserPassword :one
UPDATE
    users
SET
    update_time = now(),
    password_bcrypt = $2
WHERE
    id = $1
RETURNING
    *;

-- name: ListSAMLConnections :many
SELECT
    *
FROM
    saml_connections
WHERE
    organization_id = $1
    AND id >= $2
ORDER BY
    id
LIMIT $3;

-- name: GetSAMLConnection :one
SELECT
    saml_connections.*
FROM
    saml_connections
    JOIN organizations ON saml_connections.organization_id = organizations.id
WHERE
    saml_connections.id = $1
    AND organizations.project_id = $2;

-- name: CreateSAMLConnection :one
INSERT INTO saml_connections (id, organization_id, is_primary, idp_redirect_url, idp_x509_certificate, idp_entity_id)
    VALUES ($1, $2, $3, $4, $5, $6)
RETURNING
    *;

-- name: UpdatePrimarySAMLConnection :exec
UPDATE
    saml_connections
SET
    is_primary = (id = $1)
WHERE
    organization_id = $2;

-- name: UpdateSAMLConnection :one
UPDATE
    saml_connections
SET
    update_time = now(),
    is_primary = $1,
    idp_redirect_url = $2,
    idp_x509_certificate = $3,
    idp_entity_id = $4
WHERE
    id = $5
RETURNING
    *;

-- name: DeleteSAMLConnection :exec
DELETE FROM saml_connections
WHERE id = $1;

-- name: ListSCIMAPIKeys :many
SELECT
    *
FROM
    scim_api_keys
WHERE
    organization_id = $1
    AND id >= $2
ORDER BY
    id
LIMIT $3;

-- name: GetSCIMAPIKey :one
SELECT
    scim_api_keys.*
FROM
    scim_api_keys
    JOIN organizations ON scim_api_keys.organization_id = organizations.id
WHERE
    scim_api_keys.id = $1
    AND organizations.project_id = $2;

-- name: CreateSCIMAPIKey :one
INSERT INTO scim_api_keys (id, organization_id, display_name, secret_token_sha256)
    VALUES ($1, $2, $3, $4)
RETURNING
    *;

-- name: UpdateSCIMAPIKey :one
UPDATE
    scim_api_keys
SET
    update_time = now(),
    display_name = $1
WHERE
    id = $2
RETURNING
    *;

-- name: DeleteSCIMAPIKey :exec
DELETE FROM scim_api_keys
WHERE id = $1;

-- name: RevokeSCIMAPIKey :one
UPDATE
    scim_api_keys
SET
    update_time = now(),
    secret_token_sha256 = NULL
WHERE
    id = $1
RETURNING
    *;

-- name: ListProjectAPIKeys :many
SELECT
    *
FROM
    project_api_keys
WHERE
    project_id = $1
    AND id >= $2
ORDER BY
    id
LIMIT $3;

-- name: GetProjectAPIKey :one
SELECT
    *
FROM
    project_api_keys
WHERE
    id = $1
    AND project_id = $2;

-- name: CreateProjectAPIKey :one
INSERT INTO project_api_keys (id, project_id, display_name, secret_token_sha256)
    VALUES ($1, $2, $3, $4)
RETURNING
    *;

-- name: UpdateProjectAPIKey :one
UPDATE
    project_api_keys
SET
    update_time = now(),
    display_name = $1
WHERE
    id = $2
RETURNING
    *;

-- name: DeleteProjectAPIKey :exec
DELETE FROM project_api_keys
WHERE id = $1;

-- name: RevokeProjectAPIKey :one
UPDATE
    project_api_keys
SET
    update_time = now(),
    secret_token_sha256 = NULL
WHERE
    id = $1
RETURNING
    *;

-- name: ListPublishableKeys :many
SELECT
    *
FROM
    publishable_keys
WHERE
    project_id = $1
    AND id >= $2
ORDER BY
    id
LIMIT $3;

-- name: GetPublishableKey :one
SELECT
    *
FROM
    publishable_keys
WHERE
    id = $1
    AND project_id = $2;

-- name: CreatePublishableKey :one
INSERT INTO publishable_keys (id, project_id, display_name)
    VALUES ($1, $2, $3)
RETURNING
    *;

-- name: UpdatePublishableKey :one
UPDATE
    publishable_keys
SET
    update_time = now(),
    display_name = $1
WHERE
    id = $2
RETURNING
    *;

-- name: DeletePublishableKey :exec
DELETE FROM publishable_keys
WHERE id = $1;

-- name: ListUsers :many
SELECT
    *
FROM
    users
WHERE
    organization_id = $1
    AND id >= $2
ORDER BY
    id
LIMIT $3;

-- name: GetUser :one
SELECT
    users.*
FROM
    users
    JOIN organizations ON users.organization_id = organizations.id
WHERE
    users.id = $1
    AND organizations.project_id = $2;

-- name: CreateUser :one
INSERT INTO users (id, organization_id, google_user_id, microsoft_user_id, email, is_owner)
    VALUES ($1, $2, $3, $4, $5, $6)
RETURNING
    *;

-- name: UpdateUser :one
UPDATE
    users
SET
    update_time = now(),
    email = $2,
    google_user_id = $3,
    microsoft_user_id = $4,
    is_owner = $5
WHERE
    id = $1
RETURNING
    *;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;

-- name: ListSessions :many
SELECT
    *
FROM
    sessions
WHERE
    user_id = $1
    AND id <= $2
ORDER BY
    id DESC
LIMIT $3;

-- name: GetSession :one
SELECT
    sessions.*
FROM
    sessions
    JOIN users ON sessions.user_id = users.id
    JOIN organizations ON users.organization_id = organizations.id
WHERE
    sessions.id = $1
    AND organizations.project_id = $2;

-- name: GetProjectUISettings :one
SELECT
    *
FROM
    project_ui_settings
WHERE
    project_id = $1;

-- name: UpdateProjectUISettings :one
UPDATE
    project_ui_settings
SET
    update_time = now(),
    primary_color = $3,
    detect_dark_mode_enabled = $4,
    dark_mode_primary_color = $5
WHERE
    id = $1
    AND project_id = $2
RETURNING
    *;

-- name: GetUserForImpersonation :one
SELECT
    users.*
FROM
    users
    JOIN organizations ON users.organization_id = organizations.id
    JOIN projects ON organizations.id = projects.organization_id
WHERE
    users.id = $1
    AND projects.organization_id = @impersonator_organization_id;

-- name: CreateUserImpersonationToken :one
INSERT INTO user_impersonation_tokens (id, impersonator_id, impersonated_id, expire_time, secret_token_sha256)
    VALUES ($1, $2, $3, $4, $5)
RETURNING
    *;

-- name: RevokeAllOrganizationSessions :exec
UPDATE
    sessions
SET
    refresh_token_sha256 = NULL
WHERE
    user_id IN (
        SELECT
            id
        FROM
            users
        WHERE
            organization_id = $1);

-- name: RevokeAllProjectSessions :exec
UPDATE
    sessions
SET
    refresh_token_sha256 = NULL
WHERE
    user_id IN (
        SELECT
            id
        FROM
            users
        WHERE
            organization_id IN (
                SELECT
                    id
                FROM
                    organizations
                WHERE
                    project_id = $1));

-- name: DisableOrganizationLogins :exec
UPDATE
    organizations
SET
    logins_disabled = TRUE
WHERE
    id = $1;

-- name: EnableOrganizationLogins :exec
UPDATE
    organizations
SET
    logins_disabled = FALSE
WHERE
    id = $1;

-- name: DisableProjectLogins :exec
UPDATE
    projects
SET
    logins_disabled = TRUE
WHERE
    id = $1;

-- name: EnableProjectLogins :exec
UPDATE
    projects
SET
    logins_disabled = FALSE
WHERE
    id = $1;

-- name: GetOrganizationGoogleHostedDomains :many
SELECT
    organization_google_hosted_domains.*
FROM
    organization_google_hosted_domains
    JOIN organizations ON organization_google_hosted_domains.organization_id = organizations.id
WHERE
    public.organization_google_hosted_domains.organization_id = $1
    AND organizations.project_id = $2;

-- name: DeleteOrganizationGoogleHostedDomains :exec
DELETE FROM organization_google_hosted_domains
WHERE organization_id = $1;

-- name: CreateOrganizationGoogleHostedDomain :one
INSERT INTO organization_google_hosted_domains (id, organization_id, google_hosted_domain)
    VALUES ($1, $2, $3)
RETURNING
    *;

-- name: GetOrganizationMicrosoftTenantIDs :many
SELECT
    organization_microsoft_tenant_ids.*
FROM
    organization_microsoft_tenant_ids
    JOIN organizations ON organization_microsoft_tenant_ids.organization_id = organizations.id
WHERE
    public.organization_microsoft_tenant_ids.organization_id = $1
    AND organizations.project_id = $2;

-- name: DeleteOrganizationMicrosoftTenantIDs :exec
DELETE FROM organization_microsoft_tenant_ids
WHERE organization_id = $1;

-- name: CreateOrganizationMicrosoftTenantID :one
INSERT INTO organization_microsoft_tenant_ids (id, organization_id, microsoft_tenant_id)
    VALUES ($1, $2, $3)
RETURNING
    *;

-- name: ListPasskeys :many
SELECT
    *
FROM
    passkeys
WHERE
    user_id = $1
    AND id >= $2
ORDER BY
    id
LIMIT $3;

-- name: GetPasskey :one
SELECT
    passkeys.*
FROM
    passkeys
    JOIN users ON passkeys.user_id = users.id
    JOIN organizations ON users.organization_id = organizations.id
WHERE
    passkeys.id = $1
    AND organizations.project_id = $2;

-- name: UpdatePasskey :one
UPDATE
    passkeys
SET
    update_time = now(),
    disabled = $2
WHERE
    id = $1
RETURNING
    *;

-- name: DeletePasskey :exec
DELETE FROM passkeys
WHERE id = $1;

-- name: ListUserInvites :many
SELECT
    *
FROM
    user_invites
WHERE
    organization_id = $1
    AND id >= $2
ORDER BY
    id
LIMIT $3;

-- name: GetUserInvite :one
SELECT
    user_invites.*
FROM
    user_invites
    JOIN organizations ON user_invites.organization_id = organizations.id
WHERE
    user_invites.id = $1
    AND organizations.project_id = $2;

-- name: ExistsUserWithEmailInOrganization :one
SELECT
    EXISTS (
        SELECT
            *
        FROM
            users
        WHERE
            organization_id = $1
            AND email = $2);

-- name: CreateUserInvite :one
INSERT INTO user_invites (id, organization_id, email, is_owner)
    VALUES ($1, $2, $3, $4)
ON CONFLICT (organization_id, email)
    DO UPDATE SET
        email = excluded.email, is_owner = excluded.is_owner
    RETURNING
        *;

-- name: DeleteUserInvite :exec
DELETE FROM user_invites
WHERE id = $1;

-- name: GetVaultDomainSettings :one
SELECT
    *
FROM
    vault_domain_settings
WHERE
    project_id = $1;

-- name: UpsertVaultDomainSettings :one
INSERT INTO vault_domain_settings (project_id, pending_domain)
    VALUES ($1, $2)
ON CONFLICT (project_id)
    DO UPDATE SET
        pending_domain = excluded.pending_domain
    RETURNING
        *;

-- name: DeleteVaultDomainSettings :exec
DELETE FROM vault_domain_settings
WHERE project_id = $1;

-- name: UpdateProjectVaultDomain :one
UPDATE
    projects
SET
    vault_domain = $2
WHERE
    id = $1
RETURNING
    *;

-- name: DisablePasskeysWithOldRPID :exec
UPDATE
    passkeys
SET
    disabled = TRUE,
    update_time = now()
FROM
    users,
    organizations,
    projects
WHERE
    passkeys.rp_id != projects.vault_domain
    AND passkeys.user_id = users.id
    AND users.organization_id = organizations.id
    AND organizations.project_id = projects.id
    AND projects.id = $1;

-- name: UpdateProjectEmailSendFromDomain :one
UPDATE
    projects
SET
    email_send_from_domain = $2
WHERE
    id = $1
RETURNING
    *;

-- name: GetVaultDomainInActiveOrPendingUse :one
SELECT
    EXISTS (
        SELECT
            1
        FROM
            projects
        WHERE
            vault_domain = $1)
    OR EXISTS (
        SELECT
            1
        FROM
            vault_domain_settings
        WHERE
            pending_domain = $1);

