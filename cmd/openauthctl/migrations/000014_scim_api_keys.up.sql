create table scim_api_keys
(
    id              uuid                     not null primary key,
    organization_id uuid                     not null references organizations (id),
    create_time     timestamp with time zone not null default now(),
    revoke_time     timestamp with time zone,
    token_sha256    bytea
);
