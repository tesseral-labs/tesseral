-- name: GetProject :one
SELECT
    *
FROM
    projects
WHERE
    id = $1;

-- name: GetOIDCConnection :one
SELECT
    oidc_connections.*
FROM
    oidc_connections
    JOIN organizations ON oidc_connections.organization_id = organizations.id
WHERE
    organizations.project_id = $1
    AND organizations.log_in_with_oidc
    AND oidc_connections.id = $2;

-- name: GetOrganizationDomains :many
SELECT
    DOMAIN
FROM
    organization_domains
WHERE
    organization_id = $1;

-- name: GetIntermediateSessionByTokenSHA256AndProjectID :one
SELECT
    *
FROM
    intermediate_sessions
WHERE
    secret_token_sha256 = $1
    AND project_id = $2;

-- name: InitIntermediateSession :exec
UPDATE
    intermediate_sessions
SET
    oidc_state = $2,
    oidc_code_verifier = $3,
    organization_id = $4,
    primary_auth_factor = 'oidc'
WHERE
    id = $1;

-- name: UpdateIntermediateSession :exec
UPDATE
    intermediate_sessions
SET
    email = $2,
    verified_oidc_connection_id = $3
WHERE
    id = $1;

-- name: CreateAuditLogEvent :one
INSERT INTO audit_log_events (id, project_id, organization_id, actor_user_id, actor_session_id, resource_type, resource_id, event_name, event_time, event_details)
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, coalesce(@event_details, '{}'::jsonb))
RETURNING
    *;

