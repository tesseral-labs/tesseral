create table relayed_sessions
(
    session_id                        uuid                     not null primary key,
    relayed_session_token_expire_time timestamp with time zone not null,
    relayed_session_token_sha256      bytea unique,
    state                             varchar,
    relayed_refresh_token_sha256      bytea unique
);

alter table intermediate_sessions
    add column relayed_session_state varchar;
