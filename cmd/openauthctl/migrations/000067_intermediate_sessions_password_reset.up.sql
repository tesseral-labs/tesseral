alter table intermediate_sessions
    add column password_reset_code_sha256   bytea,
    add column password_reset_code_verified boolean not null default false;
