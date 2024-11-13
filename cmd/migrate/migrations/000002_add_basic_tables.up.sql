create table projects
(
    id                            uuid not null primary key,
    organization_id               uuid not null,

    -- todo: state related to hosted pages and email DNS settings

    -- settings related to what login methods are enabled for customers
    log_in_with_password_enabled  boolean not null,
    log_in_with_google_enabled    boolean not null,
    log_in_with_microsoft_enabled boolean not null,

    -- log in with google settings
    google_oauth_client_id        varchar,
    google_oauth_client_secret    varchar,

    -- log in with microsoft settings
    microsoft_oauth_client_id     varchar,
    microsoft_oauth_client_secret varchar
);

create table organizations
(
    id                   uuid    not null primary key,
    project_id           uuid    not null references projects (id),
    display_name         varchar not null,

    -- per-customer override of supported login methods
    override_log_in_with_password_enabled  boolean,
    override_log_in_with_google_enabled    boolean,
    override_log_in_with_microsoft_enabled boolean,

    -- corresponds to `hd` in a Google OIDC JWT: https://developers.google.com/identity/openid-connect/openid-connect
    -- (unless no `hd` is present)
    --
    -- todo: how do we deal with the fact that you can change your `hd` over
    -- time, such as if you change your primary domain in Google Workspace?
    google_hosted_domain varchar,

    -- organizations "own" a google hosted domain within a project
    unique (project_id, google_hosted_domain),

    -- corresponds to `tid` in a Microsoft OIDC JWT: https://learn.microsoft.com/en-us/entra/identity-platform/id-token-claims-reference
    -- (unless it's the documented "public" tenant)
    microsoft_tenant_id  varchar,

    -- organizations "own" a microsoft tenant ID within a project
    unique (project_id, microsoft_tenant_id)
);

create table users
(
    id                uuid    not null primary key,
    organization_id   uuid    not null references organizations (id),

    -- when a user is created with a user-provided email address, we store it here until verified
    unverified_email  varchar,
    
    unique (organization_id, unverified_email),

    -- when a user is created with a service-provided email address (Google, Microsoft), we store it here
    -- - additionally, we store the verified email address here after manual verification by the user
    verified_email    varchar,

    unique (organization_id, verified_email),

    -- todo what additional properties (name, profile picture, etc.) do we bless?

    -- only present if user is using username/password
    password_bcrypt   varchar,

    -- corresponds to `sub` in a Google OIDC JWT: https://developers.google.com/identity/openid-connect/openid-connect
    google_user_id    varchar,

    unique (organization_id, google_user_id),

    -- corresponds to `oid` in a Microsoft OIDC JWT: https://learn.microsoft.com/en-us/entra/identity-platform/id-token-claims-reference
    microsoft_user_id varchar,

    unique (organization_id, microsoft_user_id)
);

alter table projects
    add constraint projects_organization_id_fkey
        foreign key (organization_id) references organizations (id);
