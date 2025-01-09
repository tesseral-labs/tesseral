alter table projects
    add column create_time timestamp with time zone not null default now(),
    add column update_time timestamp with time zone not null default now();

alter table project_api_keys
    add column create_time timestamp with time zone not null default now(),
    add column update_time timestamp with time zone not null default now();

alter table organizations
    add column create_time timestamp with time zone not null default now(),
    add column update_time timestamp with time zone not null default now();

alter table saml_connections
    add column update_time timestamp with time zone not null default now();

alter table scim_api_keys
    add column create_time timestamp with time zone not null default now(),
    add column update_time timestamp with time zone not null default now();

alter table users
    alter column update_time set default now();

alter table users
    alter column update_time set not null;
