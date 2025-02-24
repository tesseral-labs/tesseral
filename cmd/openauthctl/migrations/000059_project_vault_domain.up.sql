alter table projects
    drop column auth_domain,
    drop column custom_auth_domain,
    add column vault_domain varchar not null unique;
