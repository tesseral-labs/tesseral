create type primary_login_factor as enum ('email', 'google_oauth', 'microsoft_oauth');

alter table intermediate_sessions add column primary_login_factor primary_login_factor;