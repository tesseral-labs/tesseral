drop table intermediate_session_signing_keys;

alter table intermediate_sessions
    add column update_time         timestamp with time zone not null default now(),
    drop column revoked,
    drop column token_sha256,
    add column secret_token_sha256 bytea;

alter table verified_emails
    drop column google_hosted_domain,
    drop column microsoft_tenant_id,
    add constraint oauth_user_ids_not_all_blank check (google_user_id is not null or
                                                       microsoft_user_id is not null);

alter table email_verification_challenges
    drop column revoked,
    drop column project_id;

alter table organizations
    drop column google_hosted_domain,
    drop column microsoft_tenant_id;

create table organization_google_hosted_domains
(
    id                   uuid    not null primary key,
    organization_id      uuid    not null references organizations (id) on delete cascade,
    google_hosted_domain varchar not null
);

create table organization_microsoft_tenant_ids
(
    id                  uuid    not null primary key,
    organization_id     uuid    not null references organizations (id) on delete cascade,
    microsoft_tenant_id varchar not null
);

alter table users
    add constraint email_not_empty_string check (email != ''),
    add constraint google_user_id_not_empty_string check (google_user_id != ''),
    add constraint microsoft_user_id_not_empty_string check (microsoft_user_id != '');
