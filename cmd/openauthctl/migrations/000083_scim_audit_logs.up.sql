alter table audit_log_events add column actor_scim_api_key_id uuid;

-- No longer needed; the DELETE operation in SCIM now deletes the user instead of deactivating them.
alter table users drop column deactivate_time;
