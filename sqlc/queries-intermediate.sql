-- name: GetIntermediateSessionByTokenSHA256AndProjectID :one
SELECT
    *
FROM
    intermediate_sessions
WHERE
    secret_token_sha256 = $1
    AND project_id = $2;

-- name: UpdateIntermediateSessionEmail :one
UPDATE
    intermediate_sessions
SET
    email = $1
WHERE
    id = $2
    AND (email IS NULL
        OR email = $1)
RETURNING
    *;

-- name: CreateEmailVerificationChallenge :one
INSERT INTO email_verification_challenges (id, intermediate_session_id, challenge_sha256, expire_time)
    VALUES ($1, $2, $3, $4)
RETURNING
    *;

-- name: GetEmailVerificationChallengeByChallengeSHA :one
SELECT
    *
FROM
    email_verification_challenges
WHERE
    intermediate_session_id = $1
    AND expire_time > now()
    AND challenge_sha256 = $2;

-- name: CompleteEmailVerificationChallenge :one
UPDATE
    email_verification_challenges
SET
    complete_time = now(),
    challenge_sha256 = NULL
WHERE
    id = $1
RETURNING
    *;

-- name: CreateVerifiedEmail :one
INSERT INTO verified_emails (id, project_id, email, google_user_id, microsoft_user_id)
    VALUES ($1, $2, $3, $4, $5)
RETURNING
    *;

-- name: CreateIntermediateSession :one
INSERT INTO intermediate_sessions (id, project_id, expire_time, email, google_user_id, microsoft_user_id, secret_token_sha256)
    VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING
    *;

-- name: CreateOrganization :one
INSERT INTO organizations (id, project_id, display_name, override_log_in_methods, override_log_in_with_google_enabled, override_log_in_with_microsoft_enabled, override_log_in_with_password_enabled, saml_enabled, scim_enabled)
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING
    *;

-- name: CreateSession :one
INSERT INTO sessions (id, user_id, expire_time, refresh_token_sha256, revoked)
    VALUES ($1, $2, $3, $4, $5)
RETURNING
    *;

-- name: CreateUser :one
INSERT INTO users (id, organization_id, email, google_user_id, microsoft_user_id, is_owner)
    VALUES ($1, $2, $3, $4, $5, $6)
RETURNING
    *;

-- name: GetEmailVerificationChallengeByID :one
SELECT
    *
FROM
    email_verification_challenges
WHERE
    id = $1;

-- name: GetEmailVerificationChallengeForCompletion :one
SELECT
    *
FROM
    email_verification_challenges
WHERE
    intermediate_session_id = $1
    AND expire_time > $2
    AND challenge_sha256 IS NOT NULL
    AND complete_time IS NULL
ORDER BY
    create_time DESC
LIMIT 1;

-- name: GetIntermediateSessionByID :one
SELECT
    *
FROM
    intermediate_sessions
WHERE
    id = $1;

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

-- name: GetOrganizationUserByEmail :one
SELECT
    *
FROM
    users
WHERE
    organization_id = $1
    AND email = $2;

-- name: GetOrganizationUserByGoogleUserID :one
SELECT
    *
FROM
    users
WHERE
    organization_id = $1
    AND google_user_id = $2;

-- name: GetOrganizationUserByMicrosoftUserID :one
SELECT
    *
FROM
    users
WHERE
    organization_id = $1
    AND microsoft_user_id = $2;

-- name: GetOrganizationPrimarySAMLConnection :one
SELECT
    *
FROM
    saml_connections
WHERE
    organization_id = $1
    AND is_primary = TRUE
LIMIT 1;

-- name: GetProjectOrganizationByID :one
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

-- name: GetProjectUISettings :one
SELECT
    *
FROM
    project_ui_settings
WHERE
    project_id = $1;

-- name: GetSessionSigningKeysByProjectID :many
SELECT
    *
FROM
    session_signing_keys
WHERE
    project_id = $1;

-- name: IsGoogleEmailVerified :one
SELECT
    count(*) > 0
FROM
    verified_emails
WHERE
    project_id = $1
    AND email = $2
    AND google_user_id = $3;

-- name: IsMicrosoftEmailVerified :one
SELECT
    count(*) > 0
FROM
    verified_emails
WHERE
    project_id = $1
    AND email = $2
    AND microsoft_user_id = $3;

-- name: ListOrganizationsByEmail :many
SELECT
    organizations.*
FROM
    organizations
    JOIN users ON organizations.id = users.organization_id
WHERE
    organizations.project_id = $1
    AND users.email = $2
    AND users.google_user_id IS NULL
    AND users.microsoft_user_id IS NULL;

-- name: ListOrganizationsByGoogleHostedDomain :many
SELECT
    organizations.*
FROM
    organizations
    JOIN organization_google_hosted_domains ON organizations.id = organization_google_hosted_domains.organization_id
WHERE
    organizations.project_id = $1
    AND organization_google_hosted_domains.google_hosted_domain = $2;

-- name: ListOrganizationsByMicrosoftTenantID :many
SELECT
    organizations.*
FROM
    organizations
    JOIN organization_microsoft_tenant_ids ON organizations.id = organization_microsoft_tenant_ids.organization_id
WHERE
    organizations.project_id = $1
    AND organization_microsoft_tenant_ids.microsoft_tenant_id = $2;

-- name: ListOrganizationsByMatchingUser :many
SELECT
    organizations.*
FROM
    organizations
    JOIN users ON organizations.id = users.organization_id
WHERE
    organizations.project_id = $1
    AND (users.email = $2
        OR (users.google_user_id IS NOT NULL
            AND users.google_user_id = $3)
        OR (users.microsoft_user_id IS NOT NULL
            AND users.microsoft_user_id = $4));

-- name: ListSAMLOrganizations :many
SELECT
    organizations.*
FROM
    organizations
    JOIN organization_domains ON organizations.id = organization_domains.organization_id
WHERE
    organizations.project_id = $1
    AND organizations.saml_enabled = TRUE
    AND organization_domains.domain = $2;

-- name: ListUsersByEmail :many
SELECT
    *
FROM
    users
WHERE
    email = $1;

-- name: RevokeEmailVerificationChallenge :one
UPDATE
    email_verification_challenges
SET
    challenge_sha256 = NULL
WHERE
    id = $1
RETURNING
    *;

-- name: RevokeIntermediateSession :one
UPDATE
    intermediate_sessions
SET
    secret_token_sha256 = NULL
WHERE
    id = $1
RETURNING
    *;

-- name: UpdateIntermediateSessionGoogleOAuthStateSHA256 :one
UPDATE
    intermediate_sessions
SET
    google_oauth_state_sha256 = $1
WHERE
    id = $2
RETURNING
    *;

-- name: UpdateIntermediateSessionGoogleDetails :one
UPDATE
    intermediate_sessions
SET
    email = $1,
    google_user_id = $2,
    google_hosted_domain = $3
WHERE
    id = $4
RETURNING
    *;

-- name: UpdateIntermediateSessionMicrosoftOAuthStateSHA256 :one
UPDATE
    intermediate_sessions
SET
    microsoft_oauth_state_sha256 = $1
WHERE
    id = $2
RETURNING
    *;

-- name: UpdateIntermediateSessionMicrosoftDetails :one
UPDATE
    intermediate_sessions
SET
    email = $1,
    microsoft_user_id = $2,
    microsoft_tenant_id = $3
WHERE
    id = $4
RETURNING
    *;

-- name: UpdateIntermediateSessionPasswordVerified :one
UPDATE
    intermediate_sessions
SET
    password_verified = TRUE,
    organization_id = $1
WHERE
    id = $2
RETURNING
    *;

-- name: GetUserByOrganizationIDAndFactors :one
SELECT
    *
FROM
    users
WHERE
    organization_id = $1
    AND (google_user_id = $2
        OR microsoft_user_id = $3
        OR email = $4);

-- name: GetEmailVerifiedByGoogleUserID :one
SELECT
    EXISTS (
        SELECT
            *
        FROM
            verified_emails
        WHERE
            project_id = $1
            AND email = $2
            AND google_user_id = $3);

-- name: GetEmailVerifiedByMicrosoftUserID :one
SELECT
    EXISTS (
        SELECT
            *
        FROM
            verified_emails
        WHERE
            project_id = $1
            AND email = $2
            AND microsoft_user_id = $3);

-- name: GetEmailVerifiedByEmailVerificationChallenge :one
SELECT
    EXISTS (
        SELECT
            *
        FROM
            email_verification_challenges
        WHERE
            complete_time IS NOT NULL
            AND intermediate_session_id = $1);

-- name: CreateOrganizationGoogleHostedDomain :one
INSERT INTO organization_google_hosted_domains (id, organization_id, google_hosted_domain)
    VALUES ($1, $2, $3)
RETURNING
    *;

-- name: CreateOrganizationMicrosoftTenantID :one
INSERT INTO organization_microsoft_tenant_ids (id, organization_id, microsoft_tenant_id)
    VALUES ($1, $2, $3)
RETURNING
    *;

-- name: UpdateUserFailedPasswordAttempts :one
UPDATE
    users
SET
    failed_password_attempts = $1
WHERE
    id = $2
RETURNING
    *;

-- name: UpdateUserPasswordLockoutExpireTime :one
UPDATE
    users
SET
    password_lockout_expire_time = $1
WHERE
    id = $2
RETURNING
    *;

