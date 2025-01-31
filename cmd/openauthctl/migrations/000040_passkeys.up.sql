create table passkeys
(
    id              uuid                     not null primary key,
    user_id         uuid                     not null references users (id),
    create_time     timestamp with time zone not null default now(),
    update_time     timestamp with time zone not null default now(),
    credential_id   bytea                    not null,
    public_key bytea                  not null,
    aaguid          varchar                  not null
);

alter table intermediate_sessions
    add column passkey_credential_id bytea,
    add column passkey_public_key bytea,
    add column passkey_aaguid varchar,
    add column passkey_verify_challenge_sha256 bytea,
    add column passkey_verified bool not null default false;
