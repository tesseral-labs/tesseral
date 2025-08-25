create table user_invites
(
    id              uuid                     not null primary key,
    organization_id uuid                     not null references organizations (id),
    create_time     timestamp with time zone not null default now(),
    update_time     timestamp with time zone not null default now(),
    email           varchar                  not null,

    unique (organization_id, email)
);
