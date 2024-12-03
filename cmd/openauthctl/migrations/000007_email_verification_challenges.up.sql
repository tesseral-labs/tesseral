-- Email verification challenges are used to verify that the user has access to the 
-- email address they are suspected of owning.
-- Either the email address or the google_user_id or the microsoft_user_id must be 
-- set for this record to be valid. This will be enforced via check queries.
create table email_verification_challenges (
  id                            uuid not null primary key,
  -- the intermediate session the challenge is associated with
  intermediate_session_id       uuid not null references intermediate_sessions(id),
  -- the project that the challenge is associated with
  project_id                    uuid not null references projects(id),
  -- the sha256 hash of the code sent to the associated email address
  challenge_sha256              bytea,
  -- the timestamp for when the challenge was completed
  complete_time                 timestamp with time zone,
  -- the timestamp for when the challenge was created
  create_time                   timestamp with time zone not null default now(),
  -- the email address that the challenge is associated with (if user-provided)
  email                         varchar,
  -- the timestamp for when the challenge expires
  expire_time                   timestamp with time zone not null,
  -- the google user id that the challenge is associated with (if google-provided)
  google_user_id                varchar,
  -- the microsoft user id that the challenge is associated with (if microsoft-provided)
  microsoft_user_id             varchar
);

create table verified_emails (
  id                            uuid not null primary key,
  -- the project that the email is associated with
  project_id                    uuid not null references projects(id),
  -- the timestamp for when the email was verified
  create_time                   timestamp with time zone not null default now(),
  -- the email address that is verified
  email                         varchar not null,
  -- the google user id that the email is associated with (if google-provided)
  google_user_id                varchar,
  -- the microsoft user id that the email is associated with (if microsoft-provided)
  microsoft_user_id             varchar
);