create type auth_method as enum ('email', 'google', 'microsoft');

-- temprorary session created after the first step of the login process.
-- this session has a much more limited permissions scope than the final session.
create table intermediate_sessions(
  id                            uuid not null primary key,
  project_id                    uuid not null references projects(id),
  unverified_email              varchar,
  verified_email                varchar,
  create_time                   timestamp with time zone not null default now(),
  -- the timestamp for when the session expires
  expire_time                   timestamp with time zone not null,
  token                         varchar not null,
  token_sha256                  bytea,
  revoked                       boolean not null default false
);

create table intermediate_session_signing_keys(
  id                          uuid not null primary key,
  project_id                  uuid not null references projects(id),
  public_key                  bytea not null,
  private_key_cipher_text     bytea not null,
  create_time                 timestamp with time zone not null default now(),
  expire_time                 timestamp with time zone not null
);

-- a user's email address is verified by sending a verification challenge to the email address
-- and then having the user provide the verification token
create table method_verification_challenges (
  id                        uuid NOT NULL,
  -- verifications are scoped to the project to ensure that each auth method
  -- has a unique set of verification challenges and is only verified once 
  -- per-project
  project_id                uuid NOT NULL references projects(id),
  complete_time             timestamp with time zone,
  -- the email address to which the verification challenge was sent
  intermediate_session_id   uuid NOT NULL references intermediate_sessions(id),
  -- the auth method for which the verification challenge was sent
  auth_method               auth_method NOT NULL,
  expire_time               timestamp with time zone NOT NULL,
  -- the verification token that the user must provide to verify the email address
  secret_token_sha256       bytea NOT NULL
);

-- final session created after the user has verified their email address and/or
-- completed the login process (if verification is no longer required).
create table sessions
(
  id                            uuid not null primary key,
  user_id                       uuid not null references users(id),
  create_time                   timestamp with time zone not null default now(),
  -- the timestamp for when the session expires
  expire_time                   timestamp with time zone default null,
  revoked                       boolean not null default false
);

create table session_signing_keys(
  id                          uuid not null primary key,
  project_id                  uuid not null references projects(id),
  public_key                  bytea not null,
  private_key_cipher_text     bytea not null,
  create_time                 timestamp with time zone not null default now(),
  expire_time                 timestamp with time zone not null
);