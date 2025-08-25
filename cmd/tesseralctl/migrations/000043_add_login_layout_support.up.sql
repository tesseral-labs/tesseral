create type log_in_layout as enum ('centered', 'side_by_side');

alter table project_ui_settings add column log_in_layout log_in_layout not null default 'centered';