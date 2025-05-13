alter table projects 
  add column log_in_with_github boolean not null default false,
  add column github_oauth_client_id varchar,
  add column github_oauth_client_secret_ciphertext bytea;

alter table organizations 
  add column log_in_with_github boolean not null default false;

alter table users
  add column github_user_id varchar;

alter table intermediate_sessions
  add column github_user_id varchar,
  add column github_oauth_state_sha256 bytea;

alter table oauth_verified_emails
  add column github_user_id varchar;

alter type primary_auth_factor add value 'github';