alter table projects
    add column custom_email_verify_email   bool not null default false,
    add column custom_email_password_reset bool not null default false,
    add column custom_email_user_invite    bool not null default false;
