drop table oidc_intermediate_sessions;

alter table intermediate_sessions
    add column oidc_state varchar,
    add column oidc_code_verifier varchar,
    add column verified_oidc_connection_id uuid;
