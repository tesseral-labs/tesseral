create table project_webhook_settings (
  id uuid not null primary key,
  project_id uuid not null,
  app_id varchar not null,
  created_at timestamp with time zone default now() not null,
  updated_at timestamp with time zone default now() not null
);