alter table organizations add column logins_disabled boolean not null default false;

alter table projects add column logins_disabled boolean not null default false;
