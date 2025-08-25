drop table project_passkey_rp_ids;

create table project_trusted_domains
(
    id         uuid    not null primary key,
    project_id uuid    not null references projects (id),
    domain     varchar not null,

    unique (project_id, domain)
);
