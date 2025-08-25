alter table projects
    add column organizations_saml_enabled_default boolean not null,
    add column organizations_scim_enabled_default boolean not null;

alter table organizations
    add column saml_enabled boolean not null,
    add column scim_enabled boolean not null;
