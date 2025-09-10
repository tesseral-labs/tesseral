alter table backend_api_keys
  add column authentication_only boolean not null default false;