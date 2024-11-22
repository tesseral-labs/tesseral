create table project_api_keys
(
    id                  uuid  not null primary key,
    project_id          uuid  not null references projects (id),
    secret_token_sha256 bytea not null
);
