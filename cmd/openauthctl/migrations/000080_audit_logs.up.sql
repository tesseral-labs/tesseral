-- Tesseral resource types for audit log events.
create type audit_log_event_resource_type as enum (
    'action',
    'api_key',
    'api_key_role_assignment',
    'audit_log_event',
    'backend_api_key',
    'email_verification_challenge',
    'intermediate_session',
    'organization',
    'organization_google_hosted_domains',
    'organization_microsoft_tenant_ids',
    'passkey',
    'password_reset_code',
    'project',
    'project_ui_settings',
    'project_webhook_settings',
    'publishable_key',
    'role',
    'saml_connection',
    'scim_api_key',
    'session',
    'user',
    'user_authenticator_app_challenge',
    'user_impersonation_token',
    'user_invite',
    'user_role_assignment'
);

-- The history of audit log events.
--
-- Note: We are not interested in foreign keys for most fields. An actor may perform auditable events and 
-- later be deleted or the modified resource may be deleted. In these cases, it is critical that their audit 
-- logs are retained without causing conflicts.
create table audit_log_events (
    -- Unique event ID, generated at time of creation.
    --
    -- This is a UUIDv7 which encodes the creation timestamp, allowing filtering/sorting by this field.
    id uuid not null primary key,

    -- The parent project ID.
    --
    -- All audit log events occur within the context of a project, both those generated for Tesseral
    -- activities and those generated as part of our Audit Logs as a Service feature.
    project_id uuid not null references projects(id),

    -- The actor's organization ID.
    --
    -- For actors operating at the project level and for unauthenticated actors, e.g. for actions performed 
    -- on an intermediate session, this will be `null`.
    -- 
    -- For actors authenticated to a particular organization, this will be the ID of that organization.
    organization_id uuid,

    -- The user who performed the action, if any.
    --
    -- We store this along with the session ID to allow for filtering/sorting by both session and user.
    user_id uuid,

    -- The session of the user who performed the action, if any.
    --
    -- We store this along with the user ID to allow for filtering/sorting by both session and user.
    session_id uuid,

    -- The API key who performed the action, if any.
    api_key_id uuid,

    -- The dogfood user who performed the action, if any.
    --
    -- This is used to track actions performed by dogfood users, which are not associated with a specific organization.
    dogfood_user_id uuid,

    -- The dogfood session of the user who performed the action, if any.
    --
    -- This is used to track actions performed by dogfood users, which are not associated with a specific organization.
    dogfood_session_id uuid,

    -- The backend API key ID of the actor, if any.
    --
    -- This is used to track actions performed by backend API keys, which are not associated with a specific organization.
    backend_api_key_id uuid,

    -- The intermediate session of the entity who performed the action, if any.
    --
    -- This is used for unauthenticated actors, e.g. those who perform actions in the login process, and which are not
    -- associated with a specific organization.
    intermediate_session_id uuid,

    -- The resource type of the Tesseral resource that was acted upon, if any.
    --
    -- For custom audit log events, this will be `null`.
    resource_type audit_log_event_resource_type,

    -- The resource ID of the Tesseral resource that was acted upon, if any.
    --
    -- For custom audit log events, this will be `null`.
    resource_id uuid,

    -- The discrete name of the event.
    --
    -- On an per-project basis, the set of all event names form an enum. Conformance to this enum
    -- is enforced at the application level, outside the database.
    event_name varchar not null,

    -- The time the event occurred.
    event_time timestamp with time zone not null,

    -- The event details as a JSON blob.
    event_details jsonb not null default '{}'::jsonb
);

-- We support filtering/sorting by id (e.g. creation time), event_name, and (actor_type, actor_id).
--
-- We further restrict secondary sorting to just `id desc`. This prevents an explosion of indexes 
-- required for more performant, dynamic sorting.

-- To allow for filtering/sorting by event time.
--
-- Postgres can intelligently use this index to order by id ASC or DESC.
-- See: https://www.postgresql.org/docs/current/indexes-ordering.html

-- The frontend audit log UI (presented in the Vault to organization owners) prioritizes querying by actor,
-- event name, and timestamp.
--
-- These queries are also supported by the backend audit log UI (presented in the Console to project owners),
-- additionally allowing filtering by organization_id (unlike the frontend where this is always singular).
create index on audit_log_events (project_id, id desc);
create index on audit_log_events (project_id, organization_id, id desc) where organization_id is not null;
create index on audit_log_events (project_id, organization_id, user_id, id desc) where organization_id is not null and user_id is not null;
create index on audit_log_events (project_id, organization_id, session_id, id desc) where organization_id is not null and session_id is not null;
create index on audit_log_events (project_id, organization_id, api_key_id, id desc) where organization_id is not null and api_key_id is not null;
create index on audit_log_events (project_id, organization_id, event_name, id desc);

-- The backend audit log UI further allows filtering by resource since the set of all Tesseral resources is fixed.
-- We don't support this on the frontend since the set of resources is dynamic on a per-organization basis and can change over time.
create index on audit_log_events (project_id, resource_type, resource_id, id desc) where resource_type is not null and resource_id is not null;
