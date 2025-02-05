alter table sessions
    add column last_active_time timestamp with time zone not null default now();
