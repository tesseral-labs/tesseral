create table vault_domain_settings
(
    project_id     uuid not null references projects (id) primary key,
    pending_domain varchar not null
);
