-- name: CreateOrganization :one
INSERT INTO organizations (id, project_id, display_name, google_hosted_domain, microsoft_tenant_id, override_log_in_with_google_enabled, override_log_in_with_microsoft_enabled, override_log_in_with_password_enabled)
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING
    *;

-- name: CreateProject :one
INSERT INTO projects (id, organization_id, log_in_with_password_enabled, log_in_with_google_enabled, log_in_with_microsoft_enabled, google_oauth_client_id, google_oauth_client_secret, microsoft_oauth_client_id, microsoft_oauth_client_secret)
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
    override_log_in_with_password_enabled = $5,
    override_log_in_with_google_enabled = $6,
    override_log_in_with_microsoft_enabled = $7
WHERE
    id = $1
RETURNING
    *;

-- name: UpdateProject :one
UPDATE
    projects
SET
    log_in_with_password_enabled = $2,
    log_in_with_google_enabled = $3,
    log_in_with_microsoft_enabled = $4,
    google_oauth_client_id = $5,
    google_oauth_client_secret = $6,
    microsoft_oauth_client_id = $7,
    microsoft_oauth_client_secret = $8
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
    unverified_email = $3,
    verified_email = $4,
    password_bcrypt = $5,
    google_user_id = $6,
    microsoft_user_id = $7
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

