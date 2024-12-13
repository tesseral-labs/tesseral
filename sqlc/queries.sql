-- name: CreateIntermediateSessionSigningKey :one
INSERT INTO intermediate_session_signing_keys (id, project_id, public_key, private_key_cipher_text, expire_time)
    VALUES ($1, $2, $3, $4, $5)
RETURNING
    *;

-- name: CreateOrganization :one
INSERT INTO organizations (id, project_id, display_name, google_hosted_domain, microsoft_tenant_id, override_log_in_with_google_enabled, override_log_in_with_microsoft_enabled, override_log_in_with_password_enabled)
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING
    *;

-- name: CreateProject :one
INSERT INTO projects (id, organization_id, google_oauth_client_id, google_oauth_client_secret_ciphertext, log_in_with_password_enabled, log_in_with_google_enabled, log_in_with_microsoft_enabled, microsoft_oauth_client_id, microsoft_oauth_client_secret_ciphertext)
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING
    *;

-- name: CreateSessionSigningKey :one
INSERT INTO session_signing_keys (id, project_id, public_key, private_key_cipher_text, expire_time)
    VALUES ($1, $2, $3, $4, $5)
RETURNING
    *;

-- name: CreateUser :one
INSERT INTO users (id, organization_id, email, password_bcrypt, google_user_id, microsoft_user_id)
    VALUES ($1, $2, $3, $4, $5, $6)
RETURNING
    *;

-- name: GetOrganizationByID :one
SELECT
    *
FROM
    organizations
WHERE
    id = $1;

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

-- name: GetSessionByID :one
SELECT
    *
FROM
    sessions
WHERE
    id = $1;

-- name: GetOrganizationByGoogleHostedDomain :one
SELECT
    *
FROM
    organizations
WHERE
    google_hosted_domain = $1;

-- name: GetUserByID :one
SELECT
    *
FROM
    users
WHERE
    id = $1;

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

-- -- name: UpdateProject :one
-- UPDATE
--     projects
-- SET
--     log_in_with_password_enabled = $2,
--     log_in_with_google_enabled = $3,
--     log_in_with_microsoft_enabled = $4,
--     google_oauth_client_id = $5,
--     google_oauth_client_secret = $6,
--     microsoft_oauth_client_id = $7,
--     microsoft_oauth_client_secret = $8
-- WHERE
--     id = $1
-- RETURNING
--     *;
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

-- name: CountAllProjects :one
SELECT
    count(*)
FROM
    projects;

-- name: CreateProjectAPIKey :one
INSERT INTO project_api_keys (id, project_id, create_time, revoked, secret_token_sha256)
    VALUES ($1, $2, $3, $4, $5)
RETURNING
    *;

-- name: GetProjectAPIKeyBySecretTokenSHA256 :one
SELECT
    *
FROM
    project_api_keys
WHERE
    secret_token_sha256 = $1;

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
    JOIN projects ON organizations.id = projects.organization_id
WHERE
    revoked = FALSE
    AND refresh_token_sha256 = $1;

