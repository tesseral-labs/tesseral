// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: queries-defaultoauth.sql

package queries

import (
	"context"
)

const getVaultDomainByGitHubOAuthStateSHA256 = `-- name: GetVaultDomainByGitHubOAuthStateSHA256 :one
SELECT
    vault_domain
FROM
    projects
    JOIN intermediate_sessions ON projects.id = intermediate_sessions.project_id
WHERE
    intermediate_sessions.github_oauth_state_sha256 = $1
`

func (q *Queries) GetVaultDomainByGitHubOAuthStateSHA256(ctx context.Context, githubOauthStateSha256 []byte) (string, error) {
	row := q.db.QueryRow(ctx, getVaultDomainByGitHubOAuthStateSHA256, githubOauthStateSha256)
	var vault_domain string
	err := row.Scan(&vault_domain)
	return vault_domain, err
}

const getVaultDomainByGoogleOAuthStateSHA256 = `-- name: GetVaultDomainByGoogleOAuthStateSHA256 :one
SELECT
    vault_domain
FROM
    projects
    JOIN intermediate_sessions ON projects.id = intermediate_sessions.project_id
WHERE
    intermediate_sessions.google_oauth_state_sha256 = $1
`

func (q *Queries) GetVaultDomainByGoogleOAuthStateSHA256(ctx context.Context, googleOauthStateSha256 []byte) (string, error) {
	row := q.db.QueryRow(ctx, getVaultDomainByGoogleOAuthStateSHA256, googleOauthStateSha256)
	var vault_domain string
	err := row.Scan(&vault_domain)
	return vault_domain, err
}

const getVaultDomainByMicrosoftOAuthStateSHA256 = `-- name: GetVaultDomainByMicrosoftOAuthStateSHA256 :one
SELECT
    vault_domain
FROM
    projects
    JOIN intermediate_sessions ON projects.id = intermediate_sessions.project_id
WHERE
    intermediate_sessions.microsoft_oauth_state_sha256 = $1
`

func (q *Queries) GetVaultDomainByMicrosoftOAuthStateSHA256(ctx context.Context, microsoftOauthStateSha256 []byte) (string, error) {
	row := q.db.QueryRow(ctx, getVaultDomainByMicrosoftOAuthStateSHA256, microsoftOauthStateSha256)
	var vault_domain string
	err := row.Scan(&vault_domain)
	return vault_domain, err
}
