alter table audit_log_events drop constraint if exists audit_log_events_project_id_fkey;

alter table audit_log_events alter column project_id drop not null;

alter table audit_log_events add constraint audit_log_events_project_id_fkey
    foreign key (project_id) references projects(id) on delete set null;
