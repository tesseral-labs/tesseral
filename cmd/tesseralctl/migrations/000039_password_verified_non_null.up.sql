alter table intermediate_sessions
    alter column password_verified set not null,
    alter column password_verified set default false;
