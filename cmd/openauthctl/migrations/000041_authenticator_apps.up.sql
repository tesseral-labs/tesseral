alter table users
    add column authenticator_app_secret_ciphertext               bytea,
    add column authenticator_app_backup_code_bcrypts             bytea[],
    add column failed_authenticator_app_backup_code_attempts     int,
    add column authenticator_app_backup_code_lockout_expire_time timestamp with time zone;

alter table intermediate_sessions
    add column authenticator_app_secret_ciphertext   bytea,
    add column authenticator_app_verified            bool not null default false,
    add column authenticator_app_backup_code_bcrypts bytea[];
