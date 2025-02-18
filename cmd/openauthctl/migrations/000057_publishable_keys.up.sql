create table publishable_keys
(
    id         uuid not null primary key,
    project_id uuid not null references projects (id),
    create_time timestamp with time zone not null default now(),
    update_time timestamp with time zone not null default now(),
    display_name varchar not null
);
