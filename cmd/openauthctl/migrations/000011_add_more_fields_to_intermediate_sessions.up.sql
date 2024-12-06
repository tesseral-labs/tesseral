alter table intermediate_sessions
add column google_user_id varchar default null,
add column microsoft_user_id varchar default null;