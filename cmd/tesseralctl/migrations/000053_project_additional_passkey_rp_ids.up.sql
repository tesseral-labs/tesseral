create table project_passkey_rp_ids
(
    project_id uuid    not null references projects (id) on delete cascade,
    rp_id      varchar not null,

    unique (project_id, rp_id)
);

alter table passkeys
    add column disabled boolean not null default false,
    add column rp_id    varchar not null;
