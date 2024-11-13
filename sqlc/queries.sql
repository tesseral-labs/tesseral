-- name: CreateOrganization :one
insert into organizations (
  id, 
  project_id, 
  display_name, 
  google_hosted_domain, 
  microsoft_tenant_id
) values (
  $1, 
  $2, 
  $3, 
  $4, 
  $5
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

-- name: CreateUser :one
insert into users (
  id, 
  organization_id,
  verified_email
) values (
  $1, 
  $2, 
  $3
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

-- name: GetOrganizationByID :one
select * from organizations where id = $1;

-- name: GetProjectByID :one
select * from projects where id = $1;

-- name: GetOrganizationByGoogleHostedDomain :one
select * from organizations where google_hosted_domain = $1;

-- name: GetUserByEmail :one
select * from users where verified_email = $1 or unverified_email = $1;

-- name: GetUserByID :one
select * from users where organization_id = $1 and id = $2;

-- name: GetUserByGoogleUserID :one
select * from users where organization_id = $1 and google_user_id = $2;

-- name: GetUserByMicrosoftUserID :one
select * from users where organization_id = $1 and microsoft_user_id = $2;

-- name: GetUserByUnverifiedEmail :one
select * from users where unverified_email = $1;

-- name: GetUserByVerifiedEmail :one
select * from users where verified_email = $1;

-- name: ListOrganizations :many
select * from organizations;

-- name: ListOrganizationsByProjectId :many
select * from organizations where project_id = $1 order by id limit $2;

-- name: ListProjects :many
select * from projects order by id limit $1;

-- name: ListUsersByEmail :many
select * from users where unverified_email = $1 or verified_email = $1;

-- name: ListUsersByOrganization :many
select * from users where organization_id = $1;

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

-- name: UpdateUserPassword :one
update users set password_bcrypt = $2 where id = $1 returning *;

-- name: UpdateUserGoogleUserID :one
update users set google_user_id = $2 where id = $1 returning *;

-- name: UpdateUserMicrosoftUserID :one
update users set microsoft_user_id = $2 where id = $1 returning *;

-- name: VerifyUserEmail :one
update users set unverified_email = null, verified_email = $2 where id = $1 returning *;
