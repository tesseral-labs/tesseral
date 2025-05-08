create table project_webhook_settings (
  id uuid not null primary key,
  project_id uuid not null references projects(id) on delete cascade,
  app_id varchar not null,
  create_time timestamp with time zone default now() not null,
  update_time timestamp with time zone default now() not null
);