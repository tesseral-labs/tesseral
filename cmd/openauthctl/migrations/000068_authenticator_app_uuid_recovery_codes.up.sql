alter table users drop column authenticator_app_recovery_code_bcrypts;
alter table users add column authenticator_app_recovery_code_sha256s bytea[];

alter table intermediate_sessions drop column authenticator_app_recovery_code_bcrypts;
alter table intermediate_sessions add column authenticator_app_recovery_code_sha256s bytea[];
