create table project_ui_settings
(
    id uuid primary key not null,
    project_id uuid unique not null references projects(id),
    logo_file_key varchar,
    favicon_file_key varchar,
    primary_color varchar,
    detect_dark_mode_enabled boolean not null default true,
    dark_mode_logo_file_key varchar,
    dark_mode_primary_color varchar,
    create_time timestamp with time zone not null default now(),
    update_time timestamp with time zone default now()
);