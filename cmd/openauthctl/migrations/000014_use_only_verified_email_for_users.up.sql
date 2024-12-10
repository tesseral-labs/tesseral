alter table users
  drop constraint users_organization_id_unverified_email_key,
  drop constraint users_organization_id_verified_email_key,
  drop column unverified_email,
  drop column verified_email,
  add column email text not null,
  add column create_time timestamp with time zone not null default now(),
  add column update_time timestamp with time zone,
  add constraint users_organization_id_email_key unique (organization_id, email);