create type audit_log_event_resource_type as enum (
    'api_key',
    'api_key_role_assignment',
    'organization',
    'passkey',
    'role',
    'saml_connection',
    'scim_api_key',
    'session',
    'user',
    'user_invite',
    'user_role_assignment'
);

create table audit_log_events (
    id uuid not null primary key,
    project_id uuid not null references projects(id),
    organization_id uuid,
    user_id uuid,
    session_id uuid,
    api_key_id uuid,
    dogfood_user_id uuid,
    dogfood_session_id uuid,
    backend_api_key_id uuid,
    intermediate_session_id uuid,
    resource_type audit_log_event_resource_type,
    resource_id uuid,
    event_name varchar not null,
    event_time timestamp with time zone not null,
    event_details jsonb not null default '{}'::jsonb
);

create index on audit_log_events (project_id, id desc);
create index on audit_log_events (project_id, organization_id, id desc) where organization_id is not null;
create index on audit_log_events (project_id, organization_id, user_id, id desc) where organization_id is not null and user_id is not null;
create index on audit_log_events (project_id, organization_id, session_id, id desc) where organization_id is not null and session_id is not null;
create index on audit_log_events (project_id, organization_id, api_key_id, id desc) where organization_id is not null and api_key_id is not null;
create index on audit_log_events (project_id, organization_id, event_name, id desc);

create index on audit_log_events (project_id, resource_type, resource_id, id desc) where resource_type is not null and resource_id is not null;

alter table projects add column audit_logs_enabled boolean not null default true;
