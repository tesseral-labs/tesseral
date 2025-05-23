-- The history of audit log events generated for an organization.
create table organization_audit_log_events (
    -- Unique event ID, generated at time of creation.
    id uuid not null primary key,

    -- The parent organization ID.
    --
    -- When organizations are deleted, we delete the associated logs. Organizations concerned
    -- with maintaining these logs after being removed from Tesseral should request a backup
    -- beforehand.
    organization_id uuid not null references organizations(id)
        on delete cascade,

    -- The time the event occurred.
    event_time timestamptz not null,

    -- The discrete name of the event.
    --
    -- On an per-project basis, the set of all event names form an enum. Conformance to this enum
    -- is enforced at the application level, outside the database.
    event_name varchar not null,

    -- The type of actor who performed the event.
    --
    -- Must be either `user` or `api_key` (enforced at application level).
    --
    -- The separation of `actor_type` and `actor_id` this like ensures that we can
    -- index efficiently when filtering/sorting by actor ID while also allowing further
    -- type/ID pairs to be used.
    actor_type varchar not null,

    -- The ID of the acting API key or user.
    --
    -- We are not interested in foreign keys for this field. A user or API key may perform
    -- auditable events and later be deleted. In this case, it is critical that their audit
    -- logs are retained without causing conflicts.
    actor_id uuid not null,

    -- The event details as a JSON blob.
    event_details jsonb not null default '{}'::jsonb
);

-- We support filtering/sorting by event_time, event_name, and (actor_type, actor_id).
--
-- We further restrict secondary sorting to just `event_time desc`. This prevents an explosion
-- of indexes required for more performant, dynamic sorting.

-- To allow for filtering/sorting by create time.
--
-- Postgres can intelligently use this index to order by event_time ASC or DESC.
-- See: https://www.postgresql.org/docs/current/indexes-ordering.html
create index on organization_audit_log_events (organization_id, event_time desc);

-- To allow for filtering/sorting by event name.
create index on organization_audit_log_events (organization_id, event_name, event_time desc);

-- To allow for filtering/sorting by actor ID.
create index on organization_audit_log_events (organization_id, actor_type, actor_id, event_time desc);
