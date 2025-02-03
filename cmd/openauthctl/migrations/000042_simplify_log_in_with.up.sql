alter table projects
    rename column log_in_with_google_enabled to log_in_with_google;

alter table projects
    rename column log_in_with_microsoft_enabled to log_in_with_microsoft;

alter table projects
    rename column log_in_with_password_enabled to log_in_with_password;

alter table projects
    add column log_in_with_authenticator_app boolean not null default false,
    add column log_in_with_passkey           boolean not null default false;

alter table organizations
    drop column override_log_in_methods,
    drop column disable_log_in_with_google,
    drop column disable_log_in_with_microsoft,
    drop column disable_log_in_with_password,
    add column log_in_with_google            boolean not null default false,
    add column log_in_with_microsoft         boolean not null default false,
    add column log_in_with_password          boolean not null default false,
    add column log_in_with_authenticator_app boolean not null default false,
    add column log_in_with_passkey           boolean not null default false,
    add column require_mfa                   boolean not null default false;
