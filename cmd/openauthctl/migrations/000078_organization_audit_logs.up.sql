-- The history of audit log events generated for an organization.
create table organization_audit_log_events (
    -- Unique event ID, generated at time of creation.
    --
    -- This is a UUIDv7 which encode the creation timestamp, allowing filtering/sorting by this field.
    id uuid not null primary key,

    -- The parent organization ID.
    organization_id uuid not null references organizations(id),

    -- The user associated with this event, if any.
    --
    -- We are not interested in foreign keys for this field. An actor may perform auditable events and 
    -- later be deleted. In this case, it is critical that their audit logs are retained without 
    -- causing conflicts.
    user_id uuid,

    -- The session associated with this event, if any.
    --
    -- We are not interested in foreign keys for this field. An actor may perform auditable events and 
    -- later be deleted. In this case, it is critical that their audit logs are retained without 
    -- causing conflicts.
    session_id uuid,

    -- The API key associated with this event, if any.
    --
    -- We are not interested in foreign keys for this field. An actor may perform auditable events and 
    -- later be deleted. In this case, it is critical that their audit logs are retained without 
    -- causing conflicts.
    api_key_id uuid,

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
create index on organization_audit_log_events (organization_id, id desc);

-- To allow for filtering/sorting by event name.
create index on organization_audit_log_events (organization_id, event_name, id desc);

-- To allow for filtering/sorting by actor ID.
create index on organization_audit_log_events (organization_id, user_id, id desc);
create index on organization_audit_log_events (organization_id, session_id, id desc);
create index on organization_audit_log_events (organization_id, api_key_id, id desc);
