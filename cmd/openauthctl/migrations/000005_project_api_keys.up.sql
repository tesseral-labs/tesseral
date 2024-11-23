create table project_api_keys
(
    id                  uuid        not null primary key,
    project_id          uuid        not null references projects (id),
    create_time         timestamptz not null,
    revoked             boolean     not null,
    secret_token_sha256 bytea       not null
);
