-- name: GetSessionSigningKeyPublicKey :one
SELECT
    public_key
FROM
    session_signing_keys
WHERE
    project_id = $1
    AND id = $2
    AND expire_time > @now;

-- name: CreateUser :one
INSERT INTO users (id, organization_id, email, password_bcrypt, google_user_id, microsoft_user_id)
    VALUES ($1, $2, $3, $4, $5, $6)
RETURNING
    *;

-- name: GetCurrentSessionKeyByProjectID :one
SELECT
    *
FROM
    session_signing_keys
WHERE
    project_id = $1
ORDER BY
    create_time DESC
LIMIT 1;

-- name: GetOrganizationByID :one
SELECT
    *
FROM
    organizations
WHERE
    id = $1;

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
    require_mfa = $9
WHERE
    id = $1
RETURNING
    *;

-- name: GetProjectByID :one
SELECT
    *
FROM
    projects
WHERE
    id = $1;

-- name: GetProjectTrustedDomains :many
SELECT
    *
FROM
    project_trusted_domains
WHERE
    project_id = $1;

-- name: GetSessionByID :one
SELECT
    *
FROM
    sessions
WHERE
    id = $1;

-- name: GetSessionDetailsByRefreshTokenSHA256 :one
SELECT
    sessions.id AS session_id,
    users.id AS user_id,
    organizations.id AS organization_id,
    projects.id AS project_id
FROM
    sessions
    JOIN users ON sessions.user_id = users.id
    JOIN organizations ON users.organization_id = organizations.id
    JOIN projects ON projects.id = organizations.project_id
WHERE
    refresh_token_sha256 = $1;

-- name: GetUserByID :one
SELECT
    *
FROM
    users
WHERE
    id = $1;

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
    *
FROM
    saml_connections
WHERE
    id = $1
    AND organization_id = $2;

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
    *
FROM
    scim_api_keys
WHERE
    id = $1
    AND organization_id = $2;

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

-- name: SetPassword :one
UPDATE
    users
SET
    update_time = now(),
    password_bcrypt = $2
WHERE
    id = $1
RETURNING
    *;

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
    *
FROM
    users
WHERE
    id = $1
    AND organization_id = $2;

-- name: UpdateUser :one
UPDATE
    users
SET
    update_time = now(),
    is_owner = $1
WHERE
    id = $2
RETURNING
    *;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;

-- name: InvalidateSession :exec
UPDATE
    sessions
SET
    expire_time = now(),
    refresh_token_sha256 = NULL
WHERE
    id = $1;

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

-- name: CreatePasskey :one
INSERT INTO passkeys (id, user_id, credential_id, public_key, aaguid, rp_id)
    VALUES ($1, $2, $3, $4, $5, $6)
RETURNING
    *;

-- name: GetUserPasskey :one
SELECT
    *
FROM
    passkeys
WHERE
    id = $1
    AND user_id = $2;

-- name: DeletePasskey :exec
DELETE FROM passkeys
WHERE id = $1;

-- name: CreateUserAuthenticatorAppChallenge :one
INSERT INTO user_authenticator_app_challenges (user_id, authenticator_app_secret_ciphertext)
    VALUES ($1, $2)
ON CONFLICT (user_id)
    DO UPDATE SET
        authenticator_app_secret_ciphertext = excluded.authenticator_app_secret_ciphertext
    RETURNING
        *;

-- name: GetUserAuthenticatorAppChallenge :one
SELECT
    *
FROM
    user_authenticator_app_challenges
WHERE
    user_id = $1;

-- name: DeleteUserAuthenticatorAppChallenge :exec
DELETE FROM user_authenticator_app_challenges
WHERE user_id = $1;

-- name: UpdateUserAuthenticatorApp :one
UPDATE
    users
SET
    authenticator_app_secret_ciphertext = $1,
    authenticator_app_recovery_code_sha256s = $2
WHERE
    id = $3
RETURNING
    *;

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
    *
FROM
    user_invites
WHERE
    id = $1
    AND organization_id = $2;

-- name: ExistsUserWithEmail :one
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

-- name: GetOrganizationGoogleHostedDomains :many
SELECT
    *
FROM
    organization_google_hosted_domains
WHERE
    organization_id = $1;

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
    *
FROM
    organization_microsoft_tenant_ids
WHERE
    organization_id = $1;

-- name: DeleteOrganizationMicrosoftTenantIDs :exec
DELETE FROM organization_microsoft_tenant_ids
WHERE organization_id = $1;

-- name: CreateOrganizationMicrosoftTenantID :one
INSERT INTO organization_microsoft_tenant_ids (id, organization_id, microsoft_tenant_id)
    VALUES ($1, $2, $3)
RETURNING
    *;

-- name: ListSwitchableOrganizations :many
SELECT
    id,
    display_name
FROM
    organizations
WHERE
    project_id = $1
    AND EXISTS (
        SELECT
            1
        FROM
            users
        WHERE
            organization_id = organizations.id
            AND users.email = $2);

-- name: GetProjectByBackingOrganizationID :one
SELECT
    *
FROM
    projects
WHERE
    organization_id = $1;

