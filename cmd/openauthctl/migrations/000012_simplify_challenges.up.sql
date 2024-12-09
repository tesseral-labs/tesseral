alter table email_verification_challenges
drop column email,
drop column google_user_id,
drop column microsoft_user_id,
add column intermediate_session_id uuid not null references intermediate_sessions(id) on delete cascade,
add column revoked boolean not null default false;

alter table verified_emails
add column google_hosted_domain varchar,
add column microsoft_tenant_id varchar;