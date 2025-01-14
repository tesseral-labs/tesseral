create table project_redirect_uris (
    id uuid primary key not null,
    project_id uuid not null references projects(id),
    uri text not null,
    is_primary boolean not null default false,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now()
);

create unique index unique_primary_per_project
on project_redirect_uris(project_id)
where is_primary = true;

create unique index unique_uri_per_project
on project_redirect_uris(project_id, uri);