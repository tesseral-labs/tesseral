alter table project_webhook_settings 
  add column direct_webhook_url varchar,
  alter column app_id drop not null;