alter table projects
    add column log_in_with_email boolean not null default false,
    add column log_in_with_saml  boolean not null default false;

alter table organizations
    drop column saml_enabled,
    add column log_in_with_email boolean not null default false,
    add column log_in_with_saml  boolean not null default false;
