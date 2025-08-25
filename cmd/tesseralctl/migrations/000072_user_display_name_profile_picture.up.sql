alter table users
    add column display_name        varchar,
    add column profile_picture_url varchar;

alter table intermediate_sessions
    add column user_display_name   varchar,
    add column profile_picture_url varchar;
