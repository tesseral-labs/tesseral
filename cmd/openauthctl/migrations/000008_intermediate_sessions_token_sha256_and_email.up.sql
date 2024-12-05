alter table intermediate_sessions
    drop column token,
    drop column verified_email,
    drop column unverified_email,
    add column email varchar,
    alter column token_sha256 set not null;
