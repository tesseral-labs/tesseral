alter table organizations
    drop column override_log_in_with_google_enabled,
    drop column override_log_in_with_microsoft_enabled,
    drop column override_log_in_with_password_enabled,

    add column disable_log_in_with_google    boolean,
    add column disable_log_in_with_microsoft boolean,
    add column disable_log_in_with_password  boolean,

    add constraint no_disable_logins_unless_override_log_in_methods check (
        override_log_in_methods or (
            disable_log_in_with_google is null and
            disable_log_in_with_microsoft is null and
            disable_log_in_with_password is null
        )
    );
