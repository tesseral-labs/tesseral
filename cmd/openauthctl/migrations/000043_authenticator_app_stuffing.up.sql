alter table users
    drop column failed_authenticator_app_backup_code_attempts,
    drop column authenticator_app_backup_code_lockout_expire_time;

alter table users
    add column failed_authenticator_app_attempts integer not null default 0,
    add column authenticator_app_lockout_expire_time timestamp with time zone;
