-- name: CreateOrganization :one
INSERT INTO organizations (id, project_id, display_name, scim_enabled, log_in_with_email, log_in_with_password, log_in_with_google, log_in_with_microsoft)
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING
    *;

-- name: CreateDogfoodProject :one
INSERT INTO projects (id, display_name, redirect_uri, log_in_with_google, log_in_with_microsoft, log_in_with_email, log_in_with_password, vault_domain, email_send_from_domain, cookie_domain)
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
RETURNING
    *;

-- name: CreateProjectTrustedDomain :one
INSERT INTO project_trusted_domains (id, project_id, DOMAIN)
    VALUES ($1, $2, $3)
RETURNING
    *;

-- name: CreateProjectUISettings :one
INSERT INTO project_ui_settings (id, project_id, primary_color, detect_dark_mode_enabled)
    VALUES ($1, $2, $3, $4)
RETURNING
    *;

-- name: CreateSessionSigningKey :one
INSERT INTO session_signing_keys (id, project_id, public_key, private_key_cipher_text, expire_time)
    VALUES ($1, $2, $3, $4, $5)
RETURNING
    *;

-- name: CreateUser :one
INSERT INTO users (id, organization_id, email, is_owner, password_bcrypt, google_user_id, microsoft_user_id)
    VALUES ($1, $2, $3, $4, $5, $6, $7)
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

-- -- name: GetOrganizationByGoogleHostedDomain :one
-- SELECT
--     *
-- FROM
--     organizations
-- WHERE
--     google_hosted_domain = $1;
-- name: GetUserByID :one
SELECT
    *
FROM
    users
WHERE
    id = $1;

-- -- name: UpdateOrganization :one
-- UPDATE
--     organizations
-- SET
--     display_name = $2,
--     google_hosted_domain = $3,
--     microsoft_tenant_id = $4,
--     override_log_in_with_password_enabled = $5,
--     override_log_in_with_google_enabled = $6,
--     override_log_in_with_microsoft_enabled = $7
-- WHERE
--     id = $1
-- RETURNING
--     *;
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

