create table oidc_connections
(
    id                          uuid                     not null primary key,
    organization_id             uuid                     not null references organizations (id),
    create_time                 timestamp with time zone not null default now(),
    update_time                 timestamp with time zone not null default now(),
    is_primary                  boolean                           default false not null,
    configuration_url           varchar not null,
    issuer                      varchar not null,
    client_id                   varchar not null,
    client_secret_ciphertext    bytea
);

create table oidc_intermediate_sessions
(
    oidc_intermediate_session_id    uuid primary key not null,
    oidc_connection_id              uuid not null references oidc_connections (id),
    code_verifier                   varchar
);

alter table projects
    add column log_in_with_oidc boolean default false not null;

alter table organizations
    add column log_in_with_oidc boolean default false not null;

alter type primary_auth_factor add value 'oidc';
