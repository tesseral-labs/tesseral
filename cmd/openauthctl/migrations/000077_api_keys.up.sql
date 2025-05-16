create table api_keys (
    id uuid not null primary key,
    organization_id uuid not null references organizations (id) on delete cascade,
    display_name varchar not null,
    secret_token_sha256 bytea,
    secret_token_suffix varchar(4),
    expire_time timestamp with time zone,
    create_time timestamp with time zone not null default now(),
    update_time timestamp with time zone not null default now()
);

create table api_key_role_assignments(
    id uuid not null primary key,
    api_key_id uuid not null references api_keys (id) on delete cascade,
    role_id uuid not null references roles (id) on delete cascade,
    create_time timestamp with time zone not null default now(),

    unique (api_key_id, role_id)
);

alter table projects 
    add column api_keys_enabled boolean not null default false,
    add column api_key_secret_token_prefix varchar;

alter table organizations
    add column api_keys_enabled boolean not null default false;