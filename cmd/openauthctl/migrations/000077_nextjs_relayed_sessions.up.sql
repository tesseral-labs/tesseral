alter table intermediate_sessions
    add column redirect_uri                                varchar,
    add column return_relayed_session_token_as_query_param boolean not null default false;
