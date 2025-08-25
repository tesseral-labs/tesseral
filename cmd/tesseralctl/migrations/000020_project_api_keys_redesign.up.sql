alter table project_api_keys
    drop column create_time,
    drop column revoked,
    add column display_name varchar not null,
    alter column secret_token_sha256 drop not null;
