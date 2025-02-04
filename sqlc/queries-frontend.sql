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
    log_in_with_password = $3,
    log_in_with_google = $4,
    log_in_with_microsoft = $5,
    log_in_with_authenticator_app = $6,
    log_in_with_passkey = $7
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
    revoked = FALSE
    AND refresh_token_sha256 = $1;

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

-- name: InvalidateSession :exec
UPDATE
    sessions
SET
    update_time = now(),
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
INSERT INTO passkeys (id, user_id, credential_id, public_key, aaguid)
    VALUES ($1, $2, $3, $4, $5)
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
    authenticator_app_recovery_code_bcrypts = $2
WHERE
    id = $3
RETURNING
    *;
