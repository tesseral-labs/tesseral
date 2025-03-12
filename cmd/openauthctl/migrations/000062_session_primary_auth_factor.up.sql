alter table intermediate_sessions drop column primary_login_factor;
drop type primary_login_factor;

create type primary_auth_factor as enum ('email', 'google', 'microsoft', 'saml', 'impersonation');

alter table intermediate_sessions add column primary_auth_factor primary_auth_factor;
alter table sessions add column primary_auth_factor primary_auth_factor not null;
