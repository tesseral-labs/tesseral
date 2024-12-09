alter table users
  drop constraint users_organization_id_unverified_email_key,
  drop constraint users_organization_id_verified_email_key,
  drop column unverified_email,
  drop column verified_email,
  add column email text not null,
  add constraint users_organization_id_email_key unique (organization_id, email);