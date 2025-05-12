create table roles
(
    id              uuid                     not null primary key,
    project_id      uuid                     not null references projects (id),
    organization_id uuid references organizations (id),
    create_time     timestamp with time zone not null default now(),
    update_time     timestamp with time zone not null default now(),
    display_name    varchar                  not null,
    description     varchar                  not null
);

create table actions
(
    id          uuid    not null primary key,
    project_id  uuid    not null references projects (id),
    name        varchar not null,
    description varchar not null,

    unique (project_id, name)
);

create table role_actions
(
    id        uuid not null primary key,
    role_id   uuid not null references roles (id) on delete cascade,
    action_id uuid not null references actions (id) on delete cascade,

    unique (role_id, action_id)
);

create table user_role_assignments
(
    id      uuid not null primary key,
    role_id uuid not null references roles (id) on delete cascade,
    user_id uuid not null references users (id) on delete cascade,

    unique (role_id, user_id)
);

alter table user_invites
    add column role_id uuid references roles (id);

alter table organizations
    add column custom_roles_enabled boolean not null default false;
