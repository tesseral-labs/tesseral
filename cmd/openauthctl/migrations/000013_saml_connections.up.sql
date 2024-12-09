create table saml_connections
(
    id                   uuid                     not null primary key,
    organization_id      uuid                     not null references organizations (id),
    create_time          timestamp with time zone not null default now(),
    is_primary           boolean                           default false not null,
    idp_redirect_url     varchar,
    idp_x509_certificate bytea,
    idp_entity_id        varchar
);

create table organization_domains
(
    id              uuid    not null primary key,
    organization_id uuid    not null references projects (id),
    domain          varchar not null,

    unique (organization_id, domain)
);
