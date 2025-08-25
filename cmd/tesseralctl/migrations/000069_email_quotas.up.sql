alter table projects
    add column email_quota_daily int;

create table project_email_quota_daily_usage
(
    project_id  uuid not null primary key references projects (id),
    date        date not null,
    quota_usage int  not null,

    unique (project_id, date)
);
