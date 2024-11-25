-- name: CompleteMethodVerificationChallenge :one
insert into method_verification_challenges (
  id,
  complete_time
) values (
  $1,
  $2
)
returning *;

-- name: CreateMethodVerificationChallenge :one
insert into method_verification_challenges (
  id,
  project_id,
  complete_time,
  intermediate_session_id,
  auth_method,
  expire_time,
  secret_token_sha256
) values (
  $1,
  $2,
  $3,
  $4,
  $5,
  $6,
  $7
)
returning *;

-- name: CreateIntermediateSession :one
insert into intermediate_sessions (
  id,
  project_id,
  unverified_email,
  verified_email,
  expire_time,
  token,
  token_sha256
) values (
  $1,
  $2,
  $3,
  $4,
  $5,
  $6,
  $7
)
returning *;

-- name: CreateIntermediateSessionSigningKey :one
insert into intermediate_session_signing_keys (
  id,
  project_id,
  public_key,
  private_key_cipher_text,
  expire_time
) values (
  $1,
  $2,
  $3,
  $4,
  $5
)
returning *;

-- name: CreateOrganization :one
insert into organizations (
  id,
  project_id,
  display_name,
  google_hosted_domain,
  microsoft_tenant_id,
  override_log_in_with_google_enabled,
  override_log_in_with_microsoft_enabled,
  override_log_in_with_password_enabled
) values (
  $1,
  $2,
  $3,
  $4,
  $5,
  $6,
  $7,
  $8
)
returning *;

-- name: CreateProject :one
insert into projects (
  id,
  organization_id,
  log_in_with_password_enabled,
  log_in_with_google_enabled,
  log_in_with_microsoft_enabled,
  google_oauth_client_id,
  google_oauth_client_secret,
  microsoft_oauth_client_id,
  microsoft_oauth_client_secret
) values (
  $1,
  $2,
  $3,
  $4,
  $5,
  $6,
  $7,
  $8,
  $9
)
returning *;

-- name: CreateSession :one
insert into sessions (
  id,
  user_id,
  create_time,
  expire_time,
  revoked
) values (
  $1,
  $2,
  $3,
  $4,
  $5
)
returning *;

-- name: CreateSessionSigningKey :one
insert into session_signing_keys (
  id,
  project_id,
  public_key,
  private_key_cipher_text,
  expire_time
) values (
  $1,
  $2,
  $3,
  $4,
  $5
)
returning *;

-- name: CreateUser :one
insert into users (
  id,
  organization_id,
  unverified_email,
  verified_email,
  password_bcrypt,
  google_user_id,
  microsoft_user_id
) values (
  $1,
  $2,
  $3,
  $4,
  $5,
  $6,
  $7
)
returning *;

-- name: CreateGoogleUser :one
insert into users (
  id,
  organization_id,
  google_user_id,
  verified_email
) values (
  $1,
  $2,
  $3,
  $4
)
returning *;

-- name: CreateMicrosoftUser :one
insert into users (
  id,
  organization_id,
  microsoft_user_id,
  verified_email
) values (
  $1,
  $2,
  $3,
  $4
)
returning *;

-- name: CreateUnverifiedUser :one
insert into users (
  id,
  organization_id,
  unverified_email
) values (
  $1,
  $2,
  $3
)
returning *;

-- name: GetIntermediateSessionByID :one
select * from intermediate_sessions where id = $1;

-- name: GetIntermediateSessionSigningKeyByID :one
select * from intermediate_session_signing_keys where id = $1;

-- name: GetIntermediateSessionSigningKeyByProjectID :one
select * from intermediate_session_signing_keys where project_id = $1 order by create_time desc limit 1;

-- name: GetMethodVerificationChallengeByID :one
select * from method_verification_challenges where id = $1;

-- name: GetOrganizationByID :one
select * from organizations where id = $1;

-- name: GetProjectByID :one
select * from projects where id = $1;

-- name: GetSessionByID :one
select * from sessions where id = $1;

-- name: GetSessionSigningKeyByID :one
select * from session_signing_keys where id = $1;

-- name: GetSessionSigningKeyByProjectID :one
select * from session_signing_keys where project_id = $1 order by create_time desc limit 1;

-- name: GetOrganizationByGoogleHostedDomain :one
select * from organizations where google_hosted_domain = $1;

-- name: GetUserByEmail :one
select * from users where verified_email = $1 or unverified_email = $1;

-- name: GetUserByID :one
select * from users where id = $1;

-- name: GetUserByGoogleUserID :one
select * from users where organization_id = $1 and google_user_id = $2;

-- name: GetUserByMicrosoftUserID :one
select * from users where organization_id = $1 and microsoft_user_id = $2;

-- name: GetUserBySessionID :one
select * from users where users.id = (select user_id from sessions where sessions.id = $1);

-- name: GetUserByUnverifiedEmail :one
select * from users where unverified_email = $1;

-- name: GetUserByVerifiedEmail :one
select * from users where verified_email = $1;

-- name: ListOrganizations :many
select org.* from organizations as org;

-- name: ListOrganizationsByProjectId :many
select o.*
from organizations as o 
join projects as p
on o.project_id = p.id
where o.project_id = $1 
order by o.display_name limit $2;

-- name: ListOrganizationsByProjectIdAndEmail :many
select o.*
from organizations as o
join users as u 
on o.id = u.organization_id
where o.project_id = $1 
and u.verified_email = $2
order by o.display_name limit $3;

-- name: ListProjects :many
select * from projects order by id limit $1;

-- name: ListUsersByEmail :many
select * from users where unverified_email = $1 or verified_email = $1;

-- name: ListUsersByOrganization :many
select * from users where organization_id = $1 order by id limit $2;

-- name: RevokeIntermediateSession :one
update intermediate_sessions set revoked = true where id = $1 returning *;

-- name: RevokeSession :one
update sessions set revoked = true where id = $1 returning *;

-- name: UpdateOrganization :one
update organizations set
  display_name = $2,
  google_hosted_domain = $3,
  microsoft_tenant_id = $4,
  override_log_in_with_password_enabled = $5,
  override_log_in_with_google_enabled = $6,
  override_log_in_with_microsoft_enabled = $7
where id = $1 returning *;

-- name: UpdateOrganizationDisplayName :one
update organizations set display_name = $2 where id = $1 returning *;

-- name: UpdateOrganizationGoogleHostedDomain :one
update organizations set google_hosted_domain = $2 where id = $1 returning *;

-- name: UpdateOrganizationMicrosoftTenantID :one
update organizations set microsoft_tenant_id = $2 where id = $1 returning *;

-- name: UpdateOrganizationOverrides :one
update organizations set
  override_log_in_with_password_enabled = $2,
  override_log_in_with_google_enabled = $3,
  override_log_in_with_microsoft_enabled = $4
where id = $1 returning *;

-- name: UpdateProject :one
update projects set
  log_in_with_password_enabled = $2,
  log_in_with_google_enabled = $3,
  log_in_with_microsoft_enabled = $4,
  google_oauth_client_id = $5,
  google_oauth_client_secret = $6,
  microsoft_oauth_client_id = $7,
  microsoft_oauth_client_secret = $8
where id = $1 returning *;

-- name: UpdateProjectOrganizationID :one
update projects set organization_id = $2 where id = $1 returning *;

-- name: UpdateProjectGoogleOAuthClient :one
update projects set google_oauth_client_id = $2, google_oauth_client_secret = $3 where id = $1 returning *;

-- name: UpdateProjectMicrosoftOAuthClient :one
update projects set microsoft_oauth_client_id = $2, microsoft_oauth_client_secret = $3 where id = $1 returning *;

-- name: UpdateProjectLoginMethods :one
update projects set
  log_in_with_password_enabled = $2,
  log_in_with_google_enabled = $3,
  log_in_with_microsoft_enabled = $4
where id = $1 returning *;

-- name: UpdateUser :one
update users set
  organization_id = $2,
  unverified_email = $3,
  verified_email = $4,
  password_bcrypt = $5,
  google_user_id = $6,
  microsoft_user_id = $7
where id = $1 returning *;

-- name: UpdateUserPassword :one
update users set password_bcrypt = $2 where id = $1 returning *;

-- name: UpdateUserGoogleUserID :one
update users set google_user_id = $2 where id = $1 returning *;

-- name: UpdateUserMicrosoftUserID :one
update users set microsoft_user_id = $2 where id = $1 returning *;

-- name: VerifyIntermediateSessionEmail :one
update intermediate_sessions set unverified_email = null, verified_email = $2 where id = $1 returning *;

-- name: VerifyUserEmail :one
update users set unverified_email = null, verified_email = $2 where id = $1 returning *;

-- name: CountAllProjects :one
select count(*) from projects;
