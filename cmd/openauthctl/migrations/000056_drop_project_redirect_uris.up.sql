drop table project_redirect_uris;

alter table projects
    add column redirect_uri              varchar,
    add column after_login_redirect_uri  varchar,
    add column after_signup_redirect_uri varchar;
