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

-- name: CreateVerifiedEmail :one
INSERT INTO oauth_verified_emails (id, project_id, email, google_user_id, microsoft_user_id)
    VALUES ($1, $2, $3, $4, $5)
RETURNING
    *;

-- name: CreateIntermediateSession :one
INSERT INTO intermediate_sessions (id, project_id, expire_time, email, google_user_id, microsoft_user_id, secret_token_sha256)
    VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING
    *;

-- name: CreateOrganization :one
INSERT INTO organizations (id, project_id, display_name, log_in_with_google, log_in_with_microsoft, log_in_with_password, log_in_with_email, scim_enabled)
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING
    *;

-- name: CreateSession :one
INSERT INTO sessions (id, user_id, expire_time, refresh_token_sha256)
    VALUES ($1, $2, $3, $4)
RETURNING
    *;

-- name: CreateUser :one
INSERT INTO users (id, organization_id, email, google_user_id, microsoft_user_id, is_owner, password_bcrypt)
    VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING
    *;

-- name: GetIntermediateSessionByID :one
SELECT
    *
FROM
    intermediate_sessions
WHERE
    id = $1;

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

-- name: GetProjectByBackingOrganizationID :one
SELECT
    *
FROM
    projects
WHERE
    organization_id = $1;

-- name: GetProjectTrustedDomains :many
SELECT
    *
FROM
    project_trusted_domains
WHERE
    project_id = $1;

-- name: GetProjectUISettings :one
SELECT
    *
FROM
    project_ui_settings
WHERE
    project_id = $1;

-- name: ListOrganizationsByGoogleHostedDomain :many
SELECT
    organizations.*
FROM
    organizations
    JOIN organization_google_hosted_domains ON organizations.id = organization_google_hosted_domains.organization_id
WHERE
    organizations.project_id = $1
    AND organization_google_hosted_domains.google_hosted_domain = $2
    AND NOT organizations.logins_disabled;

-- name: ListOrganizationsByMicrosoftTenantID :many
SELECT
    organizations.*
FROM
    organizations
    JOIN organization_microsoft_tenant_ids ON organizations.id = organization_microsoft_tenant_ids.organization_id
WHERE
    organizations.project_id = $1
    AND organization_microsoft_tenant_ids.microsoft_tenant_id = $2
    AND NOT organizations.logins_disabled;

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
            AND users.microsoft_user_id = $4))
    AND NOT organizations.logins_disabled;

-- name: ListOrganizationsByMatchingUserInvite :many
SELECT
    organizations.*
FROM
    organizations
    JOIN user_invites ON organizations.id = user_invites.organization_id
WHERE
    organizations.project_id = $1
    AND user_invites.email = $2;

-- name: ListSAMLOrganizations :many
SELECT
    organizations.*
FROM
    organizations
    JOIN organization_domains ON organizations.id = organization_domains.organization_id
WHERE
    organizations.project_id = $1
    AND organizations.log_in_with_saml = TRUE
    AND organization_domains.domain = $2;

-- name: RevokeIntermediateSession :one
UPDATE
    intermediate_sessions
SET
    secret_token_sha256 = NULL
WHERE
    id = $1
RETURNING
    *;

-- name: DeleteIntermediateSessionUserInvite :one
DELETE FROM user_invites
WHERE organization_id = $1
    AND email = $2
RETURNING
    *;

-- name: UpdateUserIsOwner :one
UPDATE
    users
SET
    is_owner = $1
WHERE
    id = $2
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

-- name: UpdateIntermediateSessionOrganizationID :one
UPDATE
    intermediate_sessions
SET
    organization_id = $1
WHERE
    id = $2
RETURNING
    *;

-- name: UpdateIntermediateSessionNewUserPasswordBcrypt :one
UPDATE
    intermediate_sessions
SET
    new_user_password_bcrypt = $1,
    password_verified = TRUE
WHERE
    id = $2
RETURNING
    *;

-- name: UpdateIntermediateSessionEmailVerificationChallengeSha256 :one
UPDATE
    intermediate_sessions
SET
    email_verification_challenge_sha256 = $1
WHERE
    id = $2
RETURNING
    *;

-- name: UpdateIntermediateSessionEmailVerificationChallengeCompleted :one
UPDATE
    intermediate_sessions
SET
    email_verification_challenge_completed = TRUE
WHERE
    id = $1
RETURNING
    *;

-- name: GetEmailVerifiedByGoogleUserID :one
SELECT
    EXISTS (
        SELECT
            *
        FROM
            oauth_verified_emails
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
            oauth_verified_emails
        WHERE
            project_id = $1
            AND email = $2
            AND microsoft_user_id = $3);

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

-- name: UpdateUserPasswordBcrypt :one
UPDATE
    users
SET
    password_bcrypt = $1
WHERE
    id = $2
RETURNING
    *;

-- name: GetUserImpersonationTokenBySecretTokenSHA256 :one
SELECT
    *
FROM
    user_impersonation_tokens
WHERE
    secret_token_sha256 = $1
    AND expire_time > now();

-- name: CreateImpersonatedSession :one
INSERT INTO sessions (id, user_id, expire_time, refresh_token_sha256, impersonator_user_id)
    VALUES ($1, $2, $3, $4, $5)
RETURNING
    *;

-- name: RevokeUserImpersonationToken :one
UPDATE
    user_impersonation_tokens
SET
    secret_token_sha256 = NULL
WHERE
    id = $1
RETURNING
    *;

-- name: GetUserHasActivePasskey :one
SELECT
    EXISTS (
        SELECT
            *
        FROM
            passkeys
        WHERE
            user_id = $1
            AND disabled = FALSE);

-- name: UpdateIntermediateSessionRegisterPasskey :one
UPDATE
    intermediate_sessions
SET
    passkey_credential_id = $1,
    passkey_public_key = $2,
    passkey_aaguid = $3,
    passkey_rp_id = $4,
    passkey_verified = TRUE,
    update_time = now()
WHERE
    id = $5
RETURNING
    *;

-- name: CreatePasskey :one
INSERT INTO passkeys (id, user_id, credential_id, public_key, aaguid, rp_id)
    VALUES ($1, $2, $3, $4, $5, $6)
RETURNING
    *;

-- name: UpdateIntermediateSessionPasskeyVerifyChallengeSHA256 :one
UPDATE
    intermediate_sessions
SET
    passkey_verify_challenge_sha256 = $1,
    update_time = now()
WHERE
    id = $2
RETURNING
    *;

-- name: UpdateIntermediateSessionPasskeyVerified :one
UPDATE
    intermediate_sessions
SET
    passkey_verify_challenge_sha256 = NULL,
    passkey_verified = TRUE,
    update_time = now()
WHERE
    id = $1
RETURNING
    *;

-- name: GetUserPasskeyCredentialIDs :many
SELECT
    credential_id
FROM
    passkeys
WHERE
    user_id = $1;

-- name: GetPasskeyByCredentialID :one
SELECT
    *
FROM
    passkeys
WHERE
    credential_id = $1
    AND user_id = $2;

-- name: UpdateIntermediateSessionAuthenticatorAppSecretCiphertext :one
UPDATE
    intermediate_sessions
SET
    authenticator_app_secret_ciphertext = $1,
    update_time = now()
WHERE
    id = $2
RETURNING
    *;

-- name: UpdateIntermediateSessionAuthenticatorAppBackupCodeBcrypts :one
UPDATE
    intermediate_sessions
SET
    authenticator_app_recovery_code_bcrypts = $1,
    update_time = now()
WHERE
    id = $2
RETURNING
    *;

-- name: UpdateIntermediateSessionAuthenticatorAppVerified :one
UPDATE
    intermediate_sessions
SET
    authenticator_app_verified = TRUE,
    update_time = now()
WHERE
    id = $1
RETURNING
    *;

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

-- name: UpdateUserAuthenticatorAppRecoveryCodeBcrypts :one
UPDATE
    users
SET
    authenticator_app_recovery_code_bcrypts = $1
WHERE
    id = $2
RETURNING
    *;

-- name: UpdateUserFailedAuthenticatorAppAttempts :one
UPDATE
    users
SET
    failed_authenticator_app_attempts = $1
WHERE
    id = $2
RETURNING
    *;

-- name: UpdateUserAuthenticatorAppLockoutExpireTime :one
UPDATE
    users
SET
    authenticator_app_lockout_expire_time = $1
WHERE
    id = $2
RETURNING
    *;

-- name: UpdateIntermediateSessionPrimaryLoginFactor :one
UPDATE
    intermediate_sessions
SET
    primary_login_factor = $1
WHERE
    id = $2
RETURNING
    *;

-- name: CreateProject :one
INSERT INTO projects (id, organization_id, display_name, vault_domain, email_send_from_domain, log_in_with_google, log_in_with_microsoft, log_in_with_password, log_in_with_saml)
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING
    *;

-- name: CreateUserInvite :one
INSERT INTO user_invites (id, organization_id, email, is_owner)
    VALUES ($1, $2, $3, $4)
RETURNING
    *;

