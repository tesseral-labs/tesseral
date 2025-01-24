alter table users
    add column failed_password_attempts     int not null default 0,
    add column password_lockout_expire_time timestamp with time zone;
