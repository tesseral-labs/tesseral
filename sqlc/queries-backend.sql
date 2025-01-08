-- name: CreateOrganization :one
INSERT INTO organizations (id, project_id, display_name, google_hosted_domain, microsoft_tenant_id, override_log_in_methods, override_log_in_with_google_enabled, override_log_in_with_microsoft_enabled, override_log_in_with_password_enabled)
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
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
    display_name = $2,
    google_hosted_domain = $3,
    microsoft_tenant_id = $4,
    override_log_in_methods = $5,
    override_log_in_with_password_enabled = $6,
    override_log_in_with_google_enabled = $7,
    override_log_in_with_microsoft_enabled = $8
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
    display_name = $2,
    log_in_with_password_enabled = $3,
    log_in_with_google_enabled = $4,
    log_in_with_microsoft_enabled = $5,
    google_oauth_client_id = $6,
    google_oauth_client_secret_ciphertext = $7,
    microsoft_oauth_client_id = $8,
    microsoft_oauth_client_secret_ciphertext = $9,
    organizations_saml_enabled_default = $10,
    organizations_scim_enabled_default = $11
WHERE
    id = $1
RETURNING
    *;

-- name: UpdateProjectOrganizationID :one
UPDATE
    projects
SET
    organization_id = $2
WHERE
    id = $1
RETURNING
    *;

-- name: UpdateUser :one
UPDATE
    users
SET
    organization_id = $2,
    email = $3,
    password_bcrypt = $4,
    google_user_id = $5,
    microsoft_user_id = $6
WHERE
    id = $1
RETURNING
    *;

-- name: UpdateUserPassword :one
UPDATE
    users
SET
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
INSERT INTO scim_api_keys (id, organization_id, display_name, token_sha256)
    VALUES ($1, $2, $3, $4)
RETURNING
    *;

-- name: UpdateSCIMAPIKey :one
UPDATE
    scim_api_keys
SET
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
    token_sha256 = NULL
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
    secret_token_sha256 = NULL
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
    users.*
FROM
    users
    JOIN organizations ON users.organization_id = organizations.id
WHERE
    users.id = $1
    AND organizations.project_id = $2;

-- name: ListSessions :many
SELECT
    *
FROM
    sessions
WHERE
    user_id = $1
    AND id >= $2
ORDER BY
    id
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

-- name: ListIntermediateSessions :many
SELECT
    *
FROM
    intermediate_sessions
WHERE
    project_id = $1
    AND id >= $2
ORDER BY
    id
LIMIT $3;

-- name: GetIntermediateSession :one
SELECT
    intermediate_sessions.*
FROM
    intermediate_sessions
WHERE
    id = $1
    AND project_id = $2;

