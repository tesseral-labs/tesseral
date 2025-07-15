-- name: GetVaultDomainByGoogleOAuthStateSHA256 :one
SELECT
    vault_domain
FROM
    projects
    JOIN intermediate_sessions ON projects.id = intermediate_sessions.project_id
WHERE
    intermediate_sessions.google_oauth_state_sha256 = $1;

-- name: GetVaultDomainByMicrosoftOAuthStateSHA256 :one
SELECT
    vault_domain
FROM
    projects
    JOIN intermediate_sessions ON projects.id = intermediate_sessions.project_id
WHERE
    intermediate_sessions.microsoft_oauth_state_sha256 = $1;

-- name: GetVaultDomainByGitHubOAuthStateSHA256 :one
SELECT
    vault_domain
FROM
    projects
    JOIN intermediate_sessions ON projects.id = intermediate_sessions.project_id
WHERE
    intermediate_sessions.github_oauth_state_sha256 = $1;

