create table project_onboarding_progress (
    id uuid primary key default gen_random_uuid(),
    project_id uuid not null references projects(id) on delete cascade,
    configure_authentication_time timestamp with time zone,
    log_in_to_vault_time timestamp with time zone,
    manage_organizations_time timestamp with time zone,
    onboarding_skipped boolean default false,
    create_time timestamp with time zone default now(),
    update_time timestamp with time zone default now(),
    unique (project_id)
);