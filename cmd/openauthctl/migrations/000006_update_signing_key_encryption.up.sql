alter table intermediate_session_signing_keys
drop column public_key;

alter table session_signing_keys
drop column public_key;
