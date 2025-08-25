alter table intermediate_sessions
    add column google_oauth_state_sha256    bytea,
    add column microsoft_oauth_state_sha256 bytea,
    add column google_hosted_domain varchar,
    add column google_user_id varchar,
    add column microsoft_tenant_id varchar,
    add column microsoft_user_id varchar;
