alter table intermediate_sessions 
  add column new_user_password_bcrypt varchar,
  add column email_verification_challenge_sha256 bytea,
  add column email_verified boolean not null default false;

drop table email_verification_challenges;